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
}

func ReadIngestStatus(ctx context.Context, db *sql.DB) (IngestStatus, error) {
	status := IngestStatus{ByProtocol: make(map[string]uint64)}
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
	return status, rows.Err()
}
