package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ClassificationCandidate struct {
	NodeConfigID                                                       uint64
	Protocol, Transport, Security, Availability, ExitCountry, ExitASN  string
	EntryHost, IPVersion                                               string
	LastSeenAt                                                         time.Time
	LastValidationAt                                                   sql.NullTime
	LatencyMS, SourceCount                                             int
	Success7D, Stability30D, Consistency, ExitStability, Compatibility float64
}

type ClassificationUpdate struct {
	NodeConfigID                                                  uint64
	Protocol, Transport, Security, FreshnessClass, StabilityClass string
	ExitCountry, ExitASN, Grade                                   string
	EntryCountry, EntryRegion, EntryCity, EntryTimeZone           string
	EntryASN, EntryOrganization, ProviderClass, IPVersion         string
	Score                                                         int
	Breakdown                                                     any
}

func CountUnclassified(ctx context.Context, db *sql.DB) (int, error) {
	var count int
	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM node_configs n
		LEFT JOIN node_classifications c ON c.node_config_id=n.node_config_id WHERE c.node_config_id IS NULL`).Scan(&count)
	return count, err
}

func ListClassificationCandidates(ctx context.Context, db *sql.DB, limit int) ([]ClassificationCandidate, error) {
	if limit < 1 || limit > 10000 {
		return nil, fmt.Errorf("classification limit must be between 1 and 10000")
	}
	ids, err := listClassificationCandidateIDs(ctx, db, limit)
	if err != nil || len(ids) == 0 {
		return nil, err
	}
	args := make([]any, len(ids))
	for index, id := range ids {
		args[index] = id
	}
	rows, err := db.QueryContext(ctx, `SELECT n.node_config_id, n.protocol, n.transport, n.security,
		COALESCE(s.availability_state,'unverified'), n.last_seen_at, s.last_validation_at,
		COALESCE(s.latency_p50_ms,0), COALESCE(s.source_count,0), COALESCE(s.exit_country,''), COALESCE(s.exit_asn,''),
		e.host, e.address_type,
		COALESCE(AVG(CASE WHEN a.started_at >= DATE_SUB(UTC_TIMESTAMP(), INTERVAL 7 DAY) THEN a.passed END),0),
		COALESCE(AVG(CASE WHEN a.started_at >= DATE_SUB(UTC_TIMESTAMP(), INTERVAL 30 DAY) THEN a.passed END),0),
		COALESCE(AVG(CASE WHEN a.started_at >= DATE_SUB(UTC_TIMESTAMP(), INTERVAL 30 DAY) THEN a.config_accepted END),0),
		CASE WHEN COUNT(DISTINCT CASE WHEN a.passed THEN a.exit_ip END) <= 1 THEN 1 ELSE 1.0/COUNT(DISTINCT CASE WHEN a.passed THEN a.exit_ip END) END,
		COALESCE(AVG(a.config_accepted),0)
		FROM node_configs n JOIN endpoints e ON e.endpoint_id=n.endpoint_id
		LEFT JOIN node_current_status s ON s.node_config_id=n.node_config_id
		LEFT JOIN validation_attempts a ON a.node_config_id=n.node_config_id
		WHERE n.node_config_id IN (`+scalarPlaceholders(len(ids))+`)
		GROUP BY n.node_config_id, n.protocol, n.transport, n.security, s.availability_state,
		 n.last_seen_at, s.last_validation_at, s.latency_p50_ms, s.source_count, s.exit_country, s.exit_asn,
		 e.host, e.address_type ORDER BY n.node_config_id`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var candidates []ClassificationCandidate
	for rows.Next() {
		var item ClassificationCandidate
		if err := rows.Scan(&item.NodeConfigID, &item.Protocol, &item.Transport, &item.Security, &item.Availability,
			&item.LastSeenAt, &item.LastValidationAt, &item.LatencyMS, &item.SourceCount, &item.ExitCountry, &item.ExitASN,
			&item.EntryHost, &item.IPVersion,
			&item.Success7D, &item.Stability30D, &item.Consistency, &item.ExitStability, &item.Compatibility); err != nil {
			return nil, err
		}
		candidates = append(candidates, item)
	}
	return candidates, rows.Err()
}

func listClassificationCandidateIDs(ctx context.Context, db *sql.DB, limit int) ([]uint64, error) {
	rows, err := db.QueryContext(ctx, `SELECT n.node_config_id FROM node_configs n
		LEFT JOIN node_classifications c ON c.node_config_id=n.node_config_id
		WHERE c.node_config_id IS NULL ORDER BY n.node_config_id LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	ids, err := scanNodeConfigIDs(rows)
	if err != nil || len(ids) == limit {
		return ids, err
	}
	rows, err = db.QueryContext(ctx, `SELECT c.node_config_id FROM node_classifications c
		LEFT JOIN node_current_status s ON s.node_config_id=c.node_config_id
		ORDER BY COALESCE(s.last_validation_at, '1970-01-01') DESC, c.classified_at, c.node_config_id LIMIT ?`, limit-len(ids))
	if err != nil {
		return nil, err
	}
	classified, err := scanNodeConfigIDs(rows)
	return append(ids, classified...), err
}

