package ingest

import (
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

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
