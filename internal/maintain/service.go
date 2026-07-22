package maintain

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

type Service struct{ DB *sql.DB }
type Report struct {
	BeforeBytes, AfterBytes uint64
	Rows                    map[string]uint64
	Policy                  Policy
	DryRun                  bool
	StorageMetricRows       int64
}

func (s Service) Run(ctx context.Context, batchSize int, dryRun bool, now time.Time) (Report, error) {
	bytes, err := store.ReadStorageBytes(ctx, s.DB)
	if err != nil {
		return Report{}, err
	}
	policy := PolicyForUsage(float64(bytes) * 100 / float64(store.DatabaseCapacityBytes))
	var rollupDates []time.Time
	if !dryRun {
		cutoff := now.Add(-time.Duration(policy.AttemptDays) * 24 * time.Hour)
		rollupDates, err = store.ValidationAttemptRollupDatesBefore(ctx, s.DB, cutoff)
		if err != nil {
			return Report{}, fmt.Errorf("list expiring validation rollups: %w", err)
		}
	}
	if err := runRequiredRollups(dryRun, rollupDates, func(date time.Time) error {
		_, rollupErr := store.FinalizeDailyStats(ctx, s.DB, date, now)
		return rollupErr
	}); err != nil {
		return Report{}, fmt.Errorf("finalize expiring validation rollup: %w", err)
	}
	stored, err := store.RunMaintenance(ctx, s.DB, store.MaintenancePolicy{RawPayloadDays: policy.RawPayloadDays, ParseErrorDays: policy.ParseErrorDays, AttemptDays: policy.AttemptDays, BatchDays: policy.BatchDays, FetchDays: policy.FetchDays, ExportDays: policy.ExportDays, PauseColdSources: policy.PauseColdSources}, now, batchSize, dryRun)
	if err != nil {
		return Report{}, err
	}
	var storageMetricRows int64
	if !dryRun {
		storageMetricRows, err = store.RecordStorageMetrics(
			ctx, s.DB, now, store.DatabaseCapacityBytes, stored.AfterBytes,
		)
		if err != nil {
			return Report{}, err
		}
	}
	return Report{
		BeforeBytes:       stored.BeforeBytes,
		AfterBytes:        stored.AfterBytes,
		Rows:              stored.Rows,
		Policy:            policy,
		DryRun:            dryRun,
		StorageMetricRows: storageMetricRows,
	}, nil
}

func runRequiredRollups(dryRun bool, dates []time.Time, rollup func(time.Time) error) error {
	if dryRun {
		return nil
	}
	for _, date := range dates {
		if err := rollup(date); err != nil {
			return err
		}
	}
	return nil
}
