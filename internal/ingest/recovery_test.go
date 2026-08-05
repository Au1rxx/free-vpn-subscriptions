package ingest

import (
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type recordingRecoveryPersister struct {
	existing map[[32]byte]bool
	writes   []FetchEnvelope
}

func (p *recordingRecoveryPersister) RecoveredFetchExists(_ context.Context, _ uint64, _ time.Time, digest [32]byte) (bool, error) {
	return p.existing[digest], nil
}

func (p *recordingRecoveryPersister) PersistRecoveredFetch(_ context.Context, envelope FetchEnvelope) error {
	p.writes = append(p.writes, envelope)
	return nil
}

func TestRecoverLegacyDirectoryNewestFirstKeepsOriginal(t *testing.T) {
	input := t.TempDir()
	older := writeLegacyEnvelope(t, input, "001.json.gz.corrupt", FetchEnvelope{
		SourceID: 555, FetchedAt: time.Unix(100, 0), StatusCode: 200, Body: []byte("older"),
	})
	newerEnvelope := FetchEnvelope{SourceID: 555, FetchedAt: time.Unix(200, 0), StatusCode: 200, Body: []byte("newer")}
	newer := writeLegacyEnvelope(t, input, "002.json.gz.corrupt", newerEnvelope)
	archive, err := NewPayloadArchive(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	persister := &recordingRecoveryPersister{existing: map[[32]byte]bool{}}

	results, err := RecoverLegacyDirectory(context.Background(), input, 555, 1, 4<<20, archive, persister)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Path != newer || results[0].Status != "recovered" || len(persister.writes) != 1 {
		t.Fatalf("results=%+v writes=%+v", results, persister.writes)
	}
	write := persister.writes[0]
	if len(write.Body) != 0 || write.ExternalPayload == nil || write.ExternalPayload.SHA256 != sha256.Sum256(newerEnvelope.Body) {
		t.Fatalf("recovered envelope=%+v", write)
	}
	for _, path := range []string{older, newer} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("original changed: %s: %v", path, err)
		}
	}
}

func TestRecoverLegacyDirectorySkipsExistingIdentity(t *testing.T) {
	input := t.TempDir()
	body := []byte("already recovered")
	path := writeLegacyEnvelope(t, input, "001.json.gz.corrupt", FetchEnvelope{
		SourceID: 555, FetchedAt: time.Unix(100, 0), StatusCode: 200, Body: body,
	})
	archive, err := NewPayloadArchive(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(body)
	persister := &recordingRecoveryPersister{existing: map[[32]byte]bool{digest: true}}
	results, err := RecoverLegacyDirectory(context.Background(), input, 555, 1, 4<<20, archive, persister)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 || len(persister.writes) != 0 {
		t.Fatalf("results=%+v writes=%+v", results, persister.writes)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("existing original changed: %v", err)
	}
}

func TestRecoverLegacyDirectoryLimitSkipsExistingNewestAndAdvances(t *testing.T) {
	input := t.TempDir()
	olderBody := []byte("older pending")
	older := writeLegacyEnvelope(t, input, "001.json.gz.corrupt", FetchEnvelope{
		SourceID: 555, FetchedAt: time.Unix(100, 0), StatusCode: 200, Body: olderBody,
	})
	newerBody := []byte("newer existing")
	writeLegacyEnvelope(t, input, "002.json.gz.corrupt", FetchEnvelope{
		SourceID: 555, FetchedAt: time.Unix(200, 0), StatusCode: 200, Body: newerBody,
	})
	archive, err := NewPayloadArchive(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	persister := &recordingRecoveryPersister{existing: map[[32]byte]bool{sha256.Sum256(newerBody): true}}
	results, err := RecoverLegacyDirectory(context.Background(), input, 555, 1, 4<<20, archive, persister)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Path != older || results[0].Status != "recovered" || len(persister.writes) != 1 {
		t.Fatalf("results=%+v writes=%+v", results, persister.writes)
	}
}

func TestRecoverLegacyDirectoryRejectsUnexpectedSourceWithStableCode(t *testing.T) {
	input := t.TempDir()
	writeLegacyEnvelope(t, input, "001.json.gz.corrupt", FetchEnvelope{
		SourceID: 556, FetchedAt: time.Unix(100, 0), StatusCode: 200, Body: []byte("wrong source"),
	})
	archive, err := NewPayloadArchive(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	_, err = RecoverLegacyDirectory(context.Background(), input, 555, 1, 4<<20, archive,
		&recordingRecoveryPersister{existing: map[[32]byte]bool{}})
	if recoveryErrorCode(err) != "recovery_source_mismatch" {
		t.Fatalf("err=%v", err)
	}
}

func writeLegacyEnvelope(t *testing.T, directory, name string, envelope FetchEnvelope) string {
	t.Helper()
	path := filepath.Join(directory, name)
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	writer := gzip.NewWriter(file)
	if err := json.NewEncoder(writer).Encode(envelope); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}
