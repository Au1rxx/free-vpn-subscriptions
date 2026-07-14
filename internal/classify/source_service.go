package classify

import (
	"context"
	"database/sql"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

// SourceReport summarizes one complete source-quality refresh.
type SourceReport struct {
	Candidates int
	Scored     int
	Written    int64
}

// RefreshSourceQualities derives and atomically persists every source score.
func RefreshSourceQualities(ctx context.Context, db *sql.DB, now time.Time) (SourceReport, error) {
	candidates, err := store.ListSourceQualityCandidates(ctx, db, now)
	if err != nil {
		return SourceReport{}, err
	}
	updates, report := sourceQualityUpdates(candidates, now)
	written, err := store.WriteSourceQualities(ctx, db, updates)
	if err != nil {
		return SourceReport{}, err
	}
	report.Written = written
	return report, nil
}

func sourceQualityUpdates(candidates []store.SourceQualityCandidate, now time.Time) ([]store.SourceQualityUpdate, SourceReport) {
	updates := make([]store.SourceQualityUpdate, 0, len(candidates))
	report := SourceReport{Candidates: len(candidates)}
	for _, candidate := range candidates {
		freshnessHours := -1.0
		if candidate.LastSuccessAt.Valid {
			freshnessHours = now.Sub(candidate.LastSuccessAt.Time).Hours()
			if freshnessHours < 0 {
				freshnessHours = 0
			}
		}
		breakdown := ScoreSource(SourceInput{
			FetchReliability: candidate.FetchReliability,
			ParseYield:       candidate.ParseYield,
			UsableNodeRate:   candidate.UsableNodeRate,
			FreshnessHours:   freshnessHours,
		})
		if breakdown.Total > 0 {
			report.Scored++
		}
		updates = append(updates, store.SourceQualityUpdate{SourceID: candidate.SourceID, Score: breakdown.Total})
	}
	return updates, report
}
