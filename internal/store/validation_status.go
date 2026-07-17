package store

import (
	"context"
	"database/sql"
	"fmt"
)

type ValidationStatus struct {
	Batches, Attempts, CurrentStatuses, PendingJobs, EligiblePendingJobs    uint64
	LeasedJobs, ExpiredLeases, OldestPendingAgeSeconds                      uint64
	Passed, Partial, Failed, Available, Available24H, Degraded, Unavailable uint64
	PerformanceAttempts, PerformanceSuccesses, AverageBytesPerSecond        uint64
	ScoredNodes, AverageQualityScore                                        uint64
	ByGrade                                                                 map[string]uint64
}

var validationStatusQueries = []string{
	`SELECT COUNT(*) FROM validation_batches`,
	`SELECT
		COALESCE(SUM(job_state='pending'),0),
		COALESCE(SUM(job_state='pending' AND next_attempt_at <= UTC_TIMESTAMP(6)),0),
		COALESCE(SUM(job_state='leased'),0),
		COALESCE(SUM(job_state='leased' AND leased_until <= UTC_TIMESTAMP(6)),0),
		CAST(COALESCE(TIMESTAMPDIFF(SECOND,
			MIN(CASE WHEN job_state='pending' AND next_attempt_at <= UTC_TIMESTAMP(6) THEN next_attempt_at END),
			UTC_TIMESTAMP(6)),0) AS UNSIGNED)
	FROM validation_queue`,
	`SELECT
		COUNT(*),
		COALESCE(SUM(passed=TRUE),0),
		COALESCE(SUM(partial_success=TRUE),0),
		COALESCE(SUM(passed=FALSE AND partial_success=FALSE),0),
		COALESCE(SUM(performance_bytes IS NOT NULL),0),
		COALESCE(SUM(performance_bytes > 0 AND performance_error_code IS NULL),0),
		CAST(COALESCE(AVG(CASE WHEN performance_error_code IS NULL THEN NULLIF(bytes_per_second,0) END),0) AS UNSIGNED)
	FROM validation_attempts`,
	`SELECT
		COUNT(*),
		COALESCE(SUM(availability_state='available'),0),
		COALESCE(SUM(availability_state='available'
			AND last_validation_at >= UTC_TIMESTAMP(6) - INTERVAL 24 HOUR),0),
		COALESCE(SUM(availability_state='degraded'),0),
		COALESCE(SUM(availability_state='unavailable'),0),
		COALESCE(SUM(last_validation_at IS NOT NULL),0),
		CAST(COALESCE(AVG(CASE WHEN last_validation_at IS NOT NULL THEN quality_score END),0) AS UNSIGNED),
		COALESCE(SUM(quality_grade='S'),0),
		COALESCE(SUM(quality_grade='A'),0),
		COALESCE(SUM(quality_grade='B'),0),
		COALESCE(SUM(quality_grade='C'),0),
		COALESCE(SUM(quality_grade='D'),0),
		COALESCE(SUM(quality_grade='U'),0)
	FROM node_current_status`,
}

func ReadValidationStatus(ctx context.Context, db *sql.DB) (ValidationStatus, error) {
	status := ValidationStatus{ByGrade: make(map[string]uint64)}
	if err := db.QueryRowContext(ctx, validationStatusQueries[0]).Scan(&status.Batches); err != nil {
		return ValidationStatus{}, fmt.Errorf("read validation batch status: %w", err)
	}
	if err := db.QueryRowContext(ctx, validationStatusQueries[1]).Scan(
		&status.PendingJobs, &status.EligiblePendingJobs, &status.LeasedJobs, &status.ExpiredLeases,
		&status.OldestPendingAgeSeconds); err != nil {
		return ValidationStatus{}, fmt.Errorf("read validation queue status: %w", err)
	}
	if err := db.QueryRowContext(ctx, validationStatusQueries[2]).Scan(
		&status.Attempts, &status.Passed, &status.Partial, &status.Failed,
		&status.PerformanceAttempts, &status.PerformanceSuccesses, &status.AverageBytesPerSecond); err != nil {
		return ValidationStatus{}, fmt.Errorf("read validation attempt status: %w", err)
	}
	var gradeS, gradeA, gradeB, gradeC, gradeD, gradeU uint64
	if err := db.QueryRowContext(ctx, validationStatusQueries[3]).Scan(
		&status.CurrentStatuses, &status.Available, &status.Available24H, &status.Degraded, &status.Unavailable,
		&status.ScoredNodes, &status.AverageQualityScore,
		&gradeS, &gradeA, &gradeB, &gradeC, &gradeD, &gradeU); err != nil {
		return ValidationStatus{}, fmt.Errorf("read current validation status: %w", err)
	}
	status.ByGrade["S"] = gradeS
	status.ByGrade["A"] = gradeA
	status.ByGrade["B"] = gradeB
	status.ByGrade["C"] = gradeC
	status.ByGrade["D"] = gradeD
	status.ByGrade["U"] = gradeU
	return status, nil
}
