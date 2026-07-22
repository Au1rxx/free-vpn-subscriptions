package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/maintain"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
	"github.com/spf13/cobra"
)

func newMaintainCmd() *cobra.Command {
	dryRun, batchSize := false, 1000
	command := &cobra.Command{Use: "maintain", Short: "Enforce TTL and the 50GB capacity policy", RunE: func(cmd *cobra.Command, _ []string) error {
		cfg, err := loadDatabaseConfig()
		if err != nil {
			return err
		}
		db, err := store.OpenMaintenance(cmd.Context(), cfg.Database, cfg.Database.Name)
		if err != nil {
			return err
		}
		defer db.Close()
		if _, err := store.CheckServer(cmd.Context(), db); err != nil {
			return err
		}
		report, err := (maintain.Service{DB: db}).Run(cmd.Context(), batchSize, dryRun, time.Now().UTC())
		if err != nil {
			return err
		}
		return writeMaintenanceReport(cmd.OutOrStdout(), report)
	}}
	command.Flags().BoolVar(&dryRun, "dry-run", false, "report without deleting or pausing")
	command.Flags().IntVar(&batchSize, "batch-size", 1000, "maximum rows deleted per transaction")
	return command
}

func writeMaintenanceReport(output io.Writer, report maintain.Report) error {
	rows, err := json.Marshal(report.Rows)
	if err != nil {
		return fmt.Errorf("encode maintenance rows: %w", err)
	}
	_, err = fmt.Fprintf(output,
		"dry_run=%t before_bytes=%d after_bytes=%d raw_ttl_days=%d rows=%s storage_metric_rows=%d\n",
		report.DryRun, report.BeforeBytes, report.AfterBytes, report.Policy.RawPayloadDays,
		rows, report.StorageMetricRows)
	return err
}
