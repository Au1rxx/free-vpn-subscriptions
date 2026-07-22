package store

import (
	"strings"
	"testing"
)

func TestMaintenanceTargetsDeleteAttemptsBeforeUnreferencedBatches(t *testing.T) {
	policy := MaintenancePolicy{
		RawPayloadDays: 30, ParseErrorDays: 30, AttemptDays: 14,
		BatchDays: 14, FetchDays: 90, ExportDays: 30,
	}
	targets := maintenanceTargets(policy)
	wantNames := []string{"raw_payloads", "parse_errors", "validation_attempts", "validation_batches", "source_fetches", "export_members"}
	if len(targets) != len(wantNames) {
		t.Fatalf("targets=%d want=%d", len(targets), len(wantNames))
	}
	for index, want := range wantNames {
		if targets[index].name != want {
			t.Fatalf("target[%d]=%q want=%q", index, targets[index].name, want)
		}
	}
	batch := targets[3]
	if batch.days != 14 || !strings.Contains(batch.countSQL, "NOT EXISTS") ||
		!strings.Contains(batch.deleteSQL, "NOT EXISTS") ||
		!strings.Contains(batch.deleteSQL, "LIMIT ?") {
		t.Fatalf("unsafe validation batch cleanup: %#v", batch)
	}
}

func TestExpiringAttemptRollupQueriesAreIndexedAndFinalizationSafe(t *testing.T) {
	for _, want := range []string{"FORCE INDEX (idx_validation_attempts_cleanup)", "finished_at < ?", "DISTINCT DATE(started_at)"} {
		if !strings.Contains(validationAttemptRollupDatesSQL, want) {
			t.Fatalf("rollup date query missing %q: %s", want, validationAttemptRollupDatesSQL)
		}
	}
	for _, want := range []string{"finalized_at", "IF(finalized_at IS NULL", "COALESCE(finalized_at"} {
		if !strings.Contains(rollupDailyStatsSQL, want) {
			t.Fatalf("daily rollup is not retry-safe; missing %q", want)
		}
	}
}
