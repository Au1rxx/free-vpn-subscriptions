package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// SourceQualityCandidate contains normalized inputs for one source score.
type SourceQualityCandidate struct {
	SourceID         uint64
	LastSuccessAt    sql.NullTime
	FetchReliability float64
	ParseYield       float64
	UsableNodeRate   float64
}

// SourceQualityUpdate is one derived source score ready for persistence.
type SourceQualityUpdate struct {
	SourceID uint64
	Score    int
}

// ListSourceQualityCandidates aggregates bounded recent source health and the
// current terminal validation state without loading raw payloads.
func ListSourceQualityCandidates(ctx context.Context, db *sql.DB, now time.Time) ([]SourceQualityCandidate, error) {
	rows, err := db.QueryContext(ctx, `SELECT s.source_id, s.last_success_at,
		COALESCE(f.reliability,0), COALESCE(p.parse_yield,0), COALESCE(v.usable_rate,0)
		FROM sources s
		LEFT JOIN (
			SELECT source_id, AVG(fetch_state IN ('success','not_modified')) reliability
			FROM source_fetches WHERE finished_at >= ? GROUP BY source_id
		) f ON f.source_id=s.source_id
		LEFT JOIN (
			SELECT sf.source_id,
				SUM(pr.success_entries)/NULLIF(SUM(pr.success_entries+pr.error_entries),0) parse_yield
			FROM parse_runs pr JOIN source_fetches sf ON sf.fetch_id=pr.fetch_id
			WHERE pr.finished_at >= ? AND pr.parse_state='success' GROUP BY sf.source_id
		) p ON p.source_id=s.source_id
		LEFT JOIN (
			SELECT ns.source_id,
				AVG(CASE cs.availability_state WHEN 'available' THEN 1 WHEN 'degraded' THEN 0.5 ELSE 0 END) usable_rate
			FROM node_current_status cs JOIN node_source_stats ns ON ns.node_config_id=cs.node_config_id
			WHERE ns.is_active=TRUE AND cs.availability_state IN ('available','degraded','unavailable')
			GROUP BY ns.source_id
		) v ON v.source_id=s.source_id
		ORDER BY s.source_id`, now.Add(-24*time.Hour), now.Add(-7*24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("query source quality candidates: %w", err)
	}
	defer rows.Close()
	var candidates []SourceQualityCandidate
	for rows.Next() {
		var candidate SourceQualityCandidate
		if err := rows.Scan(&candidate.SourceID, &candidate.LastSuccessAt, &candidate.FetchReliability,
			&candidate.ParseYield, &candidate.UsableNodeRate); err != nil {
			return nil, err
		}
		candidates = append(candidates, candidate)
	}
	return candidates, rows.Err()
}

// WriteSourceQualities atomically updates all supplied scores in bounded SQL
// batches. A failed batch rolls back every score from the refresh.
func WriteSourceQualities(ctx context.Context, db *sql.DB, updates []SourceQualityUpdate) (int64, error) {
	if len(updates) == 0 {
		return 0, nil
	}
	if len(updates) > 10000 {
		return 0, fmt.Errorf("source quality updates must not exceed 10000")
	}
	for _, update := range updates {
		if update.SourceID == 0 || update.Score < 0 || update.Score > 100 {
			return 0, fmt.Errorf("invalid source quality update: source=%d score=%d", update.SourceID, update.Score)
		}
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var written int64
	for start := 0; start < len(updates); start += classificationWriteBatchSize {
		end := start + classificationWriteBatchSize
		if end > len(updates) {
			end = len(updates)
		}
		batch := updates[start:end]
		var query strings.Builder
		query.WriteString("UPDATE sources SET quality_score=CASE source_id ")
		args := make([]any, 0, len(batch)*3)
		for _, update := range batch {
			query.WriteString("WHEN ? THEN ? ")
			args = append(args, update.SourceID, update.Score)
		}
		query.WriteString("ELSE quality_score END WHERE source_id IN (")
		query.WriteString(scalarPlaceholders(len(batch)))
		query.WriteByte(')')
		for _, update := range batch {
			args = append(args, update.SourceID)
		}
		result, err := tx.ExecContext(ctx, query.String(), args...)
		if err != nil {
			return 0, fmt.Errorf("write source quality batch: %w", err)
		}
		changed, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		written += changed
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return written, nil
}
