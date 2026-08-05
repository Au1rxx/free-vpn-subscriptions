package store

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/parse"
	"github.com/go-sql-driver/mysql"
)

const (
	nodePersistBatchSize       = 1000
	parseErrorPersistBatchSize = 500
	parseBatchMaxAttempts      = 4
)
const nodeTTL = 30 * 24 * time.Hour

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
	ProtocolHint      string
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
	if write.StatusCode == 304 && write.ErrorCode == "" {
		state = "not_modified"
	}
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
	if state != "failed" {
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
	rows, err := db.QueryContext(ctx, `SELECT f.fetch_id, f.source_id, s.format_hint, COALESCE(s.protocol_hint,''), p.compressed_body, p.original_bytes
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
		if err := rows.Scan(&input.FetchID, &input.SourceID, &input.FormatHint, &input.ProtocolHint, &compressed, &originalBytes); err != nil {
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

// RequeueSourceParses schedules successful payloads from explicitly named
// sources for a newer parser version without fetching or duplicating bodies.
func RequeueSourceParses(ctx context.Context, db *sql.DB, sourceNames []string) (int64, error) {
	if len(sourceNames) == 0 || len(sourceNames) > 100 {
		return 0, fmt.Errorf("source names must contain between 1 and 100 entries")
	}
	args := make([]any, 0, len(sourceNames))
	for _, name := range sourceNames {
		name = strings.TrimSpace(name)
		if name == "" {
			return 0, fmt.Errorf("source name must not be empty")
		}
		args = append(args, name)
	}
	result, err := db.ExecContext(ctx, `UPDATE source_fetches f JOIN sources s ON s.source_id=f.source_id
		SET f.parse_state='pending' WHERE f.fetch_state='success' AND f.payload_id IS NOT NULL
		AND s.name IN (`+scalarPlaceholders(len(args))+`)`, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// PersistParseResult durably stores identities and validation jobs in bounded
// transactions, then atomically finalizes source relations and parse metadata.
func PersistParseResult(ctx context.Context, db *sql.DB, sourceID, fetchID uint64, result parse.Result, parserVersion string) (PersistedParse, error) {
	parseRunID, completed, processedNodes, err := ensureParseRun(ctx, db, fetchID, result, parserVersion)
	if err != nil {
		return PersistedParse{}, err
	}
	if completed {
		return PersistedParse{ParseRunID: parseRunID}, nil
	}
	var sourceQuality float64
	if err := db.QueryRowContext(ctx, `SELECT quality_score FROM sources WHERE source_id=?`, sourceID).Scan(&sourceQuality); err != nil {
		return PersistedParse{}, fmt.Errorf("read source quality: %w", err)
	}
	if processedNodes > len(result.Nodes) {
		return PersistedParse{ParseRunID: parseRunID}, fmt.Errorf("parse progress %d exceeds node count %d", processedNodes, len(result.Nodes))
	}
	started := time.Now().UTC()
	report := PersistedParse{ParseRunID: parseRunID}
	for start := processedNodes; start < len(result.Nodes); start += nodePersistBatchSize {
		end := start + nodePersistBatchSize
		if end > len(result.Nodes) {
			end = len(result.Nodes)
		}
		var batchReport PersistedParse
		for attempt := 1; ; attempt++ {
			batchReport, err = persistNodeBatchTransaction(ctx, db, sourceID, fetchID, sourceQuality,
				result.Nodes[start:end], parserVersion, started, parseRunID, end)
			if err == nil {
				break
			}
			if !shouldRetryParseBatch(err, attempt) {
				return report, err
			}
			timer := time.NewTimer(time.Duration(1<<(attempt-1)) * 50 * time.Millisecond)
			select {
			case <-ctx.Done():
				timer.Stop()
				return report, ctx.Err()
			case <-timer.C:
			}
		}
		report.NewEndpoints += batchReport.NewEndpoints
		report.NewConfigs += batchReport.NewConfigs
		report.QueueJobs += batchReport.QueueJobs
	}
	if err := finishParseRun(ctx, db, sourceID, fetchID, parseRunID, result, started); err != nil {
		return report, err
	}
	report.Errors = len(result.Errors)
	return report, nil
}

func persistNodeBatchTransaction(ctx context.Context, db *sql.DB, sourceID, fetchID uint64, sourceQuality float64,
	nodes []*node.Node, parserVersion string, started time.Time, parseRunID uint64, processedNodes int) (PersistedParse, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return PersistedParse{}, err
	}
	defer tx.Rollback()
	batchReport, err := persistNodeBatch(ctx, tx, sourceID, fetchID, sourceQuality, nodes, parserVersion, started)
	if err != nil {
		return PersistedParse{}, err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE parse_runs SET error_summary=?
		WHERE parse_run_id=? AND parse_state='running'`, formatParseProgress(processedNodes), parseRunID); err != nil {
		return PersistedParse{}, fmt.Errorf("checkpoint parse progress: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return PersistedParse{}, err
	}
	return batchReport, nil
}

func shouldRetryParseBatch(err error, attempt int) bool {
	if attempt >= parseBatchMaxAttempts {
		return false
	}
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && (mysqlErr.Number == 1213 || mysqlErr.Number == 1205)
}

func ensureParseRun(ctx context.Context, db *sql.DB, fetchID uint64, result parse.Result, parserVersion string) (uint64, bool, int, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, false, 0, err
	}
	defer tx.Rollback()
	var existing uint64
	var state string
	var summary sql.NullString
	err = tx.QueryRowContext(ctx, `SELECT parse_run_id, parse_state, error_summary
		FROM parse_runs WHERE fetch_id=? AND parser_version=?`, fetchID, parserVersion).Scan(&existing, &state, &summary)
	if err == nil {
		if err := tx.Commit(); err != nil {
			return 0, false, 0, err
		}
		return existing, state == "success", parseProgress(summary), nil
	}
	if err != sql.ErrNoRows {
		return 0, false, 0, err
	}
	started := time.Now().UTC()
	insert, err := tx.ExecContext(ctx, `INSERT INTO parse_runs
		(fetch_id, parser_version, detected_format, started_at, finished_at, input_entries,
		 success_entries, error_entries, discovered_urls, parse_state)
		VALUES (?, ?, ?, ?, NULL, ?, ?, ?, ?, 'running')`, fetchID, parserVersion, string(result.Format), started,
		len(result.Nodes)+len(result.Errors)+len(result.DiscoveredURLs), len(result.Nodes), len(result.Errors), len(result.DiscoveredURLs))
	if err != nil {
		return 0, false, 0, fmt.Errorf("insert parse run: %w", err)
	}
	parseRunID, err := insert.LastInsertId()
	if err != nil {
		return 0, false, 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, false, 0, err
	}
	return uint64(parseRunID), false, 0, nil
}

