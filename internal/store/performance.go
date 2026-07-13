package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func PerformanceSampleDue(ctx context.Context, db *sql.DB, nodeConfigID uint64, now time.Time) (bool, error) {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM validation_attempts
		WHERE node_config_id=? AND performance_bytes IS NOT NULL AND started_at >= ?`,
		nodeConfigID, now.UTC().Add(-24*time.Hour)).Scan(&count); err != nil {
		return false, fmt.Errorf("read performance sampling window: %w", err)
	}
	return count == 0, nil
}

func RecordPerformanceSample(ctx context.Context, db *sql.DB, attemptID uint64, bytes, bytesPerSecond int64) error {
	result, err := db.ExecContext(ctx, `UPDATE validation_attempts SET performance_bytes=?, bytes_per_second=?
		WHERE validation_attempt_id=? AND performance_bytes IS NULL`, bytes, bytesPerSecond, attemptID)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected != 1 {
		return fmt.Errorf("performance sample was already recorded or attempt is missing")
	}
	return nil
}
