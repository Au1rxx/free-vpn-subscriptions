package store

import (
	"context"
	"database/sql"
	"fmt"
)

type ValidationStatus struct {
	Batches, Attempts, CurrentStatuses, PendingJobs, LeasedJobs, ExpiredLeases uint64
	Passed, Partial, Failed, Available, Degraded, Unavailable                  uint64
	PerformanceAttempts, PerformanceSuccesses, AverageBytesPerSecond           uint64
}

func ReadValidationStatus(ctx context.Context, db *sql.DB) (ValidationStatus, error) {
	var status ValidationStatus
	err := db.QueryRowContext(ctx, `SELECT
		(SELECT COUNT(*) FROM validation_batches),
		(SELECT COUNT(*) FROM validation_attempts),
		(SELECT COUNT(*) FROM node_current_status),
		(SELECT COUNT(*) FROM validation_queue WHERE job_state='pending'),
		(SELECT COUNT(*) FROM validation_queue WHERE job_state='leased'),
		(SELECT COUNT(*) FROM validation_queue WHERE job_state='leased' AND leased_until <= UTC_TIMESTAMP(6)),
		(SELECT COUNT(*) FROM validation_attempts WHERE passed=TRUE),
		(SELECT COUNT(*) FROM validation_attempts WHERE partial_success=TRUE),
		(SELECT COUNT(*) FROM validation_attempts WHERE passed=FALSE AND partial_success=FALSE),
		(SELECT COUNT(*) FROM node_current_status WHERE availability_state='available'),
		(SELECT COUNT(*) FROM node_current_status WHERE availability_state='degraded'),
		(SELECT COUNT(*) FROM node_current_status WHERE availability_state='unavailable'),
		(SELECT COUNT(*) FROM validation_attempts WHERE performance_bytes IS NOT NULL),
		(SELECT COUNT(*) FROM validation_attempts WHERE performance_bytes > 0 AND performance_error_code IS NULL),
		(SELECT CAST(COALESCE(AVG(NULLIF(bytes_per_second,0)),0) AS UNSIGNED) FROM validation_attempts
			WHERE performance_error_code IS NULL)`).Scan(
		&status.Batches, &status.Attempts, &status.CurrentStatuses, &status.PendingJobs,
		&status.LeasedJobs, &status.ExpiredLeases, &status.Passed, &status.Partial, &status.Failed,
		&status.Available, &status.Degraded, &status.Unavailable,
		&status.PerformanceAttempts, &status.PerformanceSuccesses, &status.AverageBytesPerSecond)
	if err != nil {
		return ValidationStatus{}, fmt.Errorf("read validation status: %w", err)
	}
	return status, nil
}
