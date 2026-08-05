package ingest

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type recordingPersister struct{ IDs []uint64 }

func (p *recordingPersister) PersistFetch(_ context.Context, envelope FetchEnvelope) error {
	p.IDs = append(p.IDs, envelope.SourceID)
	return nil
}

func TestSpoolEnqueueReplayOrderAndAtomicFiles(t *testing.T) {
	spool, err := NewSpool(t.TempDir(), 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []uint64{1, 2} {
		if err := spool.Enqueue(FetchEnvelope{SourceID: id, FetchedAt: time.Unix(int64(id), 0), Body: []byte("body")}); err != nil {
			t.Fatal(err)
		}
	}
	if matches, _ := filepath.Glob(filepath.Join(spool.Dir, "*.tmp")); len(matches) != 0 {
		t.Fatalf("temporary files remain: %v", matches)
	}
	persister := &recordingPersister{}
	report, err := spool.Replay(context.Background(), persister)
	if err != nil || report.Persisted != 2 || len(persister.IDs) != 2 || persister.IDs[0] != 1 || persister.IDs[1] != 2 {
		t.Fatalf("report=%+v ids=%v err=%v", report, persister.IDs, err)
	}
}

func TestSpoolQuarantinesCorruptionAndEnforcesLimit(t *testing.T) {
	spool, err := NewSpool(t.TempDir(), 100)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(spool.Dir, "000-bad.json.gz"), []byte("bad"), 0o600); err != nil {
		t.Fatal(err)
	}
	report, err := spool.Replay(context.Background(), &recordingPersister{})
	if err != nil || report.Quarantined != 1 {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	if err := spool.Enqueue(FetchEnvelope{SourceID: 3, FetchedAt: time.Now(), Body: make([]byte, 4096)}); spoolErrorCode(err) != "spool_full" {
		t.Fatalf("expected spool_full, got %v", err)
	}
}

func TestSpoolEnvelopeLimitCoversFiftyMiBBody(t *testing.T) {
	spool, err := NewSpool(t.TempDir(), 1<<20, 50<<20)
	if err != nil {
		t.Fatal(err)
	}
	if spool.MaxEnvelopeBytes <= 64<<20 {
		t.Fatalf("max envelope bytes=%d", spool.MaxEnvelopeBytes)
	}
}

func TestSpoolRetainsOversizedEnvelopeWithoutQuarantine(t *testing.T) {
	spool, err := NewSpool(t.TempDir(), 1<<20, 16)
	if err != nil {
		t.Fatal(err)
	}
	if err := spool.Enqueue(FetchEnvelope{SourceID: 5, FetchedAt: time.Now(), Body: make([]byte, 2<<20)}); err != nil {
		t.Fatal(err)
	}
	report, err := spool.Replay(context.Background(), &recordingPersister{})
	if spoolErrorCode(err) != "spool_envelope_too_large" || report.Quarantined != 0 || report.Failed != 1 {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	entries, globErr := filepath.Glob(filepath.Join(spool.Dir, "*.json.gz"))
	if globErr != nil || len(entries) != 1 {
		t.Fatalf("entries=%v err=%v", entries, globErr)
	}
}

func TestSpoolQuarantinesTrailingJSON(t *testing.T) {
	spool, err := NewSpool(t.TempDir(), 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(spool.Dir, "000-trailing.json.gz")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	writer := gzip.NewWriter(file)
	envelope := FetchEnvelope{SourceID: 6, FetchedAt: time.Now()}
	encoded, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write(append(encoded, []byte("{}")...)); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	report, err := spool.Replay(context.Background(), &recordingPersister{})
	if err != nil || report.Quarantined != 1 {
		t.Fatalf("report=%+v err=%v", report, err)
	}
}
