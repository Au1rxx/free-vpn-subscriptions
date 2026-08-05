// Package ingest orchestrates durable source collection.
package ingest

import (
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const defaultSpoolBytes = int64(2 << 30)

// FetchEnvelope is the credential-sensitive response persisted during a
// database outage. The spool directory and files are owner-only.
type FetchEnvelope struct {
	SourceID        uint64    `json:"source_id"`
	FetchedAt       time.Time `json:"fetched_at"`
	StatusCode      int       `json:"status_code"`
	FinalURL        string    `json:"final_url,omitempty"`
	ETag            string    `json:"etag,omitempty"`
	LastModified    string    `json:"last_modified,omitempty"`
	ContentType     string    `json:"content_type,omitempty"`
	ContentEncoding string    `json:"content_encoding,omitempty"`
	Body            []byte    `json:"body,omitempty"`
	DurationMS      uint64    `json:"duration_ms,omitempty"`
	ErrorCode       string    `json:"error_code,omitempty"`
	ErrorSummary    string    `json:"error_summary,omitempty"`
}

// Persister accepts one replayed fetch. Success authorizes file deletion.
type Persister interface {
	PersistFetch(context.Context, FetchEnvelope) error
}

type ReplayReport struct {
	Persisted, Quarantined, Failed int
}

type Spool struct {
	Dir      string
	MaxBytes int64
}

type SpoolError struct {
	Code string
	Err  error
}

func (e *SpoolError) Error() string { return e.Code + ": " + e.Err.Error() }
func (e *SpoolError) Unwrap() error { return e.Err }

func spoolErrorCode(err error) string {
	if typed, ok := err.(*SpoolError); ok {
		return typed.Code
	}
	return ""
}

func NewSpool(directory string, maximumBytes int64) (*Spool, error) {
	if directory == "" {
		return nil, fmt.Errorf("spool directory is required")
	}
	if maximumBytes <= 0 {
		maximumBytes = defaultSpoolBytes
	}
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return nil, err
	}
	if err := os.Chmod(directory, 0o700); err != nil {
		return nil, err
	}
	return &Spool{Dir: directory, MaxBytes: maximumBytes}, nil
}

// Enqueue atomically writes gzip JSON and rejects writes that exceed the
// configured aggregate capacity.
func (s *Spool) Enqueue(envelope FetchEnvelope) error {
	if envelope.FetchedAt.IsZero() {
		envelope.FetchedAt = time.Now().UTC()
	}
	encoded, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	digest := sha256.Sum256(encoded)
	name := envelope.FetchedAt.UTC().Format("20060102T150405.000000000Z") + "-" + hex.EncodeToString(digest[:8]) + ".json.gz"
	finalPath := filepath.Join(s.Dir, name)
	if _, err := os.Stat(finalPath); err == nil {
		return nil
	}
	temporary, err := os.CreateTemp(s.Dir, ".spool-*.tmp")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return err
	}
	writer := gzip.NewWriter(temporary)
	if _, err := writer.Write(encoded); err != nil {
		writer.Close()
		temporary.Close()
		return err
	}
	if err := writer.Close(); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return err
	}
	info, err := temporary.Stat()
	if err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	used, err := s.bytesUsed()
	if err != nil {
		return err
	}
	if used+info.Size() > s.MaxBytes {
		return &SpoolError{Code: "spool_full", Err: fmt.Errorf("spool would use %d of %d bytes", used+info.Size(), s.MaxBytes)}
	}
	if err := os.Rename(temporaryPath, finalPath); err != nil {
		return err
	}
	return syncDirectory(s.Dir)
}

// Replay processes files in lexical timestamp order. Corrupt files are
// isolated; persistence errors retain the current and subsequent files.
func (s *Spool) Replay(ctx context.Context, persister Persister) (ReplayReport, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		return ReplayReport{}, err
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json.gz") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	report := ReplayReport{}
	for _, name := range names {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		path := filepath.Join(s.Dir, name)
		envelope, err := readEnvelope(path)
		if err != nil {
			if quarantineErr := s.quarantine(path, name); quarantineErr != nil {
				return report, quarantineErr
			}
			report.Quarantined++
			continue
		}
		if err := persister.PersistFetch(ctx, envelope); err != nil {
			report.Failed++
			return report, err
		}
		if err := os.Remove(path); err != nil {
			return report, err
		}
		report.Persisted++
	}
	if err := syncDirectory(s.Dir); err != nil {
		return report, err
	}
	return report, nil
}

func readEnvelope(path string) (FetchEnvelope, error) {
	file, err := os.Open(path)
	if err != nil {
		return FetchEnvelope{}, err
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		return FetchEnvelope{}, err
	}
	defer reader.Close()
	decoder := json.NewDecoder(io.LimitReader(reader, 64<<20))
	var envelope FetchEnvelope
	if err := decoder.Decode(&envelope); err != nil {
		return FetchEnvelope{}, err
	}
	if envelope.SourceID == 0 || envelope.FetchedAt.IsZero() {
		return FetchEnvelope{}, fmt.Errorf("spool envelope is missing identity")
	}
	return envelope, nil
}

func (s *Spool) quarantine(path, name string) error {
	directory := filepath.Join(s.Dir, "quarantine")
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return err
	}
	target := filepath.Join(directory, name+".corrupt")
	return os.Rename(path, target)
}

func (s *Spool) bytesUsed() (int64, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		return 0, err
	}
	var total int64
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json.gz") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return 0, err
		}
		total += info.Size()
	}
	return total, nil
}

func syncDirectory(path string) error {
	directory, err := os.Open(path)
	if err != nil {
		return err
	}
	defer directory.Close()
	return directory.Sync()
}
