package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/maintain"
	"github.com/spf13/cobra"
)

func newMaintainCmd() *cobra.Command {
	dryRun, batchSize := false, 1000
	command := &cobra.Command{Use: "maintain", Short: "Enforce TTL and the 50GB capacity policy", RunE: func(cmd *cobra.Command, _ []string) error {
		_, db, _, err := openIngestService(cmd.Context())
		if err != nil {
			return err
		}
		defer db.Close()
		report, err := (maintain.Service{DB: db}).Run(cmd.Context(), batchSize, dryRun, time.Now().UTC())
		if err != nil {
			return err
		}
		rows, _ := json.Marshal(report.Rows)
		fmt.Fprintf(cmd.OutOrStdout(), "dry_run=%t before_bytes=%d after_bytes=%d raw_ttl_days=%d rows=%s\n", report.DryRun, report.BeforeBytes, report.AfterBytes, report.Policy.RawPayloadDays, rows)
		return nil
	}}
	command.Flags().BoolVar(&dryRun, "dry-run", false, "report without deleting or pausing")
	command.Flags().IntVar(&batchSize, "batch-size", 1000, "maximum rows deleted per transaction")
	return command
}
