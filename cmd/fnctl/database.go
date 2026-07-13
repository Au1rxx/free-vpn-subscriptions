package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"

	dbmigrations "github.com/Au1rxx/free-vpn-subscriptions/db/migrations"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

func newMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Create and migrate the node database",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := loadDatabaseConfig()
			if err != nil {
				return err
			}
			ctx := cmd.Context()
			admin, err := store.Open(ctx, cfg.Database, "")
			if err != nil {
				return err
			}
			defer admin.Close()
			if _, err := store.CheckServer(ctx, admin); err != nil {
				return err
			}
			results, err := store.Migrate(ctx, admin, func(ctx context.Context, database string) (*sql.DB, error) {
				return store.Open(ctx, cfg.Database, database)
			}, dbmigrations.Files, cfg.Database.Name)
			if err != nil {
				return err
			}
			applied := 0
			for _, result := range results {
				if result.Applied {
					applied++
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "migrations_total=%d applied=%d skipped=%d\n", len(results), applied, len(results)-applied)
			return nil
		},
	}
}

func newDBStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "db-status",
		Short: "Show bounded MySQL and schema health",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := loadDatabaseConfig()
			if err != nil {
				return err
			}
			db, err := store.Open(cmd.Context(), cfg.Database, cfg.Database.Name)
			if err != nil {
				return err
			}
			defer db.Close()
			server, err := store.CheckServer(cmd.Context(), db)
			if err != nil {
				return err
			}
			status, err := store.ReadDatabaseStatus(cmd.Context(), db, cfg.Database.Name, server)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), formatDatabaseStatus(status))
			return nil
		},
	}
}

func loadDatabaseConfig() (*config.Config, error) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, err
	}
	if !cfg.Database.Enabled {
		return nil, fmt.Errorf("database is disabled in %s", cfgPath)
	}
	if cfg.Database.PasswordFile == "" {
		return nil, fmt.Errorf("database password_file is required")
	}
	if cfg.Database.Name != "vpn_nodes" {
		return nil, fmt.Errorf("database name must be vpn_nodes")
	}
	return cfg, nil
}

func formatDatabaseStatus(status store.DatabaseStatus) string {
	return fmt.Sprintf(
		"version=%s tls=%s read_only=%t timezone=%s charset=%s collation=%s\n"+
			"migrations=%d tables=%d empty_table_comments=%d empty_column_comments=%d enabled_policies=%d\n"+
			"data_bytes=%d index_bytes=%d total_bytes=%d\n",
		status.Server.Version, status.Server.Cipher, status.Server.ReadOnly,
		status.Server.TimeZone, status.Server.Charset, status.Server.Collation,
		status.AppliedMigrations, status.BusinessTables, status.EmptyTableComments,
		status.EmptyColumnComments, status.EnabledPolicies, status.DataBytes,
		status.IndexBytes, status.DataBytes+status.IndexBytes,
	)
}