func formatParseProgress(processedNodes int) string {
	return fmt.Sprintf("processed_nodes=%d", processedNodes)
}

func parseProgress(summary sql.NullString) int {
	if !summary.Valid || !strings.HasPrefix(summary.String, "processed_nodes=") {
		return 0
	}
	processed, err := strconv.Atoi(strings.TrimPrefix(summary.String, "processed_nodes="))
	if err != nil || processed < 0 {
		return 0
	}
	return processed
}

func finishParseRun(ctx context.Context, db *sql.DB, sourceID, fetchID, parseRunID uint64, result parse.Result, started time.Time) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if len(result.Nodes) > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE node_source_stats SET is_active=FALSE
			WHERE source_id=? AND last_fetch_id<>? AND is_active=TRUE`, sourceID, fetchID); err != nil {
			return fmt.Errorf("deactivate missing source nodes: %w", err)
		}
	}
	if err := insertParseErrors(ctx, tx, parseRunID, sourceID, fetchID, result.Errors, started); err != nil {
		return err
	}
	finished := time.Now().UTC()
	if _, err := tx.ExecContext(ctx, `UPDATE parse_runs SET detected_format=?, finished_at=?,
		input_entries=?, success_entries=?, error_entries=?, discovered_urls=?, parse_state='success', error_summary=NULL
		WHERE parse_run_id=?`, string(result.Format), finished,
		len(result.Nodes)+len(result.Errors)+len(result.DiscoveredURLs), len(result.Nodes), len(result.Errors), len(result.DiscoveredURLs), parseRunID); err != nil {
		return fmt.Errorf("finish parse run: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE source_fetches SET parse_state='success' WHERE fetch_id=?`, fetchID); err != nil {
		return err
	}
	return tx.Commit()
}

