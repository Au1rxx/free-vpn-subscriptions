package store

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/parse"
)

// FetchRecord is a persisted terminal source fetch.
type FetchRecord struct {
	ID, SourceID uint64
	StatusCode   int
	PayloadHash  [32]byte
	ErrorCode    string
	State        string
	StartedAt    time.Time
	FinishedAt   time.Time
}

// FetchWrite contains the bounded response body and safe metadata.
type FetchWrite struct {
	SourceID                                                   uint64
	StartedAt, FinishedAt                                      time.Time
	StatusCode                                                 int
	FinalURL, ETag, LastModified, ContentType, ContentEncoding string
	Body                                                       []byte
	Duration                                                   time.Duration
	ErrorCode, ErrorSummary                                    string
}

// ParseInput is a claimed successful response ready for parsing.
type ParseInput struct {
	FetchID, SourceID uint64
	FormatHint        string
	Body              []byte
}

// PersistedParse summarizes one all-or-nothing parse transaction.
type PersistedParse struct {
	ParseRunID                                  uint64
	NewEndpoints, NewConfigs, Errors, QueueJobs int
}

// FinishFetch stores a terminal fetch and, for successful bodies, a
// content-addressed gzip payload with a 30-day expiry.
func FinishFetch(ctx context.Context, db *sql.DB, write FetchWrite) (FetchRecord, error) {
	if write.StartedAt.IsZero() {
		write.StartedAt = time.Now().UTC()
	}
	if write.FinishedAt.IsZero() {
		write.FinishedAt = time.Now().UTC()
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return FetchRecord{}, err
	}
	defer tx.Rollback()
	var payloadID any
	digest := sha256.Sum256(write.Body)
	state, parseState := "failed", "skipped"
	if len(write.Body) > 0 && write.ErrorCode == "" && write.StatusCode >= 200 && write.StatusCode < 300 {
		compressed, err := compressPayload(write.Body)
		if err != nil {
			return FetchRecord{}, err
		}
		result, err := tx.ExecContext(ctx, `
			INSERT INTO raw_payloads
			  (content_sha256, content_type, content_encoding, compression, original_bytes,
			   compressed_bytes, compressed_body, first_seen_at, last_seen_at, expires_at)
			VALUES (?, ?, ?, 'gzip', ?, ?, ?, ?, ?, DATE_ADD(?, INTERVAL 30 DAY))
			ON DUPLICATE KEY UPDATE payload_id=LAST_INSERT_ID(payload_id),
			  last_seen_at=VALUES(last_seen_at), expires_at=VALUES(expires_at), reference_count=reference_count+1`,
			digest[:], nullString(write.ContentType), nullString(write.ContentEncoding), len(write.Body), len(compressed), compressed,
			write.FinishedAt, write.FinishedAt, write.FinishedAt)
		if err != nil {
			return FetchRecord{}, fmt.Errorf("upsert payload: %w", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return FetchRecord{}, err
		}
		payloadID, state, parseState = id, "success", "pending"
	}
	result, err := tx.ExecContext(ctx, `
		INSERT INTO source_fetches
		  (source_id, payload_id, started_at, finished_at, http_status, final_url, etag,
		   last_modified, content_type, content_encoding, response_bytes, duration_ms,
		   fetch_state, parse_state, error_code, error_summary)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, write.SourceID, payloadID,
		write.StartedAt, write.FinishedAt, nullInt(write.StatusCode), nullString(write.FinalURL),
		nullString(write.ETag), nullString(write.LastModified), nullString(write.ContentType),
		nullString(write.ContentEncoding), len(write.Body), durationMilliseconds(write.Duration),
		state, parseState, nullString(write.ErrorCode), nullString(bounded(write.ErrorSummary, 1024)))
	if err != nil {
		return FetchRecord{}, fmt.Errorf("insert source fetch: %w", err)
	}
	fetchID, err := result.LastInsertId()
	if err != nil {
		return FetchRecord{}, err
	}
	if state == "success" {
		_, err = tx.ExecContext(ctx, `UPDATE sources SET etag=?, last_modified=?, consecutive_failures=0,
			last_http_status=?, last_success_at=?, updated_at=UTC_TIMESTAMP(6) WHERE source_id=?`,
			nullString(write.ETag), nullString(write.LastModified), write.StatusCode, write.FinishedAt, write.SourceID)
	} else {
		_, err = tx.ExecContext(ctx, `UPDATE sources SET consecutive_failures=consecutive_failures+1,
			last_http_status=?, last_failure_at=?, updated_at=UTC_TIMESTAMP(6) WHERE source_id=?`,
			nullInt(write.StatusCode), write.FinishedAt, write.SourceID)
	}
	if err != nil {
		return FetchRecord{}, fmt.Errorf("update source fetch state: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return FetchRecord{}, err
	}
	return FetchRecord{ID: uint64(fetchID), SourceID: write.SourceID, StatusCode: write.StatusCode, PayloadHash: digest, ErrorCode: write.ErrorCode, State: state, StartedAt: write.StartedAt, FinishedAt: write.FinishedAt}, nil
}

// ClaimUnparsedFetches loads a bounded set of deduplicated payloads.
func ClaimUnparsedFetches(ctx context.Context, db *sql.DB, limit int) ([]ParseInput, error) {
	if limit < 1 || limit > 1000 {
		return nil, fmt.Errorf("parse claim limit must be between 1 and 1000")
	}
	rows, err := db.QueryContext(ctx, `SELECT f.fetch_id, f.source_id, s.format_hint, p.compressed_body, p.original_bytes
		FROM source_fetches f JOIN sources s ON s.source_id=f.source_id JOIN raw_payloads p ON p.payload_id=f.payload_id
		WHERE f.fetch_state='success' AND f.parse_state='pending' ORDER BY f.finished_at LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var inputs []ParseInput
	for rows.Next() {
		var input ParseInput
		var compressed []byte
		var originalBytes int64
		if err := rows.Scan(&input.FetchID, &input.SourceID, &input.FormatHint, &compressed, &originalBytes); err != nil {
			return nil, err
		}
		body, err := decompressPayload(compressed, originalBytes)
		if err != nil {
			return nil, fmt.Errorf("decompress fetch %d: %w", input.FetchID, err)
		}
		input.Body = body
		inputs = append(inputs, input)
	}
	return inputs, rows.Err()
}

// PersistParseResult atomically stores all identities, source relations,
// errors and validation jobs produced by one parse run.
func PersistParseResult(ctx context.Context, db *sql.DB, sourceID, fetchID uint64, result parse.Result, parserVersion string) (PersistedParse, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return PersistedParse{}, err
	}
	defer tx.Rollback()
	var existing uint64
	err = tx.QueryRowContext(ctx, `SELECT parse_run_id FROM parse_runs WHERE fetch_id=? AND parser_version=?`, fetchID, parserVersion).Scan(&existing)
	if err == nil {
		return PersistedParse{ParseRunID: existing}, tx.Commit()
	}
	if err != sql.ErrNoRows {
		return PersistedParse{}, err
	}
	started := time.Now().UTC()
	insert, err := tx.ExecContext(ctx, `INSERT INTO parse_runs
		(fetch_id, parser_version, detected_format, started_at, finished_at, input_entries,
		 success_entries, error_entries, discovered_urls, parse_state)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'success')`, fetchID, parserVersion, string(result.Format), started, time.Now().UTC(),
		len(result.Nodes)+len(result.Errors)+len(result.DiscoveredURLs), len(result.Nodes), len(result.Errors), len(result.DiscoveredURLs))
	if err != nil {
		return PersistedParse{}, fmt.Errorf("insert parse run: %w", err)
	}
	parseRunID, err := insert.LastInsertId()
	if err != nil {
		return PersistedParse{}, err
	}
	report := PersistedParse{ParseRunID: uint64(parseRunID)}
	for _, n := range result.Nodes {
		endpointID, isNew, err := upsertEndpoint(ctx, tx, n, started)
		if err != nil {
			return PersistedParse{}, err
		}
		if isNew {
			report.NewEndpoints++
		}
		nodeID, isNew, err := upsertNodeConfig(ctx, tx, endpointID, n, parserVersion, started)
		if err != nil {
			return PersistedParse{}, err
		}
		if isNew {
			report.NewConfigs++
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO node_source_stats
			(node_config_id, source_id, last_fetch_id, first_seen_at, last_seen_at)
			VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE last_fetch_id=VALUES(last_fetch_id),
			last_seen_at=VALUES(last_seen_at), seen_count=seen_count+1, is_active=TRUE`, nodeID, sourceID, fetchID, started, started); err != nil {
			return PersistedParse{}, fmt.Errorf("upsert node source relation: %w", err)
		}
		queueResult, err := tx.ExecContext(ctx, `INSERT INTO validation_queue
			(node_config_id, stage, priority, job_state, next_attempt_at)
			VALUES (?, 'connectivity', 0, 'pending', UTC_TIMESTAMP(6))
			ON DUPLICATE KEY UPDATE node_config_id=VALUES(node_config_id)`, nodeID)
		if err != nil {
			return PersistedParse{}, fmt.Errorf("queue validation: %w", err)
		}
		if affected, _ := queueResult.RowsAffected(); affected == 1 {
			report.QueueJobs++
		}
	}
	for _, entry := range result.Errors {
		sample := sha256.Sum256([]byte(entry.SampleHash))
		if _, err := tx.ExecContext(ctx, `INSERT INTO parse_errors
			(parse_run_id, source_id, fetch_id, line_number, scheme_hint, error_code,
			 sample_sha256, error_message, first_seen_at, last_seen_at, expires_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, DATE_ADD(?, INTERVAL 90 DAY))
			ON DUPLICATE KEY UPDATE last_seen_at=VALUES(last_seen_at), seen_count=seen_count+1`,
			parseRunID, sourceID, fetchID, nullInt(entry.Line), nullString(entry.Scheme), entry.Code,
			sample[:], bounded(entry.Message, 1024), started, started, started); err != nil {
			return PersistedParse{}, fmt.Errorf("insert parse error: %w", err)
		}
		report.Errors++
	}
	if _, err := tx.ExecContext(ctx, `UPDATE source_fetches SET parse_state='success' WHERE fetch_id=?`, fetchID); err != nil {
		return PersistedParse{}, err
	}
	if err := tx.Commit(); err != nil {
		return PersistedParse{}, err
	}
	return report, nil
}

func upsertEndpoint(ctx context.Context, tx *sql.Tx, n *node.Node, seen time.Time) (uint64, bool, error) {
	host, fingerprint := endpointIdentity(n)
	addressType := "domain"
	if ip := net.ParseIP(host); ip != nil {
		addressType = "ipv6"
		if ip.To4() != nil {
			addressType = "ipv4"
		}
	}
	result, err := tx.ExecContext(ctx, `INSERT INTO endpoints
		(host, host_hash, port, address_type, first_seen_at, last_seen_at)
		VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE endpoint_id=LAST_INSERT_ID(endpoint_id),
		last_seen_at=VALUES(last_seen_at)`, host, fingerprint[:], n.Port, addressType, seen, seen)
	if err != nil {
		return 0, false, fmt.Errorf("upsert endpoint: %w", err)
	}
	id, err := result.LastInsertId()
	affected, _ := result.RowsAffected()
	return uint64(id), affected == 1, err
}

func upsertNodeConfig(ctx context.Context, tx *sql.Tx, endpointID uint64, n *node.Node, parserVersion string, seen time.Time) (uint64, bool, error) {
	canonical, err := n.CanonicalJSON()
	if err != nil {
		return 0, false, err
	}
	fingerprint := n.ConfigFingerprint()
	result, err := tx.ExecContext(ctx, `INSERT INTO node_configs
		(endpoint_id, config_fingerprint, protocol, transport, security, normalized_config,
		 config_bytes, parser_version, first_seen_at, last_seen_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE node_config_id=LAST_INSERT_ID(node_config_id),
		last_seen_at=VALUES(last_seen_at), endpoint_id=VALUES(endpoint_id)`, endpointID, fingerprint[:],
		n.Protocol, defaultString(n.Network, "tcp"), defaultString(n.Security, "none"), canonical,
		len(canonical), parserVersion, seen, seen)
	if err != nil {
		return 0, false, fmt.Errorf("upsert node config: %w", err)
	}
	id, err := result.LastInsertId()
	affected, _ := result.RowsAffected()
	return uint64(id), affected == 1, err
}

func endpointIdentity(n *node.Node) (string, [32]byte) {
	host := strings.TrimSuffix(strings.ToLower(strings.TrimSpace(n.Server)), ".")
	return host, sha256.Sum256([]byte(host))
}

func compressPayload(body []byte) ([]byte, error) {
	var output bytes.Buffer
	writer, _ := gzip.NewWriterLevel(&output, gzip.BestSpeed)
	if _, err := writer.Write(body); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func decompressPayload(body []byte, expected int64) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(io.LimitReader(reader, expected+1))
}

func nullString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func nullInt(value int) any {
	if value == 0 {
		return nil
	}
	return value
}

func bounded(value string, maximum int) string {
	if len(value) > maximum {
		return value[:maximum]
	}
	return value
}

func durationMilliseconds(value time.Duration) any {
	if value <= 0 {
		return nil
	}
	return uint64(value / time.Millisecond)
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