func scanNodeConfigIDs(rows *sql.Rows) ([]uint64, error) {
	defer rows.Close()
	var ids []uint64
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func WriteClassifications(ctx context.Context, db *sql.DB, updates []ClassificationUpdate, classifiedAt time.Time) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for start := 0; start < len(updates); start += 500 {
		end := start + 500
		if end > len(updates) {
			end = len(updates)
		}
		if err := writeClassificationBatch(ctx, tx, updates[start:end], classifiedAt); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func writeClassificationBatch(ctx context.Context, tx *sql.Tx, updates []ClassificationUpdate, classifiedAt time.Time) error {
	classificationValues := make([]string, 0, len(updates))
	classificationArgs := make([]any, 0, len(updates)*9)
	statusValues := make([]string, 0, len(updates))
	statusArgs := make([]any, 0, len(updates)*4)
	ids := make([]any, 0, len(updates))
	for _, update := range updates {
		breakdown, err := json.Marshal(update.Breakdown)
		if err != nil {
			return err
		}
		classificationValues = append(classificationValues, "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?, 'fnctl-2', ?)")
		classificationArgs = append(classificationArgs, update.NodeConfigID, update.Protocol, update.Transport, update.Security,
			nullString(update.IPVersion), nullString(update.EntryCountry), nullString(update.EntryRegion), nullString(update.EntryCity),
			nullString(update.EntryTimeZone), nullString(update.EntryASN), nullString(update.EntryOrganization), nullString(update.ProviderClass),
			nullString(update.ExitCountry), nullString(update.ExitASN), update.FreshnessClass, update.StabilityClass, classifiedAt)
		statusValues = append(statusValues, "(?,?,?,?)")
		statusArgs = append(statusArgs, update.NodeConfigID, update.Score, update.Grade, breakdown)
		ids = append(ids, update.NodeConfigID)
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO node_classifications
		(node_config_id, protocol, transport, security, ip_version, entry_country, entry_region, entry_city,
		 entry_timezone, entry_asn, entry_organization, provider_class, exit_country, exit_asn, freshness_class,
		 stability_class, classifier_version, classified_at) VALUES `+strings.Join(classificationValues, ",")+`
		ON DUPLICATE KEY UPDATE protocol=VALUES(protocol), transport=VALUES(transport), security=VALUES(security),
		 ip_version=VALUES(ip_version), entry_country=VALUES(entry_country), entry_region=VALUES(entry_region),
		 entry_city=VALUES(entry_city), entry_timezone=VALUES(entry_timezone), entry_asn=VALUES(entry_asn),
		 entry_organization=VALUES(entry_organization), provider_class=VALUES(provider_class),
		 exit_country=VALUES(exit_country), exit_asn=VALUES(exit_asn), freshness_class=VALUES(freshness_class),
		 stability_class=VALUES(stability_class), classifier_version=VALUES(classifier_version), classified_at=VALUES(classified_at)`, classificationArgs...)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO node_current_status
		(node_config_id, quality_score, quality_grade, score_breakdown) VALUES `+strings.Join(statusValues, ",")+`
		ON DUPLICATE KEY UPDATE quality_score=VALUES(quality_score), quality_grade=VALUES(quality_grade),
		 score_breakdown=VALUES(score_breakdown)`, statusArgs...)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE node_configs n JOIN node_current_status s ON s.node_config_id=n.node_config_id
		SET n.is_exportable=(s.quality_grade <> 'U') WHERE n.node_config_id IN (`+scalarPlaceholders(len(ids))+`)`, ids...)
	return err
}

func RollupDailyStats(ctx context.Context, db *sql.DB, date time.Time) (int64, error) {
	result, err := db.ExecContext(ctx, `INSERT INTO node_daily_stats
		(stat_date, node_config_id, validation_count, success_count, partial_count, failure_count,
		 success_rate, latency_p50_ms, latency_p95_ms, source_count, quality_score, quality_grade)
		SELECT DATE(?), a.node_config_id, COUNT(*), SUM(a.passed), SUM(a.partial_success),
		 SUM(NOT a.passed AND NOT a.partial_success), AVG(a.passed),
		 AVG(a.http_median_ms), MAX(a.http_median_ms), MAX(COALESCE(s.source_count,0)),
		 MAX(COALESCE(s.quality_score,0)), MAX(COALESCE(s.quality_grade,'U'))
		FROM validation_attempts a LEFT JOIN node_current_status s ON s.node_config_id=a.node_config_id
		WHERE a.started_at >= DATE(?) AND a.started_at < DATE_ADD(DATE(?), INTERVAL 1 DAY)
		GROUP BY a.node_config_id
		ON DUPLICATE KEY UPDATE validation_count=VALUES(validation_count), success_count=VALUES(success_count),
		 partial_count=VALUES(partial_count), failure_count=VALUES(failure_count), success_rate=VALUES(success_rate),
		 latency_p50_ms=VALUES(latency_p50_ms), latency_p95_ms=VALUES(latency_p95_ms),
		 source_count=VALUES(source_count), quality_score=VALUES(quality_score), quality_grade=VALUES(quality_grade)`, date, date, date)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