func insertParseErrors(ctx context.Context, tx *sql.Tx, parseRunID, sourceID, fetchID uint64, entries []parse.EntryError, seen time.Time) error {
	for start := 0; start < len(entries); start += parseErrorPersistBatchSize {
		end := start + parseErrorPersistBatchSize
		if end > len(entries) {
			end = len(entries)
		}
		var query strings.Builder
		query.WriteString(`INSERT INTO parse_errors
			(parse_run_id, source_id, fetch_id, line_number, scheme_hint, error_code,
			 sample_sha256, error_message, first_seen_at, last_seen_at, expires_at) VALUES `)
		args := make([]any, 0, (end-start)*11)
		for index, entry := range entries[start:end] {
			if index > 0 {
				query.WriteByte(',')
			}
			query.WriteString("(?,?,?,?,?,?,?,?,?,?,DATE_ADD(?, INTERVAL 90 DAY))")
			sample := sha256.Sum256([]byte(entry.SampleHash))
			args = append(args, parseRunID, sourceID, fetchID, nullInt(entry.Line), nullString(entry.Scheme),
				entry.Code, sample[:], bounded(entry.Message, 1024), seen, seen, seen)
		}
		query.WriteString(` ON DUPLICATE KEY UPDATE
			last_seen_at=VALUES(last_seen_at), seen_count=seen_count+1`)
		if _, err := tx.ExecContext(ctx, query.String(), args...); err != nil {
			return fmt.Errorf("insert parse error batch: %w", err)
		}
	}
	return nil
}

type preparedNode struct {
	n                 *node.Node
	host              string
	hostHash          [32]byte
	configFingerprint [32]byte
	canonical         []byte
}

