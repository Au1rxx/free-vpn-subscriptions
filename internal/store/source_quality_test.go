package store

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
)

func TestSourceQualityUpdatesAllowMoreThanTenThousandSources(t *testing.T) {
	updates := make([]SourceQualityUpdate, 10001)
	for index := range updates {
		updates[index] = SourceQualityUpdate{SourceID: uint64(index + 1), Score: 50}
	}
	if err := validateSourceQualityUpdates(updates); err != nil {
		t.Fatalf("large bounded update set rejected: %v", err)
	}
}

func TestSourceQualityRoundTripIntegration(t *testing.T) {
	configPath := os.Getenv("VPN_NODE_TEST_CONFIG")
	if configPath == "" {
		t.Skip("VPN_NODE_TEST_CONFIG is not set")
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	db, err := Open(ctx, cfg.Database, cfg.Database.Name)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	source, err := UpsertSource(ctx, db, SourceRecord{
		Name:    "source-quality-integration",
		URL:     fmt.Sprintf("https://example.invalid/source-quality/%d", time.Now().UnixNano()),
		Enabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = db.ExecContext(context.Background(), `DELETE FROM sources WHERE source_id=?`, source.ID)
	}()

	candidates, err := ListSourceQualityCandidates(ctx, db, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, candidate := range candidates {
		if candidate.SourceID == source.ID {
			found = true
			if candidate.FetchReliability != 0 || candidate.ParseYield != 0 || candidate.UsableNodeRate != 0 || candidate.LastSuccessAt.Valid {
				t.Fatalf("empty source candidate=%+v", candidate)
			}
		}
	}
	if !found {
		t.Fatalf("source %d missing from quality candidates", source.ID)
	}

	written, err := WriteSourceQualities(ctx, db, []SourceQualityUpdate{{SourceID: source.ID, Score: 73}})
	if err != nil {
		t.Fatal(err)
	}
	if written != 1 {
		t.Fatalf("written=%d, want 1", written)
	}
	var score float64
	if err := db.QueryRowContext(ctx, `SELECT quality_score FROM sources WHERE source_id=?`, source.ID).Scan(&score); err != nil {
		t.Fatal(err)
	}
	if score != 73 {
		t.Fatalf("score=%v, want 73", score)
	}
}
