package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrLeaseOwnership = errors.New("validation job lease ownership mismatch")

type ValidationJob struct {
	ID, NodeConfigID           uint64
	Stage, Protocol, State     string
	Priority, Attempts         int
	LeaseOwner                 string
	LeasedUntil, NextAttemptAt time.Time
	NormalizedConfig           []byte
}

func ClaimValidationJobs(ctx context.Context, db *sql.DB, owner string, limit int, lease time.Duration) ([]ValidationJob, error) {
	if owner == "" {
		return nil, fmt.Errorf("validation lease owner is required")
	}
	if limit < 1 || limit > 1000 {
		return nil, fmt.Errorf("validation claim limit must be between 1 and 1000")
	}
	if lease <= 0 || lease > time.Hour {
		return nil, fmt.Errorf("validation lease must be between zero and one hour")
	}
	if db == nil {
		return nil, fmt.Errorf("validation database is required")
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, `SELECT q.validation_job_id, q.node_config_id, q.stage,
		q.priority, q.job_state, q.attempts, COALESCE(q.lease_owner,''),
		q.next_attempt_at, q.leased_until, n.protocol, n.normalized_config
		FROM validation_queue q JOIN node_configs n ON n.node_config_id=q.node_config_id
		WHERE (q.job_state='pending' AND q.next_attempt_at <= UTC_TIMESTAMP(6))
		   OR (q.job_state='leased' AND q.leased_until <= UTC_TIMESTAMP(6))
		ORDER BY q.priority DESC, q.next_attempt_at ASC, q.validation_job_id ASC
		LIMIT ? FOR UPDATE SKIP LOCKED`, limit)
	if err != nil {
		return nil, fmt.Errorf("select validation jobs: %w", err)
	}
	var jobs []ValidationJob
	for rows.Next() {
		var job ValidationJob
		var leasedUntil sql.NullTime
		if err := rows.Scan(&job.ID, &job.NodeConfigID, &job.Stage, &job.Priority, &job.State,
			&job.Attempts, &job.LeaseOwner, &job.NextAttemptAt, &leasedUntil,
			&job.Protocol, &job.NormalizedConfig); err != nil {
			rows.Close()
			return nil, err
		}
		if leasedUntil.Valid {
			job.LeasedUntil = leasedUntil.Time
		}
		jobs = append(jobs, job)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(jobs) == 0 {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	ids := make([]any, len(jobs)+2)
	ids[0], ids[1] = owner, time.Now().UTC().Add(lease)
	for index := range jobs {
		ids[index+2] = jobs[index].ID
	}
	query := `UPDATE validation_queue SET lease_owner=?, leased_until=?, job_state='leased', attempts=attempts+1 WHERE validation_job_id IN (` + scalarPlaceholders(len(jobs)) + `)`
	result, err := tx.ExecContext(ctx, query, ids...)
	if err != nil {
		return nil, err
	}
	affected, _ := result.RowsAffected()
	if affected != int64(len(jobs)) {
		return nil, fmt.Errorf("claimed %d jobs but updated %d", len(jobs), affected)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	leasedUntil := ids[1].(time.Time)
	for index := range jobs {
		jobs[index].State, jobs[index].LeaseOwner = "leased", owner
		jobs[index].Attempts++
		jobs[index].LeasedUntil = leasedUntil
	}
	return jobs, nil
}

func ExtendLease(ctx context.Context, db *sql.DB, jobID uint64, owner string, lease time.Duration) error {
	if owner == "" || lease <= 0 {
		return fmt.Errorf("owner and positive lease are required")
	}
	result, err := db.ExecContext(ctx, `UPDATE validation_queue SET leased_until=?
		WHERE validation_job_id=? AND lease_owner=? AND job_state='leased'`, time.Now().UTC().Add(lease), jobID, owner)
	return requireLeaseRow(result, err)
}

func CompleteValidationJob(ctx context.Context, db *sql.DB, jobID uint64, owner string) error {
	result, err := db.ExecContext(ctx, `UPDATE validation_queue SET job_state='completed', completed_at=UTC_TIMESTAMP(6),
		lease_owner=NULL, leased_until=NULL, last_error_code=NULL, last_error_summary=NULL
		WHERE validation_job_id=? AND lease_owner=? AND job_state='leased'`, jobID, owner)
	return requireLeaseRow(result, err)
}

func FailValidationJob(ctx context.Context, db *sql.DB, jobID uint64, owner, code, summary string, nextAttempt time.Time) error {
	result, err := db.ExecContext(ctx, `UPDATE validation_queue SET job_state='pending', next_attempt_at=?,
		lease_owner=NULL, leased_until=NULL, last_error_code=?, last_error_summary=?
		WHERE validation_job_id=? AND lease_owner=? AND job_state='leased'`, nextAttempt.UTC(),
		nullString(code), nullString(bounded(summary, 1024)), jobID, owner)
	return requireLeaseRow(result, err)
}

func requireLeaseRow(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return ErrLeaseOwnership
	}
	return nil
}
