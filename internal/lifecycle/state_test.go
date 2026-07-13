package lifecycle

import (
	"testing"
	"time"
)

func TestLifecycleDecisionsAndActiveSourceProtection(t *testing.T) {
	now := time.Date(2026, 7, 13, 0, 0, 0, 0, time.UTC)
	recent := now.Add(-time.Hour)
	oldSuccess := now.Add(-100 * 24 * time.Hour)
	oldFirst := now.Add(-40 * 24 * time.Hour)
	oldArchive := now.Add(-100 * 24 * time.Hour)
	tests := []struct {
		input Input
		state string
	}{
		{Input{Now: now, FirstSeenAt: now, ActiveSources: 1}, "pending"},
		{Input{Now: now, LastSuccessAt: &recent, ActiveSources: 1}, "active"},
		{Input{Now: now, LastSuccessAt: &recent, ActiveSources: 1, ConsecutiveFailures: 1}, "degraded"},
		{Input{Now: now, LastSuccessAt: &oldSuccess, ActiveSources: 1}, "stale"},
		{Input{Now: now, FirstSeenAt: oldFirst, ActiveSources: 0}, "dead"},
		{Input{Now: now, State: "archived", ArchivedAt: &oldArchive, ActiveSources: 0}, "purged"},
		{Input{Now: now, State: "archived", ArchivedAt: &oldArchive, ActiveSources: 1}, "pending"},
	}
	for _, tt := range tests {
		if decision := Decide(tt.input); decision.State != tt.state {
			t.Fatalf("input=%+v decision=%+v want=%s", tt.input, decision, tt.state)
		}
	}
}
