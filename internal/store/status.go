package store

import (
	"context"
	"database/sql"
	"fmt"
)

// DatabaseStatus is a bounded operational summary without connection secrets.
type DatabaseStatus struct {
	Server            ServerInfo
	AppliedMigrations int
	BusinessTables    int
	DataBytes         uint64
	IndexBytes        uint64
}

// ReadDatabaseStatus returns migration and capacity counters for one schema.
func ReadDatabaseStatus(ctx context.Context, db *sql.DB, database string, server ServerInfo) (DatabaseStatus, error) {
	status := DatabaseStatus{Server: server}
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM schema_migrations`).Scan(&status.AppliedMigrations); err != nil {
		return DatabaseStatus{}, fmt.Errorf("count migrations: %w", err)
	}
	if err := db.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(SUM(data_length), 0), COALESCE(SUM(index_length), 0)
		FROM information_schema.tables WHERE table_schema = ?`, database).Scan(
		&status.BusinessTables, &status.DataBytes, &status.IndexBytes); err != nil {
		return DatabaseStatus{}, fmt.Errorf("read database capacity: %w", err)
	}
	return status, nil
}
