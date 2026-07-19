package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"
)

const statementDelimiter = "-- fnctl:statement"

// Migration is one immutable, checksummed schema change.
type Migration struct {
	Version    string
	Name       string
	Statements []string
	Checksum   [32]byte
}

// AppliedMigration is the database record for an executed migration.
type AppliedMigration struct {
	Version   string
	Name      string
	Checksum  [32]byte
	AppliedAt time.Time
}

// MigrationResult reports whether a version was applied or already present.
type MigrationResult struct {
	Version string
	Name    string
	Applied bool
}

func splitStatements(source string) ([]string, error) {
	normalized := strings.ReplaceAll(source, "\r\n", "\n")
	normalized = strings.TrimPrefix(normalized, "\ufeff")
	trimmed := strings.TrimSpace(normalized)
	if !strings.HasPrefix(trimmed, statementDelimiter) {
		return nil, fmt.Errorf("migration must begin with %q", statementDelimiter)
	}
	parts := strings.Split(trimmed, statementDelimiter)
	statements := make([]string, 0, len(parts)-1)
	for _, part := range parts[1:] {
		statement := strings.TrimSpace(part)
		if statement == "" {
			return nil, fmt.Errorf("migration contains an empty statement")
		}
		statements = append(statements, statement)
	}
	return statements, nil
}

func loadMigrations(files fs.FS) ([]Migration, error) {
	entries, err := fs.ReadDir(files, ".")
	if err != nil {
		return nil, fmt.Errorf("read migrations: %w", err)
	}
	var migrations []Migration
	seen := make(map[string]bool)
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".sql") {
			continue
		}
		if len(name) < 5 || name[4] != '_' {
			return nil, fmt.Errorf("invalid migration filename %q", name)
		}
		version := name[:4]
		if seen[version] {
			return nil, fmt.Errorf("duplicate migration version %s", version)
		}
		seen[version] = true
		body, err := fs.ReadFile(files, name)
		if err != nil {
			return nil, fmt.Errorf("read migration %s: %w", name, err)
		}
		statements, err := splitStatements(string(body))
		if err != nil {
			return nil, fmt.Errorf("parse migration %s: %w", name, err)
		}
		migrations = append(migrations, Migration{
			Version: version, Name: name, Statements: statements, Checksum: sha256.Sum256(body),
		})
	}
	sort.Slice(migrations, func(i, j int) bool { return migrations[i].Version < migrations[j].Version })
	return migrations, nil
}

func verifyMigration(migration Migration, applied AppliedMigration) error {
	if migration.Checksum != applied.Checksum {
		return fmt.Errorf("migration checksum mismatch for %s", migration.Name)
	}
	return nil
}

// Migrate applies trusted embedded SQL one statement at a time. MySQL DDL
// implicitly commits, so every statement must be idempotent and the checksum
// record is written only after every statement in a version succeeds.
func Migrate(
	ctx context.Context,
	admin *sql.DB,
	openDatabase func(context.Context, string) (*sql.DB, error),
	files fs.FS,
	database string,
) ([]MigrationResult, error) {
	migrations, err := loadMigrations(files)
	if err != nil {
		return nil, err
	}
	if len(migrations) == 0 || migrations[0].Version != "0001" {
		return nil, fmt.Errorf("migration 0001 is required")
	}
	adminStatements, databaseStatements, err := bootstrapStatements(migrations[0])
	if err != nil {
		return nil, err
	}
	for _, statement := range adminStatements {
		if _, err := admin.ExecContext(ctx, statement); err != nil {
			return nil, fmt.Errorf("apply %s: %w", migrations[0].Name, err)
		}
	}

	db, err := openDatabase(ctx, database)
	if err != nil {
		return nil, fmt.Errorf("open migrated database: %w", err)
	}
	defer db.Close()
	for _, statement := range databaseStatements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			return nil, fmt.Errorf("apply %s: %w", migrations[0].Name, err)
		}
	}
	applied, err := readAppliedMigrations(ctx, db)
	if err != nil {
		return nil, err
	}

	results := make([]MigrationResult, 0, len(migrations))
	for _, migration := range migrations {
		if existing, ok := applied[migration.Version]; ok {
			if err := verifyMigration(migration, existing); err != nil {
				return nil, err
			}
			results = append(results, MigrationResult{Version: migration.Version, Name: migration.Name})
			continue
		}
		if migration.Version != "0001" {
			for _, statement := range migration.Statements {
				if _, err := db.ExecContext(ctx, statement); err != nil {
					return nil, fmt.Errorf("apply %s: %w", migration.Name, err)
				}
			}
		}
		if err := recordMigration(ctx, db, migration); err != nil {
			return nil, err
		}
		results = append(results, MigrationResult{Version: migration.Version, Name: migration.Name, Applied: true})
	}
	return results, nil
}

func bootstrapStatements(migration Migration) ([]string, []string, error) {
	if migration.Version != "0001" || len(migration.Statements) < 2 {
		return nil, nil, fmt.Errorf("migration 0001 must create the database and migration table")
	}
	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(migration.Statements[0])), "CREATE DATABASE") {
		return nil, nil, fmt.Errorf("migration 0001 first statement must create the database")
	}
	return migration.Statements[:1], migration.Statements[1:], nil
}

func readAppliedMigrations(ctx context.Context, db *sql.DB) (map[string]AppliedMigration, error) {
	rows, err := db.QueryContext(ctx, `SELECT version, name, checksum, applied_at FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("read applied migrations: %w", err)
	}
	defer rows.Close()
	applied := make(map[string]AppliedMigration)
	for rows.Next() {
		var item AppliedMigration
		var checksum []byte
		if err := rows.Scan(&item.Version, &item.Name, &checksum, &item.AppliedAt); err != nil {
			return nil, fmt.Errorf("scan applied migration: %w", err)
		}
		if len(checksum) != sha256.Size {
			return nil, fmt.Errorf("migration %s has invalid checksum length", item.Version)
		}
		copy(item.Checksum[:], checksum)
		applied[item.Version] = item
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read applied migrations: %w", err)
	}
	return applied, nil
}

func recordMigration(ctx context.Context, db *sql.DB, migration Migration) error {
	_, err := db.ExecContext(ctx,
		`INSERT INTO schema_migrations (version, name, checksum, applied_at) VALUES (?, ?, ?, UTC_TIMESTAMP(6))`,
		migration.Version, migration.Name, migration.Checksum[:])
	if err != nil {
		return fmt.Errorf("record migration %s: %w", migration.Name, err)
	}
	return nil
}
