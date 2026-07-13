package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/classify"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

func newClassifyCmd() *cobra.Command {
	limit, daily := 1000, ""
	command := &cobra.Command{Use: "classify", Short: "Score nodes and roll up daily validation metrics", RunE: func(cmd *cobra.Command, _ []string) error {
		_, db, _, err := openIngestService(cmd.Context())
		if err != nil {
			return err
		}
		defer db.Close()
		if daily != "" {
			date, err := time.Parse("2006-01-02", daily)
			if err != nil {
				return err
			}
			rows, err := store.RollupDailyStats(cmd.Context(), db, date)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "daily=%s affected=%d\n", daily, rows)
			return nil
		}
		report, err := (classify.Service{DB: db}).Run(cmd.Context(), limit, time.Now().UTC())
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "candidates=%d classified=%d\n", report.Candidates, report.Classified)
		return nil
	}}
	command.Flags().IntVar(&limit, "limit", 1000, "maximum nodes to classify")
	command.Flags().StringVar(&daily, "daily", "", "roll up one UTC date (YYYY-MM-DD)")
	return command
}
