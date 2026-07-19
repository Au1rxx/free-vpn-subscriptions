package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/verify"
)

var ErrStaleValidation = errors.New("validation result is older than current status")

type ValidationWrite struct {
	JobID, NodeConfigID       uint64
	Owner, ValidatorID, Stage string
	Engine, EngineVersion     string
	StartedAt, FinishedAt     time.Time
	Round                     int
	Result                    verify.Result
}

func PersistValidationResult(ctx context.Context, db *sql.DB, write ValidationWrite) error {
	if db == nil || write.JobID == 0 || write.NodeConfigID == 0 || write.Owner == "" || write.ValidatorID == "" {
		return fmt.Errorf("validation job, node, owner and validator identity are required")
	}
	if write.Stage == "" {
		write.Stage = "connectivity"
	}
	if write.Engine == "" {
		write.Engine = "sing-box"
	}
	if write.Round < 1 {
		write.Round = 1
	}
	if write.StartedAt.IsZero() {
		write.StartedAt = time.Now().UTC()
	}
	if write.FinishedAt.IsZero() {
		write.FinishedAt = time.Now().UTC()
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var leasedNode uint64
	var leaseOwner, jobState string
	if err := tx.QueryRowContext(ctx, `SELECT node_config_id, COALESCE(lease_owner,''), job_state
		FROM validation_queue WHERE validation_job_id=? FOR UPDATE`, write.JobID).Scan(&leasedNode, &leaseOwner, &jobState); err != nil {
		return err
	}
	if leasedNode != write.NodeConfigID || leaseOwner != write.Owner || jobState != "leased" {
		return ErrLeaseOwnership
	}
	var currentValidation sql.NullTime
	if err := tx.QueryRowContext(ctx, `SELECT last_validation_at FROM node_current_status WHERE node_config_id=?`, write.NodeConfigID).Scan(&currentValidation); err != nil && err != sql.ErrNoRows {
		return err
	}
	if currentValidation.Valid && currentValidation.Time.After(write.FinishedAt) {
		return ErrStaleValidation
	}
	state, availability, quality := validationState(write.Result)
	var failures int
	if err := tx.QueryRowContext(ctx, `SELECT consecutive_failures FROM node_configs WHERE node_config_id=?`, write.NodeConfigID).Scan(&failures); err != nil {
		return err
	}
	if write.Result.Passed || write.Result.PartialSuccess {
		failures = 0
	} else {
		failures++
	}
	nextAttempt := nextValidationAttempt(state, failures, write.FinishedAt)
	batchResult, err := tx.ExecContext(ctx, `INSERT INTO validation_batches
		(validator_id, stage, engine, engine_version, config_snapshot, started_at, finished_at,
		 claimed_count, success_count, partial_count, failure_count, batch_state)
		VALUES (?, ?, ?, ?, JSON_OBJECT('single_job', TRUE), ?, ?, 1, ?, ?, ?, 'completed')`,
		write.ValidatorID, write.Stage, write.Engine, nullString(write.EngineVersion), write.StartedAt, write.FinishedAt,
		boolInt(write.Result.Passed), boolInt(write.Result.PartialSuccess), boolInt(!write.Result.Passed && !write.Result.PartialSuccess))
	if err != nil {
		return fmt.Errorf("insert validation batch: %w", err)
	}
	batchID, err := batchResult.LastInsertId()
	if err != nil {
		return err
	}
	targetJSON, err := json.Marshal(write.Result.Targets)
	if err != nil {
		return err
	}
	metrics := aggregateTargetMetrics(write.Result.Targets)
	_, err = tx.ExecContext(ctx, `INSERT INTO validation_attempts
		(validation_job_id, validation_batch_id, node_config_id, validator_id, stage, round_number,
		 started_at, finished_at, config_accepted, proxy_started, passed, partial_success,
		 target_count, success_count, dns_ms, connect_ms, tls_ms, proxy_start_ms, ttfb_ms,
		 total_ms, http_median_ms, performance_bytes, bytes_per_second, performance_error_code,
		 exit_ip, exit_country, exit_asn, target_results, error_code, error_summary)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		write.JobID, batchID, write.NodeConfigID, write.ValidatorID, write.Stage, write.Round,
		write.StartedAt, write.FinishedAt, write.Result.ConfigAccepted, write.Result.ProxyStarted,
		write.Result.Passed, write.Result.PartialSuccess, write.Result.Attempts, write.Result.Successes,
		nullInt(metrics.DNSMS), nullInt(metrics.ConnectMS), nullInt(metrics.TLSMS), nullInt(write.Result.StartMS),
		nullInt(metrics.TTFBMS), nullInt(metrics.TotalMS), nullInt(write.Result.HTTPMedianMS),
		performanceValue(write.Result.Performance.Attempted, write.Result.Performance.Bytes),
		performanceValue(write.Result.Performance.Attempted, write.Result.Performance.BytesPerSecond),
		nullString(write.Result.Performance.ErrorCode),
		nullString(write.Result.ExitIP), nullString(write.Result.ExitCountry), nullString(write.Result.ExitASN),
		targetJSON, nullString(write.Result.ErrorCode), nullString(bounded(write.Result.ErrorSummary, 1024)))
	if err != nil {
		return fmt.Errorf("insert validation attempt: %w", err)
	}
	lastSuccess := any(nil)
	if write.Result.Passed || write.Result.PartialSuccess {
		lastSuccess = write.FinishedAt
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO node_current_status
		(node_config_id, lifecycle_state, availability_state, quality_score, quality_grade,
		 last_validation_at, last_success_at, last_failure_at, latency_p50_ms, source_count,
		 consecutive_successes, consecutive_failures, exit_ip, exit_country, exit_asn,
		 last_error_code, last_error_summary, next_check_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,
		 (SELECT COUNT(*) FROM node_source_stats WHERE node_config_id=? AND is_active=TRUE),
		 ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE lifecycle_state=VALUES(lifecycle_state), availability_state=VALUES(availability_state),
		 quality_score=VALUES(quality_score), quality_grade=VALUES(quality_grade),
		 last_validation_at=VALUES(last_validation_at),
		 last_success_at=COALESCE(VALUES(last_success_at), last_success_at),
		 last_failure_at=COALESCE(VALUES(last_failure_at), last_failure_at),
		 latency_p50_ms=VALUES(latency_p50_ms), source_count=VALUES(source_count),
		 consecutive_successes=VALUES(consecutive_successes), consecutive_failures=VALUES(consecutive_failures),
		 exit_ip=COALESCE(VALUES(exit_ip), exit_ip), exit_country=COALESCE(VALUES(exit_country), exit_country),
		 exit_asn=COALESCE(VALUES(exit_asn), exit_asn), last_error_code=VALUES(last_error_code),
		 last_error_summary=VALUES(last_error_summary), next_check_at=VALUES(next_check_at)`,
		write.NodeConfigID, state, availability, quality, qualityGrade(quality), write.FinishedAt,
		lastSuccess, failureTime(write.Result, write.FinishedAt), nullInt(write.Result.HTTPMedianMS), write.NodeConfigID,
		boolInt(write.Result.Passed || write.Result.PartialSuccess), failures,
		nullString(write.Result.ExitIP), nullString(write.Result.ExitCountry), nullString(write.Result.ExitASN),
		nullString(write.Result.ErrorCode), nullString(bounded(write.Result.ErrorSummary, 1024)), nextAttempt)
	if err != nil {
		return fmt.Errorf("update current validation status: %w", err)
	}
	_, err = tx.ExecContext(ctx, `UPDATE node_configs SET lifecycle_state=?, is_exportable=?,
		 last_success_at=COALESCE(?, last_success_at), consecutive_failures=?, updated_at=UTC_TIMESTAMP(6)
		 WHERE node_config_id=?`, state, write.Result.Passed, lastSuccess, failures, write.NodeConfigID)
	if err != nil {
		return err
	}
	queueResult, err := tx.ExecContext(ctx, `UPDATE validation_queue SET job_state='pending', next_attempt_at=?,
		 lease_owner=NULL, leased_until=NULL, completed_at=?, last_error_code=?, last_error_summary=?
		 WHERE validation_job_id=? AND lease_owner=? AND job_state='leased'`, nextAttempt, write.FinishedAt,
		nullString(write.Result.ErrorCode), nullString(bounded(write.Result.ErrorSummary, 1024)), write.JobID, write.Owner)
	if err != nil {
		return err
	}
	if affected, _ := queueResult.RowsAffected(); affected != 1 {
		return ErrLeaseOwnership
	}
	return tx.Commit()
}

func performanceValue(attempted bool, value int64) any {
	if !attempted {
		return nil
	}
	return value
}

type targetMetrics struct{ DNSMS, ConnectMS, TLSMS, TTFBMS, TotalMS int }

func aggregateTargetMetrics(targets []verify.TargetResult) targetMetrics {
	var result targetMetrics
	for _, target := range targets {
		result.DNSMS = max(result.DNSMS, target.DNSMS)
		result.ConnectMS = max(result.ConnectMS, target.ConnectMS)
		result.TLSMS = max(result.TLSMS, target.TLSMS)
		result.TTFBMS = max(result.TTFBMS, target.TTFBMS)
		result.TotalMS = max(result.TotalMS, target.TotalMS)
	}
	return result
}

func validationState(result verify.Result) (string, string, int) {
	if result.Passed {
		return "active", "available", 85
	}
	if result.PartialSuccess {
		return "active", "degraded", 55
	}
	return "pending", "unavailable", 10
}

func qualityGrade(score int) string {
	switch {
	case score >= 80:
		return "A"
	case score >= 60:
		return "B"
	case score >= 40:
		return "C"
	default:
		return "D"
	}
}

func nextValidationAttempt(state string, failures int, now time.Time) time.Time {
	if state == "active" && failures == 0 {
		return now.Add(6 * time.Hour)
	}
	intervals := []time.Duration{15 * time.Minute, time.Hour, 6 * time.Hour, 24 * time.Hour, 3 * 24 * time.Hour, 7 * 24 * time.Hour}
	if failures < 1 {
		failures = 1
	}
	if failures > len(intervals) {
		failures = len(intervals)
	}
	return now.Add(intervals[failures-1])
}

func failureTime(result verify.Result, finished time.Time) any {
	if result.Passed || result.PartialSuccess {
		return nil
	}
	return finished
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
