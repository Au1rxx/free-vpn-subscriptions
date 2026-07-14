package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/validation"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/verify"
)

func newValidateWorkerCmd() *cobra.Command {
	var once bool
	limit, concurrency := 40, 20
	validatorID := "ai-a1"
	command := &cobra.Command{Use: "validate-worker", Short: "Lease and deeply validate queued nodes", RunE: func(cmd *cobra.Command, _ []string) error {
		cfg, db, _, err := openIngestService(cmd.Context())
		if err != nil {
			return err
		}
		defer db.Close()
		worker := validation.Worker{Queue: validation.SQLQueue{DB: db}, ValidatorID: validatorID,
			Concurrency: concurrency, Lease: 2 * time.Minute,
			Performance: validation.SampleRequest{URL: cfg.Verify.PerformanceURL, Bytes: cfg.Verify.PerformanceBytes,
				Timeout: time.Duration(cfg.Verify.PerformanceTimeoutMS) * time.Millisecond},
			Request: verify.Request{Targets: cfg.Verify.Targets,
				Timeout:    time.Duration(cfg.Verify.TimeoutMS) * time.Millisecond,
				SingBoxBin: cfg.Verify.SingBoxBin, StartupTimeout: time.Duration(cfg.Verify.StartupTimeoutMS) * time.Millisecond}}
		if once {
			report, err := worker.RunOnce(cmd.Context(), limit)
			fmt.Fprintf(cmd.OutOrStdout(), "claimed=%d passed=%d partial=%d failed=%d persist_errors=%d\n",
				report.Claimed, report.Passed, report.Partial, report.Failed, report.PersistErrors)
			return err
		}
		return worker.Run(cmd.Context(), limit, 5*time.Second)
	}}
	command.Flags().BoolVar(&once, "once", false, "process one batch and exit")
	command.Flags().IntVar(&limit, "limit", 40, "jobs claimed per batch (1-1000)")
	command.Flags().IntVar(&concurrency, "concurrency", 20, "maximum parallel validations")
	command.Flags().StringVar(&validatorID, "validator-id", "ai-a1", "stable validator identity")
	return command
}

func newValidationStatusCmd() *cobra.Command {
	return &cobra.Command{Use: "validation-status", Short: "Show bounded validation counters", RunE: func(cmd *cobra.Command, _ []string) error {
		_, db, _, err := openIngestService(cmd.Context())
		if err != nil {
			return err
		}
		defer db.Close()
		status, err := store.ReadValidationStatus(cmd.Context(), db)
		if err != nil {
			return err
		}
		fmt.Fprint(cmd.OutOrStdout(), formatValidationStatus(status))
		return nil
	}}
}

func formatValidationStatus(status store.ValidationStatus) string {
	return fmt.Sprintf("batches=%d attempts=%d current_statuses=%d pending_jobs=%d eligible_pending_jobs=%d leased_jobs=%d expired_leases=%d oldest_pending_age_seconds=%d\n"+
		"passed=%d partial=%d failed=%d available=%d available_24h=%d degraded=%d unavailable=%d\n"+
		"performance_attempts=%d performance_successes=%d average_bytes_per_second=%d\n"+
		"scored_nodes=%d average_quality_score=%d grade_S=%d grade_A=%d grade_B=%d grade_C=%d grade_D=%d grade_U=%d\n",
		status.Batches, status.Attempts, status.CurrentStatuses, status.PendingJobs, status.EligiblePendingJobs,
		status.LeasedJobs, status.ExpiredLeases, status.OldestPendingAgeSeconds,
		status.Passed, status.Partial, status.Failed, status.Available, status.Available24H, status.Degraded, status.Unavailable,
		status.PerformanceAttempts, status.PerformanceSuccesses, status.AverageBytesPerSecond,
		status.ScoredNodes, status.AverageQualityScore, status.ByGrade["S"], status.ByGrade["A"],
		status.ByGrade["B"], status.ByGrade["C"], status.ByGrade["D"], status.ByGrade["U"])
}
