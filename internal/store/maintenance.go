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
	RawPayloadDays, ParseErrorDays, AttemptDays, FetchDays, ExportDays int
	PauseColdSources, StoreRawBodies                                   bool
}
type MaintenanceReport struct {
	BeforeBytes, AfterBytes uint64
	Rows                    map[string]uint64
	DryRun                  bool
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
	targets := []struct {
		name, countSQL, deleteSQL string
		days                      int
	}{
		{"raw_payloads", `SELECT COUNT(*) FROM raw_payloads WHERE last_seen_at < DATE_SUB(?, INTERVAL ? DAY)`, `DELETE FROM raw_payloads WHERE last_seen_at < DATE_SUB(?, INTERVAL ? DAY) LIMIT ?`, policy.RawPayloadDays},
		{"parse_errors", `SELECT COUNT(*) FROM parse_errors WHERE last_seen_at < DATE_SUB(?, INTERVAL ? DAY)`, `DELETE FROM parse_errors WHERE last_seen_at < DATE_SUB(?, INTERVAL ? DAY) LIMIT ?`, policy.ParseErrorDays},
		{"validation_attempts", `SELECT COUNT(*) FROM validation_attempts WHERE finished_at < DATE_SUB(?, INTERVAL ? DAY)`, `DELETE FROM validation_attempts WHERE finished_at < DATE_SUB(?, INTERVAL ? DAY) LIMIT ?`, policy.AttemptDays},
		{"source_fetches", `SELECT COUNT(*) FROM source_fetches WHERE finished_at < DATE_SUB(?, INTERVAL ? DAY)`, `DELETE FROM source_fetches WHERE finished_at < DATE_SUB(?, INTERVAL ? DAY) LIMIT ?`, policy.FetchDays},
		{"export_members", `SELECT COUNT(*) FROM export_members WHERE created_at < DATE_SUB(?, INTERVAL ? DAY)`, `DELETE FROM export_members WHERE created_at < DATE_SUB(?, INTERVAL ? DAY) LIMIT ?`, policy.ExportDays},
	}
	for _, target := range targets {
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
	var purgeCount uint64
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM node_configs n WHERE n.expires_at < ? AND NOT EXISTS (SELECT 1 FROM node_source_stats s WHERE s.node_config_id=n.node_config_id AND s.is_active=TRUE)`, now).Scan(&purgeCount); err != nil {
		return report, err
	}
	report.Rows["node_configs"] = purgeCount
	if !dryRun && purgeCount > 0 {
		_, err := db.ExecContext(ctx, `INSERT INTO node_tombstones (config_fingerprint, endpoint_fingerprint, protocol, first_seen_at, last_seen_at, last_success_at, ever_succeeded, best_quality_score, source_count, purge_reason, purged_at)
			SELECT n.config_fingerprint, UNHEX(SHA2(CONCAT(e.host,':',e.port),256)), n.protocol, n.first_seen_at, n.last_seen_at, n.last_success_at, n.last_success_at IS NOT NULL, COALESCE(s.quality_score,0), COALESCE(s.source_count,0), 'ttl_expired', ?
			FROM node_configs n JOIN endpoints e ON e.endpoint_id=n.endpoint_id LEFT JOIN node_current_status s ON s.node_config_id=n.node_config_id
			WHERE n.expires_at < ? AND NOT EXISTS (SELECT 1 FROM node_source_stats ns WHERE ns.node_config_id=n.node_config_id AND ns.is_active=TRUE)
			ON DUPLICATE KEY UPDATE last_seen_at=VALUES(last_seen_at), purged_at=VALUES(purged_at)`, now, now)
		if err != nil {
			return report, err
		}
		for {
			result, err := db.ExecContext(ctx, `DELETE FROM node_configs WHERE node_config_id IN (SELECT node_config_id FROM (SELECT n.node_config_id FROM node_configs n WHERE n.expires_at < ? AND NOT EXISTS (SELECT 1 FROM node_source_stats s WHERE s.node_config_id=n.node_config_id AND s.is_active=TRUE) LIMIT ?) expired)`, now, batchSize)
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