func persistNodeBatch(ctx context.Context, tx *sql.Tx, sourceID, fetchID uint64, sourceQuality float64, nodes []*node.Node, parserVersion string, seen time.Time) (PersistedParse, error) {
	prepared := make([]preparedNode, 0, len(nodes))
	for _, n := range nodes {
		canonical, err := n.CanonicalJSON()
		if err != nil {
			return PersistedParse{}, err
		}
		host, hostHash := endpointIdentity(n)
		prepared = append(prepared, preparedNode{n: n, host: host, hostHash: hostHash, configFingerprint: n.ConfigFingerprint(), canonical: canonical})
	}
	existingEndpoints, err := selectEndpointIDs(ctx, tx, prepared)
	if err != nil {
		return PersistedParse{}, err
	}
	report := PersistedParse{}
	uniqueEndpoints := make(map[string]bool)
	var endpointSQL strings.Builder
	endpointSQL.WriteString(`INSERT INTO endpoints (host, host_hash, port, address_type, first_seen_at, last_seen_at) VALUES `)
	endpointArgs := make([]any, 0, len(prepared)*6)
	for index, item := range prepared {
		if index > 0 {
			endpointSQL.WriteByte(',')
		}
		endpointSQL.WriteString("(?,?,?,?,?,?)")
		addressType := "domain"
		if ip := net.ParseIP(item.host); ip != nil {
			addressType = "ipv6"
			if ip.To4() != nil {
				addressType = "ipv4"
			}
		}
		endpointArgs = append(endpointArgs, item.host, item.hostHash[:], item.n.Port, addressType, seen, seen)
		key := endpointMapKey(item.hostHash[:], item.n.Port)
		if !uniqueEndpoints[key] {
			uniqueEndpoints[key] = true
			if existingEndpoints[key] == 0 {
				report.NewEndpoints++
			}
		}
	}
	endpointSQL.WriteString(` ON DUPLICATE KEY UPDATE last_seen_at=VALUES(last_seen_at)`)
	if _, err := tx.ExecContext(ctx, endpointSQL.String(), endpointArgs...); err != nil {
		return PersistedParse{}, fmt.Errorf("batch upsert endpoints: %w", err)
	}
	endpointIDs, err := selectEndpointIDs(ctx, tx, prepared)
	if err != nil {
		return PersistedParse{}, err
	}
	existingConfigs, err := selectConfigIDs(ctx, tx, prepared)
	if err != nil {
		return PersistedParse{}, err
	}
	var configSQL strings.Builder
	configSQL.WriteString(`INSERT INTO node_configs (endpoint_id, config_fingerprint, protocol, transport, security, normalized_config, config_bytes, parser_version, first_seen_at, last_seen_at, expires_at) VALUES `)
	configArgs := make([]any, 0, len(prepared)*11)
	uniqueConfigs := make(map[string]bool)
	for index, item := range prepared {
		if index > 0 {
			configSQL.WriteByte(',')
		}
		configSQL.WriteString("(?,?,?,?,?,?,?,?,?,?,?)")
		endpointID := endpointIDs[endpointMapKey(item.hostHash[:], item.n.Port)]
		configArgs = append(configArgs, endpointID, item.configFingerprint[:], item.n.Protocol,
			classificationValue(item.n.Network, "tcp"), classificationValue(item.n.Security, "none"), item.canonical,
			len(item.canonical), parserVersion, seen, seen, nodeExpiresAt(seen))
		key := string(item.configFingerprint[:])
		if !uniqueConfigs[key] {
			uniqueConfigs[key] = true
			if existingConfigs[key] == 0 {
				report.NewConfigs++
			}
		}
	}
	configSQL.WriteString(` ON DUPLICATE KEY UPDATE last_seen_at=VALUES(last_seen_at), expires_at=VALUES(expires_at), endpoint_id=VALUES(endpoint_id)`)
	if _, err := tx.ExecContext(ctx, configSQL.String(), configArgs...); err != nil {
		return PersistedParse{}, fmt.Errorf("batch upsert node configs: %w", err)
	}
	configIDs, err := selectConfigIDs(ctx, tx, prepared)
	if err != nil {
		return PersistedParse{}, err
	}
	nodeIDs := make([]uint64, 0, len(uniqueConfigs))
	seenNodeIDs := make(map[uint64]bool)
	for fingerprint := range uniqueConfigs {
		id := configIDs[fingerprint]
		if id == 0 {
			return PersistedParse{}, fmt.Errorf("node config identity was not returned")
		}
		if !seenNodeIDs[id] {
			seenNodeIDs[id] = true
			nodeIDs = append(nodeIDs, id)
		}
	}
	if err := upsertNodeRelations(ctx, tx, nodeIDs, sourceID, fetchID, sourceQuality, seen); err != nil {
		return PersistedParse{}, err
	}
	existingQueue, err := selectQueuedNodeIDs(ctx, tx, nodeIDs)
	if err != nil {
		return PersistedParse{}, err
	}
	if err := upsertValidationQueue(ctx, tx, nodeIDs); err != nil {
		return PersistedParse{}, err
	}
	for _, id := range nodeIDs {
		if !existingQueue[id] {
			report.QueueJobs++
		}
	}
	return report, nil
}

func classificationValue(value, fallback string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return fallback
	}
	if len(value) > 32 {
		return "other"
	}
	for _, char := range value {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '-' && char != '_' && char != '.' {
			return "other"
		}
	}
	return value
}

func nodeExpiresAt(seen time.Time) time.Time { return seen.Add(nodeTTL) }

func selectEndpointIDs(ctx context.Context, tx *sql.Tx, prepared []preparedNode) (map[string]uint64, error) {
	query := `SELECT endpoint_id, host_hash, port FROM endpoints WHERE (host_hash, port) IN (` + rowPlaceholders(len(prepared), 2) + `)`
	args := make([]any, 0, len(prepared)*2)
	for _, item := range prepared {
		args = append(args, item.hostHash[:], item.n.Port)
	}
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select endpoint identities: %w", err)
	}
	defer rows.Close()
	ids := make(map[string]uint64)
	for rows.Next() {
		var id uint64
		var hash []byte
		var port int
		if err := rows.Scan(&id, &hash, &port); err != nil {
			return nil, err
		}
		ids[endpointMapKey(hash, port)] = id
	}
	return ids, rows.Err()
}

