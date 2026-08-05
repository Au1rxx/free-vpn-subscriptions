package ingest

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type RecoveryResult struct {
	Path, ArchiveKey, Status string
	SourceID                 uint64
}

type RecoveryError struct {
	Code string
	Err  error
}

func (e *RecoveryError) Error() string { return e.Code + ": " + e.Err.Error() }
func (e *RecoveryError) Unwrap() error { return e.Err }

func recoveryErrorCode(err error) string {
	if typed, ok := err.(*RecoveryError); ok {
		return typed.Code
	}
	return ""
}

type RecoveryPersister interface {
	RecoveredFetchExists(context.Context, uint64, time.Time, [32]byte) (bool, error)
	PersistRecoveredFetch(context.Context, FetchEnvelope) error
}

func RecoverLegacyDirectory(ctx context.Context, inputDirectory string, sourceID uint64, limit int,
	maximumEnvelopeBytes int64, archive *PayloadArchive, persister RecoveryPersister, newestFirst ...bool) ([]RecoveryResult, error) {
	if strings.TrimSpace(inputDirectory) == "" || sourceID == 0 {
		return nil, fmt.Errorf("recovery input directory and source id are required")
	}
	if limit < 1 || limit > 100 {
		return nil, fmt.Errorf("recovery limit must be between 1 and 100")
	}
	if maximumEnvelopeBytes < 1 || archive == nil || persister == nil {
		return nil, fmt.Errorf("recovery dependencies are invalid")
	}
	entries, err := os.ReadDir(inputDirectory)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".json.gz") || strings.HasSuffix(name, ".json.gz.corrupt") {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	if len(newestFirst) == 0 || newestFirst[0] {
		sort.Sort(sort.Reverse(sort.StringSlice(names)))
	}
	results := make([]RecoveryResult, 0, len(names))
	for _, name := range names {
		if len(results) >= limit {
			break
		}
		if err := ctx.Err(); err != nil {
			return results, err
		}
		path := filepath.Join(inputDirectory, name)
		envelope, err := readEnvelope(path, maximumEnvelopeBytes)
		if err != nil {
			return results, fmt.Errorf("read legacy envelope %s: %w", name, err)
		}
		if envelope.SourceID != sourceID {
			return results, &RecoveryError{Code: "recovery_source_mismatch", Err: fmt.Errorf("legacy envelope %s source id=%d want=%d", name, envelope.SourceID, sourceID)}
		}
		if len(envelope.Body) == 0 {
			return results, fmt.Errorf("legacy envelope %s has no body", name)
		}
		digest := sha256.Sum256(envelope.Body)
		exists, err := persister.RecoveredFetchExists(ctx, sourceID, envelope.FetchedAt, digest)
		if err != nil {
			return results, err
		}
		if exists {
			continue
		}
		archived, err := archiveFetchWrite(ctx, archive, writeFromEnvelope(envelope))
		if err != nil {
			return results, fmt.Errorf("archive legacy envelope %s: %w", name, err)
		}
		recovered := envelopeFromWrite(archived)
		if err := persister.PersistRecoveredFetch(ctx, recovered); err != nil {
			return results, err
		}
		results = append(results, RecoveryResult{
			Path: path, ArchiveKey: archived.ExternalPayload.ArchiveKey, Status: "recovered", SourceID: sourceID,
		})
	}
	return results, nil
}
