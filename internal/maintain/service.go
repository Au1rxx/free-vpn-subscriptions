package maintain

import (
	"context"
	"database/sql"
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
	if !dryRun {
		_, _ = store.RollupDailyStats(ctx, s.DB, now)
	}
	stored, err := store.RunMaintenance(ctx, s.DB, store.MaintenancePolicy{RawPayloadDays: policy.RawPayloadDays, ParseErrorDays: policy.ParseErrorDays, AttemptDays: policy.AttemptDays, FetchDays: policy.FetchDays, ExportDays: policy.ExportDays, PauseColdSources: policy.PauseColdSources, StoreRawBodies: policy.StoreRawBodies}, now, batchSize, dryRun)
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
