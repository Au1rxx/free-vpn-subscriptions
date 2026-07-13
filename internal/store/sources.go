package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// SourceRecord is one fetchable or discoverable public source.
type SourceRecord struct {
	ID                                               uint64
	Kind, Name, URL, CanonicalURL, FormatHint, State string
	DiscoveryMethod                                  string
	Depth                                            int
	Enabled                                          bool
	Priority                                         int
	FetchInterval                                    time.Duration
	ETag, LastModified                               string
}

// CanonicalizeSourceURL normalizes identity without discarding query tokens
// that may be required by subscription endpoints.
func CanonicalizeSourceURL(raw string) (string, [32]byte, error) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", [32]byte{}, err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" && parsed.Scheme != "HTTPS" && parsed.Scheme != "HTTP" {
		return "", [32]byte{}, fmt.Errorf("unsupported source scheme %q", parsed.Scheme)
	}
	if parsed.Host == "" {
		return "", [32]byte{}, fmt.Errorf("source URL has no host")
	}
	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = strings.ToLower(parsed.Host)
	parsed.Fragment = ""
	parsed.RawQuery = parsed.Query().Encode()
	canonical := parsed.String()
	return canonical, sha256.Sum256([]byte(canonical)), nil
}

// UpsertSource inserts a source by canonical URL or refreshes its metadata.
func UpsertSource(ctx context.Context, db *sql.DB, source SourceRecord) (SourceRecord, error) {
	canonical, digest, err := CanonicalizeSourceURL(source.URL)
	if err != nil {
		return SourceRecord{}, err
	}
	if source.Kind == "" {
		source.Kind = "subscription"
	}
	if source.FormatHint == "" {
		source.FormatHint = "auto"
	}
	if source.State == "" {
		source.State = "active"
	}
	if source.DiscoveryMethod == "" {
		source.DiscoveryMethod = "seed"
	}
	if source.FetchInterval <= 0 {
		source.FetchInterval = time.Hour
	}
	result, err := db.ExecContext(ctx, `
		INSERT INTO sources
		  (name, kind, url, canonical_url, canonical_url_hash, format_hint, discovery_method,
		   state, enabled, priority, depth, fetch_interval_seconds, next_fetch_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(6))
		ON DUPLICATE KEY UPDATE
		  source_id=LAST_INSERT_ID(source_id), name=VALUES(name), kind=VALUES(kind),
		  url=VALUES(url), format_hint=VALUES(format_hint), state=VALUES(state),
		  enabled=VALUES(enabled), priority=VALUES(priority), depth=LEAST(depth, VALUES(depth)),
		  updated_at=UTC_TIMESTAMP(6)`, source.Name, source.Kind, source.URL, canonical, digest[:],
		source.FormatHint, source.DiscoveryMethod, source.State, source.Enabled, source.Priority,
		source.Depth, uint64(source.FetchInterval/time.Second))
	if err != nil {
		return SourceRecord{}, fmt.Errorf("upsert source: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return SourceRecord{}, fmt.Errorf("read source id: %w", err)
	}
	source.ID, source.CanonicalURL = uint64(id), canonical
	return source, nil
}

// ClaimDueSources returns a bounded priority-ordered fetch batch. The caller
// records a terminal fetch immediately, while next_fetch_at prevents overlap.
func ClaimDueSources(ctx context.Context, db *sql.DB, limit int) ([]SourceRecord, error) {
	if limit < 1 || limit > 1000 {
		return nil, fmt.Errorf("source claim limit must be between 1 and 1000")
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, `
		SELECT source_id, kind, name, url, canonical_url, format_hint, state, depth,
		       enabled, priority, fetch_interval_seconds, COALESCE(etag,''), COALESCE(last_modified,'')
		FROM sources
		WHERE enabled=TRUE AND state='active' AND (next_fetch_at IS NULL OR next_fetch_at <= UTC_TIMESTAMP(6))
		ORDER BY priority DESC, COALESCE(next_fetch_at, '1970-01-01') ASC
		LIMIT ? FOR UPDATE SKIP LOCKED`, limit)
	if err != nil {
		return nil, fmt.Errorf("claim due sources: %w", err)
	}
	defer rows.Close()
	var sources []SourceRecord
	for rows.Next() {
		var source SourceRecord
		var interval uint64
		if err := rows.Scan(&source.ID, &source.Kind, &source.Name, &source.URL, &source.CanonicalURL,
			&source.FormatHint, &source.State, &source.Depth, &source.Enabled, &source.Priority,
			&interval, &source.ETag, &source.LastModified); err != nil {
			return nil, err
		}
		source.FetchInterval = time.Duration(interval) * time.Second
		sources = append(sources, source)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for _, source := range sources {
		if _, err := tx.ExecContext(ctx, `UPDATE sources SET next_fetch_at=DATE_ADD(UTC_TIMESTAMP(6), INTERVAL ? SECOND) WHERE source_id=?`, uint64(source.FetchInterval/time.Second), source.ID); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return sources, nil
}

func DisableDiscoveredSources(ctx context.Context, db *sql.DB, kind, method string) (int64, error) {
	if strings.TrimSpace(kind) == "" || strings.TrimSpace(method) == "" {
		return 0, fmt.Errorf("source kind and discovery method are required")
	}
	result, err := db.ExecContext(ctx, `UPDATE sources SET enabled=FALSE, state='paused',
		updated_at=UTC_TIMESTAMP(6) WHERE kind=? AND discovery_method=?`, kind, method)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
