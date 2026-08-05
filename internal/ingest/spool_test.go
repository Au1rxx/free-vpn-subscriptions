package ingest

import (
	"context"
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