func selectConfigIDs(ctx context.Context, tx *sql.Tx, prepared []preparedNode) (map[string]uint64, error) {
	query := `SELECT node_config_id, config_fingerprint FROM node_configs WHERE config_fingerprint IN (` + scalarPlaceholders(len(prepared)) + `)`
	args := make([]any, 0, len(prepared))
	for _, item := range prepared {
		args = append(args, item.configFingerprint[:])
	}
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select config identities: %w", err)
	}
	defer rows.Close()
	ids := make(map[string]uint64)
	for rows.Next() {
		var id uint64
		var fingerprint []byte
		if err := rows.Scan(&id, &fingerprint); err != nil {
			return nil, err
		}
		ids[string(fingerprint)] = id
	}
	return ids, rows.Err()
}

func upsertNodeRelations(ctx context.Context, tx *sql.Tx, nodeIDs []uint64, sourceID, fetchID uint64, sourceQuality float64, seen time.Time) error {
	var query strings.Builder
	query.WriteString(`INSERT INTO node_source_stats (node_config_id, source_id, last_fetch_id, first_seen_at, last_seen_at, source_quality) VALUES `)
	args := make([]any, 0, len(nodeIDs)*6)
	for index, id := range nodeIDs {
		if index > 0 {
			query.WriteByte(',')
		}
		query.WriteString("(?,?,?,?,?,?)")
		args = append(args, id, sourceID, fetchID, seen, seen, sourceQuality)
	}
	query.WriteString(` ON DUPLICATE KEY UPDATE
		seen_count=seen_count+IF(last_fetch_id<>VALUES(last_fetch_id),1,0),
		last_fetch_id=VALUES(last_fetch_id), last_seen_at=VALUES(last_seen_at),
		source_quality=VALUES(source_quality), is_active=TRUE`)
	_, err := tx.ExecContext(ctx, query.String(), args...)
	return err
}

func selectQueuedNodeIDs(ctx context.Context, tx *sql.Tx, nodeIDs []uint64) (map[uint64]bool, error) {
	query := `SELECT node_config_id FROM validation_queue WHERE stage='connectivity' AND node_config_id IN (` + scalarPlaceholders(len(nodeIDs)) + `)`
	args := make([]any, len(nodeIDs))
	for index, id := range nodeIDs {
		args[index] = id
	}
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := make(map[uint64]bool)
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids[id] = true
	}
	return ids, rows.Err()
}

func upsertValidationQueue(ctx context.Context, tx *sql.Tx, nodeIDs []uint64) error {
	var query strings.Builder
	query.WriteString(`INSERT INTO validation_queue (node_config_id, stage, priority, job_state, next_attempt_at) VALUES `)
	args := make([]any, 0, len(nodeIDs))
	for index, id := range nodeIDs {
		if index > 0 {
			query.WriteByte(',')
		}
		query.WriteString("(?,'connectivity',0,'pending',UTC_TIMESTAMP(6))")
		args = append(args, id)
	}
	query.WriteString(` ON DUPLICATE KEY UPDATE node_config_id=VALUES(node_config_id)`)
	_, err := tx.ExecContext(ctx, query.String(), args...)
	return err
}

func rowPlaceholders(rows, columns int) string {
	row := "(" + strings.TrimSuffix(strings.Repeat("?,", columns), ",") + ")"
	return strings.TrimSuffix(strings.Repeat(row+",", rows), ",")
}

func scalarPlaceholders(count int) string {
	return strings.TrimSuffix(strings.Repeat("?,", count), ",")
}

func endpointMapKey(hash []byte, port int) string {
	return string(hash) + fmt.Sprintf("/%d", port)
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
		 config_bytes, parser_version, first_seen_at, last_seen_at, expires_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE node_config_id=LAST_INSERT_ID(node_config_id),
		last_seen_at=VALUES(last_seen_at), expires_at=VALUES(expires_at), endpoint_id=VALUES(endpoint_id)`, endpointID, fingerprint[:],
		n.Protocol, defaultString(n.Network, "tcp"), defaultString(n.Security, "none"), canonical,
		len(canonical), parserVersion, seen, seen, nodeExpiresAt(seen))
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
