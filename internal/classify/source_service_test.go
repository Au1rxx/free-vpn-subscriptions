package classify

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

func TestRefreshSourceQualityMapsCandidatesToUpdates(t *testing.T) {
	now := time.Date(2026, 7, 14, 6, 0, 0, 0, time.UTC)
	updates, report := sourceQualityUpdates([]store.SourceQualityCandidate{
		{SourceID: 11, LastSuccessAt: sql.NullTime{Time: now, Valid: true}, FetchReliability: 1, ParseYield: 1, UsableNodeRate: 1},
		{SourceID: 12},
	}, now)
	if report.Candidates != 2 || report.Scored != 1 || len(updates) != 2 {
		t.Fatalf("report=%+v updates=%+v", report, updates)
	}
	if updates[0].SourceID != 11 || updates[0].Score != 100 {
		t.Fatalf("first update=%+v", updates[0])
	}
	if updates[1].SourceID != 12 || updates[1].Score != 0 {
		t.Fatalf("second update=%+v", updates[1])
	}
}
