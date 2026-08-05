package ingest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

func TestPersistPreparedFetchRecordsArchiveFallbackFailureWithoutLosingIdentity(t *testing.T) {
	write := store.FetchWrite{SourceID: 77, StatusCode: 200, Body: []byte("inline body"), FinishedAt: time.Unix(300, 0)}
	persistCalls := 0
	got, err := persistPreparedFetch(context.Background(), write, true, 1<<20,
		func(context.Context, store.FetchWrite) (store.FetchWrite, error) {
			return store.FetchWrite{}, errors.New("archive unavailable")
		},
		func(_ context.Context, candidate store.FetchWrite) error {
			persistCalls++
			if persistCalls == 1 {
				return &store.PayloadStorageError{Code: "database_payload_too_large", Err: errors.New("mediumblob")}
			}
			if candidate.SourceID != write.SourceID || candidate.ResponseBytes != int64(len(write.Body)) || len(candidate.Body) != 0 || candidate.ErrorCode != "payload_storage_required" {
				t.Fatalf("failure candidate=%+v", candidate)
			}
			return nil
		})
	if err != nil {
		t.Fatal(err)
	}
	if persistCalls != 2 || got.SourceID != write.SourceID || got.ErrorCode != "payload_storage_required" {
		t.Fatalf("calls=%d write=%+v", persistCalls, got)
	}
}

func TestPersistPreparedFetchReturnsExternalReferenceOnDatabaseOutage(t *testing.T) {
	write := store.FetchWrite{SourceID: 78, StatusCode: 200, Body: []byte("large body")}
	databaseErr := errors.New("database unavailable")
	got, err := persistPreparedFetch(context.Background(), write, true, 1,
		func(_ context.Context, candidate store.FetchWrite) (store.FetchWrite, error) {
			candidate.ResponseBytes = int64(len(candidate.Body))
			candidate.Body = nil
			candidate.ExternalPayload = &store.ExternalPayload{OriginalBytes: 10, CompressedBytes: 5, Compression: "gzip", ArchiveKey: "sha256/ab/cd/body.gz"}
			candidate.ExternalPayload.SHA256[0] = 1
			return candidate, nil
		},
		func(context.Context, store.FetchWrite) error { return databaseErr })
	if !errors.Is(err, databaseErr) || got.ExternalPayload == nil || len(got.Body) != 0 || got.SourceID != write.SourceID {
		t.Fatalf("write=%+v err=%v", got, err)
	}
}

func TestArchiveFetchWriteRemovesInlineBody(t *testing.T) {
	archive, err := NewPayloadArchive(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	body := bytes.Repeat([]byte("vless://large\n"), 128)
	write := store.FetchWrite{SourceID: 9, StatusCode: 200, Body: body, FinishedAt: time.Unix(200, 0)}
	archived, err := archiveFetchWrite(context.Background(), archive, write)
	if err != nil {
		t.Fatal(err)
	}
	if len(archived.Body) != 0 || archived.ExternalPayload == nil || archived.ResponseBytes != int64(len(body)) {
		t.Fatalf("archived write=%+v", archived)
	}
	if archived.ExternalPayload.OriginalBytes != int64(len(body)) || archived.ExternalPayload.ArchiveKey == "" {
		t.Fatalf("external payload=%+v", archived.ExternalPayload)
	}
	if !shouldArchiveBody(int64(len(body)), int64(len(body)-1)) {
		t.Fatal("body above threshold was not selected")
	}
	if shouldArchiveBody(int64(len(body)), int64(len(body))) {
		t.Fatal("body equal to threshold was selected")
	}
}

func TestIngestExternalEnvelopeRoundTripHasNoBody(t *testing.T) {
	payload := &store.ExternalPayload{OriginalBytes: 100, CompressedBytes: 20, Compression: "gzip", ArchiveKey: "sha256/ab/cd/body.gz"}
	payload.SHA256[0] = 1
	write := store.FetchWrite{SourceID: 8, FinishedAt: time.Unix(101, 0), StatusCode: 200, ExternalPayload: payload, ResponseBytes: 100}
	decoded := writeFromEnvelope(envelopeFromWrite(write))
	if len(decoded.Body) != 0 || decoded.ExternalPayload == nil || decoded.ExternalPayload.ArchiveKey != payload.ArchiveKey || decoded.ResponseBytes != 100 {
		t.Fatalf("external round trip mismatch: %+v", decoded)
	}
	encoded, err := json.Marshal(envelopeFromWrite(write))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(encoded, []byte(`"body"`)) || !bytes.Contains(encoded, []byte(`"archive_key"`)) || bytes.Contains(encoded, []byte(`"ArchiveKey"`)) {
		t.Fatalf("unstable external envelope=%s", encoded)
	}
}

func TestServiceReplaySpoolReportsEmptyQueue(t *testing.T) {
	spool, err := NewSpool(t.TempDir(), 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	service := Service{Spool: spool}
	report, err := service.ReplaySpool(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if report != (ReplayReport{}) {
		t.Fatalf("unexpected replay report: %+v", report)
	}
}

func TestIngestEnvelopeRoundTrip(t *testing.T) {
	write := store.FetchWrite{SourceID: 7, FinishedAt: time.Unix(100, 0), StatusCode: 200, FinalURL: "https://example.test/sub", Body: []byte("nodes"), Duration: 25 * time.Millisecond}
	decoded := writeFromEnvelope(envelopeFromWrite(write))
	if decoded.SourceID != write.SourceID || decoded.StatusCode != 200 || string(decoded.Body) != "nodes" || decoded.Duration != write.Duration {
		t.Fatalf("round trip mismatch: %+v", decoded)
	}
}

func TestParserVersionIncludesHTTPProxySemantics(t *testing.T) {
	if parserVersion != "fnctl-4" {
		t.Fatalf("parserVersion=%q, want fnctl-4", parserVersion)
	}
}
