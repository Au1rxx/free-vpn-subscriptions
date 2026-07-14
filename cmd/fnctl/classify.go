package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/classify"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/geoip"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

func newClassifyCmd() *cobra.Command {
	limit, daily, all := 1000, "", false
	command := &cobra.Command{Use: "classify", Short: "Score nodes and roll up daily validation metrics", RunE: func(cmd *cobra.Command, _ []string) error {
		cfg, db, _, err := openIngestService(cmd.Context())
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
		var network *geoip.NetworkClassifier
		if cfg.GeoIP.Enabled {
			if err := geoip.EnsureDB(cfg.GeoIP.CityDBURL, cfg.GeoIP.CityDBPath); err != nil {
				return err
			}
			if err := geoip.EnsureDB(cfg.GeoIP.ASNDBURL, cfg.GeoIP.ASNDBPath); err != nil {
				return err
			}
			network, err = geoip.OpenNetwork(cfg.GeoIP.CityDBPath, cfg.GeoIP.ASNDBPath)
			if err != nil {
				return err
			}
			defer network.Close()
		}
		remaining, err := store.CountUnclassified(cmd.Context(), db)
		if err != nil {
			return err
		}
		totalCandidates, totalClassified := 0, 0
		for {
			batch := classificationBatchSize(limit, remaining, all)
			if batch == 0 {
				break
			}
			report, err := (classify.Service{DB: db, Network: network}).Run(cmd.Context(), batch, time.Now().UTC())
			if err != nil {
				return err
			}
			totalCandidates += report.Candidates
			totalClassified += report.Classified
			if all {
				remaining = remainingAfterClassification(remaining, report.Classified)
			} else {
				remaining, err = store.CountUnclassified(cmd.Context(), db)
				if err != nil {
					return err
				}
			}
			if !all || remaining == 0 || report.Classified == 0 {
				break
			}
		}
		fmt.Fprintf(cmd.OutOrStdout(), "candidates=%d classified=%d unclassified_remaining=%d\n", totalCandidates, totalClassified, remaining)
		sourceReport, err := classify.RefreshSourceQualities(cmd.Context(), db, time.Now().UTC())
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "source_candidates=%d source_scored=%d source_written=%d\n",
			sourceReport.Candidates, sourceReport.Scored, sourceReport.Written)
		return nil
	}}
	command.Flags().IntVar(&limit, "limit", 1000, "maximum nodes to classify")
	command.Flags().StringVar(&daily, "daily", "", "roll up one UTC date (YYYY-MM-DD)")
	command.Flags().BoolVar(&all, "all", false, "classify every currently unclassified node in bounded batches")
	return command
}

func classificationBatchSize(limit, remaining int, all bool) int {
	if all && remaining > 0 && remaining < limit {
		return remaining
	}
	return limit
}

func remainingAfterClassification(remaining, classified int) int {
	if classified >= remaining {
		return 0
	}
	return remaining - classified
}
