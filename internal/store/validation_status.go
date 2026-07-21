package store

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
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
	FROM validation_attempts FORCE INDEX (idx_validation_attempts_status)`,
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
	var batches uint64
	var queue struct {
		pending, eligible, leased, expired, oldest uint64
	}
	var attempts struct {
		total, passed, partial, failed, performance, performanceSuccesses, averageBytesPerSecond uint64
	}
	var current struct {
		total, available, available24H, degraded, unavailable, scored, averageScore uint64
	}
	var grades [6]uint64
	queries := []func() error{
		func() error {
			return db.QueryRowContext(ctx, validationStatusQueries[0]).Scan(&batches)
		},
		func() error {
			return db.QueryRowContext(ctx, validationStatusQueries[1]).Scan(
				&queue.pending, &queue.eligible, &queue.leased, &queue.expired, &queue.oldest)
		},
		func() error {
			return db.QueryRowContext(ctx, validationStatusQueries[2]).Scan(
				&attempts.total, &attempts.passed, &attempts.partial, &attempts.failed,
				&attempts.performance, &attempts.performanceSuccesses, &attempts.averageBytesPerSecond)
		},
		func() error {
			return db.QueryRowContext(ctx, validationStatusQueries[3]).Scan(
				&current.total, &current.available, &current.available24H, &current.degraded, &current.unavailable,
				&current.scored, &current.averageScore,
				&grades[0], &grades[1], &grades[2], &grades[3], &grades[4], &grades[5])
		},
	}
	queryNames := []string{"validation batch status", "validation queue status", "validation attempt status", "current validation status"}
	queryErrors := make([]error, len(queries))
	var wait sync.WaitGroup
	wait.Add(len(queries))
	for index, query := range queries {
		go func() {
			defer wait.Done()
			queryErrors[index] = query()
		}()
	}
	wait.Wait()
	for index, err := range queryErrors {
		if err != nil {
			return ValidationStatus{}, fmt.Errorf("read %s: %w", queryNames[index], err)
		}
	}
	status.Batches = batches
	status.PendingJobs, status.EligiblePendingJobs = queue.pending, queue.eligible
	status.LeasedJobs, status.ExpiredLeases, status.OldestPendingAgeSeconds = queue.leased, queue.expired, queue.oldest
	status.Attempts, status.Passed, status.Partial, status.Failed = attempts.total, attempts.passed, attempts.partial, attempts.failed
	status.PerformanceAttempts, status.PerformanceSuccesses = attempts.performance, attempts.performanceSuccesses
	status.AverageBytesPerSecond = attempts.averageBytesPerSecond
	status.CurrentStatuses, status.Available, status.Available24H = current.total, current.available, current.available24H
	status.Degraded, status.Unavailable = current.degraded, current.unavailable
	status.ScoredNodes, status.AverageQualityScore = current.scored, current.averageScore
	for index, grade := range []string{"S", "A", "B", "C", "D", "U"} {
		status.ByGrade[grade] = grades[index]
	}
	return status, nil
}
