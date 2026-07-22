package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const DatabaseCapacityBytes = uint64(50 << 30)

func ReadStorageBytes(ctx context.Context, db *sql.DB) (uint64, error) {
	var logical, allocated uint64
	if err := db.QueryRowContext(ctx, `SELECT COALESCE(SUM(data_length+index_length),0) FROM information_schema.tables WHERE table_schema=DATABASE()`).Scan(&logical); err != nil {
		return 0, err
	}
	if err := db.QueryRowContext(ctx, `SELECT COALESCE(SUM(allocated_size),0)
		FROM information_schema.innodb_tablespaces WHERE name LIKE CONCAT(DATABASE(), '/%')`).Scan(&allocated); err != nil {
		return logical, nil
	}
	if allocated > logical {
		return allocated, nil
	}
	return logical, nil
}

type MaintenancePolicy struct {
	RawPayloadDays, ParseErrorDays, AttemptDays, BatchDays, FetchDays, ExportDays int
	PauseColdSources                                                              bool
}
type MaintenanceReport struct {
	BeforeBytes, AfterBytes uint64
	Rows                    map[string]uint64
	DryRun                  bool
}

type maintenanceTarget struct {
	name, countSQL, deleteSQL string
	days                      int
}

const validationAttemptRollupDatesSQL = `SELECT DISTINCT DATE(started_at)
	FROM validation_attempts FORCE INDEX (idx_validation_attempts_cleanup)
	WHERE finished_at < ? ORDER BY DATE(started_at)`

func ValidationAttemptRollupDatesBefore(ctx context.Context, db *sql.DB, cutoff time.Time) ([]time.Time, error) {
	rows, err := db.QueryContext(ctx, validationAttemptRollupDatesSQL, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dates := make([]time.Time, 0)
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			return nil, err
		}
		dates = append(dates, date.UTC())
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return dates, nil
}

func maintenanceTargets(policy MaintenancePolicy) []maintenanceTarget {
	return []maintenanceTarget{
		{"raw_payloads", `SELECT COUNT(*) FROM raw_payloads WHERE last_seen_at < DATE_SUB(?, INTERVAL ? DAY)`, `DELETE FROM raw_payloads WHERE last_seen_at < DATE_SUB(?, INTERVAL ? DAY) LIMIT ?`, policy.RawPayloadDays},
		{"parse_errors", `SELECT COUNT(*) FROM parse_errors WHERE last_seen_at < DATE_SUB(?, INTERVAL ? DAY)`, `DELETE FROM parse_errors WHERE last_seen_at < DATE_SUB(?, INTERVAL ? DAY) LIMIT ?`, policy.ParseErrorDays},
		{"validation_attempts", `SELECT COUNT(*) FROM validation_attempts WHERE finished_at < DATE_SUB(?, INTERVAL ? DAY)`, `DELETE FROM validation_attempts WHERE finished_at < DATE_SUB(?, INTERVAL ? DAY) LIMIT ?`, policy.AttemptDays},
		{"validation_batches", `SELECT COUNT(*) FROM validation_batches b WHERE b.started_at < DATE_SUB(?, INTERVAL ? DAY) AND NOT EXISTS (SELECT 1 FROM validation_attempts a WHERE a.validation_batch_id=b.validation_batch_id)`, `DELETE FROM validation_batches WHERE validation_batch_id IN (SELECT validation_batch_id FROM (SELECT b.validation_batch_id FROM validation_batches b FORCE INDEX (idx_validation_batches_cleanup) WHERE b.started_at < DATE_SUB(?, INTERVAL ? DAY) AND NOT EXISTS (SELECT 1 FROM validation_attempts a WHERE a.validation_batch_id=b.validation_batch_id) ORDER BY b.started_at,b.validation_batch_id LIMIT ?) expired)`, policy.BatchDays},
		{"source_fetches", `SELECT COUNT(*) FROM source_fetches WHERE finished_at < DATE_SUB(?, INTERVAL ? DAY)`, `DELETE FROM source_fetches WHERE finished_at < DATE_SUB(?, INTERVAL ? DAY) LIMIT ?`, policy.FetchDays},
		{"export_members", `SELECT COUNT(*) FROM export_members WHERE created_at < DATE_SUB(?, INTERVAL ? DAY)`, `DELETE FROM export_members WHERE created_at < DATE_SUB(?, INTERVAL ? DAY) LIMIT ?`, policy.ExportDays},
	}
}

