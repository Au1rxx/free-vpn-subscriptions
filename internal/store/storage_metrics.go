package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// RecordStorageMetrics persists one idempotent per-table capacity snapshot for
// the current database schema.
func RecordStorageMetrics(
	ctx context.Context,
	db *sql.DB,
	sampledAt time.Time,
	capacityBytes uint64,
	usageBytes uint64,
) (int64, error) {
	if capacityBytes == 0 {
		return 0, fmt.Errorf("storage metric capacity must be greater than zero")
	}
	result, err := db.ExecContext(ctx, `
		INSERT INTO storage_metrics (
			sampled_at, table_schema, table_name, table_rows_estimate,
			data_bytes, index_bytes, total_bytes, capacity_bytes, usage_percent
		)
		SELECT ?, table_schema, table_name, COALESCE(table_rows, 0),
			COALESCE(data_length, 0), COALESCE(index_length, 0),
			COALESCE(data_length, 0) + COALESCE(index_length, 0),
			?, ROUND((? * 100.0) / ?, 3)
		FROM information_schema.tables
		WHERE table_schema = DATABASE()
		ON DUPLICATE KEY UPDATE
			table_rows_estimate=VALUES(table_rows_estimate),
			data_bytes=VALUES(data_bytes),
			index_bytes=VALUES(index_bytes),
			total_bytes=VALUES(total_bytes),
			capacity_bytes=VALUES(capacity_bytes),
			usage_percent=VALUES(usage_percent)`,
		sampledAt.UTC(), capacityBytes, usageBytes, capacityBytes)
	if err != nil {
		return 0, fmt.Errorf("record storage metrics: %w", err)
	}
	written, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("read storage metric affected rows: %w", err)
	}
	return written, nil
}
