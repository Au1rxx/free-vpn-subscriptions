package store

import (
	"context"
	"database/sql"
	"fmt"
)

type IngestStatus struct {
	Sources, EnabledSources, Fetches, PendingFetches, ParseRuns, Endpoints, Configs, ParseErrors, QueuePending uint64
	Fetches24H, SuccessfulFetches24H, FailedFetches24H                                                         uint64
	ByProtocol                                                                                                 map[string]uint64
	SourceKinds                                                                                                map[string]SourceKindStatus
	FetchErrorCounts24H, ParseErrorCounts                                                                      []NamedCount
}

type SourceKindStatus struct{ Total, Enabled uint64 }
type NamedCount struct {
	Name  string
	Count uint64
}

func ReadIngestStatus(ctx context.Context, db *sql.DB) (IngestStatus, error) {
	status := IngestStatus{ByProtocol: make(map[string]uint64), SourceKinds: make(map[string]SourceKindStatus)}
	err := db.QueryRowContext(ctx, `SELECT
		(SELECT COUNT(*) FROM sources),
		(SELECT COUNT(*) FROM sources WHERE enabled=TRUE AND state='active'),
		(SELECT COUNT(*) FROM source_fetches),
		(SELECT COUNT(*) FROM source_fetches WHERE parse_state='pending'),
		(SELECT COUNT(*) FROM parse_runs),
		(SELECT COUNT(*) FROM endpoints),
		(SELECT COUNT(*) FROM node_configs),
		(SELECT COUNT(*) FROM parse_errors),
		(SELECT COUNT(*) FROM validation_queue WHERE job_state='pending'),
		(SELECT COUNT(*) FROM source_fetches WHERE finished_at >= UTC_TIMESTAMP(6) - INTERVAL 24 HOUR),
		(SELECT COUNT(*) FROM source_fetches WHERE finished_at >= UTC_TIMESTAMP(6) - INTERVAL 24 HOUR
			AND fetch_state IN ('success','not_modified')),
		(SELECT COUNT(*) FROM source_fetches WHERE finished_at >= UTC_TIMESTAMP(6) - INTERVAL 24 HOUR
			AND fetch_state NOT IN ('success','not_modified'))`).Scan(
		&status.Sources, &status.EnabledSources, &status.Fetches, &status.PendingFetches, &status.ParseRuns,
		&status.Endpoints, &status.Configs, &status.ParseErrors, &status.QueuePending,
		&status.Fetches24H, &status.SuccessfulFetches24H, &status.FailedFetches24H)
	if err != nil {
		return IngestStatus{}, fmt.Errorf("read ingest counters: %w", err)
	}
	rows, err := db.QueryContext(ctx, `SELECT protocol, COUNT(*) FROM node_configs GROUP BY protocol ORDER BY protocol`)
	if err != nil {
		return IngestStatus{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var protocol string
		var count uint64
		if err := rows.Scan(&protocol, &count); err != nil {
			return IngestStatus{}, err
		}
		status.ByProtocol[protocol] = count
	}
	if err := rows.Close(); err != nil {
		return IngestStatus{}, err
	}
	if err := rows.Err(); err != nil {
		return IngestStatus{}, err
	}
	rows, err = db.QueryContext(ctx, `SELECT COALESCE(NULLIF(kind,''),'unknown'), COUNT(*),
		SUM(enabled=TRUE AND state='active') FROM sources GROUP BY kind ORDER BY kind`)
	if err != nil {
		return IngestStatus{}, fmt.Errorf("read source kind distribution: %w", err)
	}
	for rows.Next() {
		var kind string
		var item SourceKindStatus
		if err := rows.Scan(&kind, &item.Total, &item.Enabled); err != nil {
			rows.Close()
			return IngestStatus{}, err
		}
		status.SourceKinds[kind] = item
	}
	if err := rows.Close(); err != nil {
		return IngestStatus{}, err
	}
	if err := rows.Err(); err != nil {
		return IngestStatus{}, err
	}
	rows, err = db.QueryContext(ctx, `SELECT CONCAT(COALESCE(error_code,'unknown'),
		IF(http_status IS NULL,'',CONCAT('_http_',http_status))), COUNT(*) AS occurrences
		FROM source_fetches WHERE fetch_state='failed'
		AND finished_at >= UTC_TIMESTAMP(6) - INTERVAL 24 HOUR
		GROUP BY error_code, http_status ORDER BY occurrences DESC, error_code LIMIT 20`)
	if err != nil {
		return IngestStatus{}, fmt.Errorf("read fetch error distribution: %w", err)
	}
	status.FetchErrorCounts24H, err = scanNamedCounts(rows)
	if err != nil {
		return IngestStatus{}, err
	}
	rows, err = db.QueryContext(ctx, `SELECT error_code, SUM(seen_count) AS occurrences
		FROM parse_errors GROUP BY error_code ORDER BY occurrences DESC, error_code LIMIT 20`)
	if err != nil {
		return IngestStatus{}, fmt.Errorf("read parse error distribution: %w", err)
	}
	status.ParseErrorCounts, err = scanNamedCounts(rows)
	return status, err
}

func scanNamedCounts(rows *sql.Rows) ([]NamedCount, error) {
	defer rows.Close()
	var counts []NamedCount
	for rows.Next() {
		var item NamedCount
		if err := rows.Scan(&item.Name, &item.Count); err != nil {
			return nil, err
		}
		counts = append(counts, item)
	}
	return counts, rows.Err()
}