func RunMaintenance(ctx context.Context, db *sql.DB, policy MaintenancePolicy, now time.Time, batchSize int, dryRun bool) (MaintenanceReport, error) {
	if batchSize < 1 || batchSize > 10000 {
		return MaintenanceReport{}, fmt.Errorf("maintenance batch size must be between 1 and 10000")
	}
	before, err := ReadStorageBytes(ctx, db)
	if err != nil {
		return MaintenanceReport{}, err
	}
	report := MaintenanceReport{BeforeBytes: before, AfterBytes: before, Rows: map[string]uint64{}, DryRun: dryRun}
	for _, target := range maintenanceTargets(policy) {
		var count uint64
		if err := db.QueryRowContext(ctx, target.countSQL, now, target.days).Scan(&count); err != nil {
			return report, err
		}
		report.Rows[target.name] = count
		if dryRun {
			continue
		}
		for {
			result, err := db.ExecContext(ctx, target.deleteSQL, now, target.days, batchSize)
			if err != nil {
				return report, err
			}
			affected, _ := result.RowsAffected()
			if affected < int64(batchSize) {
				break
			}
		}
	}
	var missingExpiry, staleRelations uint64
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM node_configs WHERE expires_at IS NULL`).Scan(&missingExpiry); err != nil {
		return report, err
	}
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM node_source_stats
		WHERE is_active=TRUE AND last_seen_at < DATE_SUB(?, INTERVAL 30 DAY)`, now).Scan(&staleRelations); err != nil {
		return report, err
	}
	report.Rows["node_configs_missing_expiry"] = missingExpiry
	report.Rows["node_source_stats_stale"] = staleRelations
	if !dryRun {
		for {
			result, err := db.ExecContext(ctx, `UPDATE node_configs
				SET expires_at=DATE_ADD(last_seen_at, INTERVAL 30 DAY) WHERE expires_at IS NULL LIMIT ?`, batchSize)
			if err != nil {
				return report, err
			}
			affected, _ := result.RowsAffected()
			if affected < int64(batchSize) {
				break
			}
		}
		for {
			result, err := db.ExecContext(ctx, `UPDATE node_source_stats SET is_active=FALSE
				WHERE is_active=TRUE AND last_seen_at < DATE_SUB(?, INTERVAL 30 DAY) LIMIT ?`, now, batchSize)
			if err != nil {
				return report, err
			}
			affected, _ := result.RowsAffected()
			if affected < int64(batchSize) {
				break
			}
		}
	}
	var purgeCount uint64
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM node_configs n
		WHERE COALESCE(n.expires_at, DATE_ADD(n.last_seen_at, INTERVAL 30 DAY)) < ?
		AND NOT EXISTS (SELECT 1 FROM node_source_stats s WHERE s.node_config_id=n.node_config_id
			AND s.is_active=TRUE AND s.last_seen_at >= DATE_SUB(?, INTERVAL 30 DAY))`, now, now).Scan(&purgeCount); err != nil {
		return report, err
	}
	report.Rows["node_configs"] = purgeCount
	if !dryRun && purgeCount > 0 {
		_, err := db.ExecContext(ctx, `INSERT INTO node_tombstones (config_fingerprint, endpoint_fingerprint, protocol, first_seen_at, last_seen_at, last_success_at, ever_succeeded, best_quality_score, source_count, purge_reason, purged_at)
			SELECT n.config_fingerprint, UNHEX(SHA2(CONCAT(e.host,':',e.port),256)), n.protocol, n.first_seen_at, n.last_seen_at, n.last_success_at, n.last_success_at IS NOT NULL, COALESCE(s.quality_score,0), COALESCE(s.source_count,0), 'ttl_expired', ?
			FROM node_configs n JOIN endpoints e ON e.endpoint_id=n.endpoint_id LEFT JOIN node_current_status s ON s.node_config_id=n.node_config_id
			WHERE COALESCE(n.expires_at, DATE_ADD(n.last_seen_at, INTERVAL 30 DAY)) < ?
			AND NOT EXISTS (SELECT 1 FROM node_source_stats ns WHERE ns.node_config_id=n.node_config_id
				AND ns.is_active=TRUE AND ns.last_seen_at >= DATE_SUB(?, INTERVAL 30 DAY))
			ON DUPLICATE KEY UPDATE last_seen_at=VALUES(last_seen_at), purged_at=VALUES(purged_at)`, now, now, now)
		if err != nil {
			return report, err
		}
		for {
			result, err := db.ExecContext(ctx, `DELETE FROM node_configs WHERE node_config_id IN
				(SELECT node_config_id FROM (SELECT n.node_config_id FROM node_configs n
				WHERE COALESCE(n.expires_at, DATE_ADD(n.last_seen_at, INTERVAL 30 DAY)) < ?
				AND NOT EXISTS (SELECT 1 FROM node_source_stats s WHERE s.node_config_id=n.node_config_id
					AND s.is_active=TRUE AND s.last_seen_at >= DATE_SUB(?, INTERVAL 30 DAY)) LIMIT ?) expired)`, now, now, batchSize)
			if err != nil {
				return report, err
			}
			affected, _ := result.RowsAffected()
			if affected < int64(batchSize) {
				break
			}
		}
	}
	if !dryRun && policy.PauseColdSources {
		if _, err := db.ExecContext(ctx, `UPDATE sources SET state='paused_capacity' WHERE enabled=TRUE AND quality_score < 10 AND last_success_at < DATE_SUB(?, INTERVAL 30 DAY)`, now); err != nil {
			return report, err
		}
	}
	if !dryRun {
		report.AfterBytes, err = ReadStorageBytes(ctx, db)
	}
	return report, err
}
