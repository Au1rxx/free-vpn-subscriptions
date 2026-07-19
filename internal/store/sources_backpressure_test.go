package store

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
)

func TestClaimDueSourcesDefersSourceWithPendingParse(t *testing.T) {
	configPath := os.Getenv("VPN_NODE_TEST_CONFIG")
	if configPath == "" {
		t.Skip("VPN_NODE_TEST_CONFIG is not set")
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	admin, err := Open(ctx, cfg.Database, "")
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()

	database := fmt.Sprintf("vpn_nodes_backpressure_%d", time.Now().UnixNano())
	if _, err := admin.ExecContext(ctx, "CREATE DATABASE `"+database+"` CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), "DROP DATABASE IF EXISTS `"+database+"`")
	})
	db, err := Open(ctx, cfg.Database, database)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for _, statement := range []string{
		`CREATE TABLE sources (
			source_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			kind VARCHAR(32) NOT NULL, name VARCHAR(255) NOT NULL, url TEXT NOT NULL,
			canonical_url TEXT NOT NULL, format_hint VARCHAR(32) NOT NULL,
			protocol_hint VARCHAR(32) NULL, state VARCHAR(32) NOT NULL,
			depth TINYINT UNSIGNED NOT NULL, enabled BOOLEAN NOT NULL,
			priority SMALLINT NOT NULL, fetch_interval_seconds INT UNSIGNED NOT NULL,
			next_fetch_at DATETIME(6) NULL, etag VARCHAR(512) NULL,
			last_modified VARCHAR(255) NULL
		)`,
		`CREATE TABLE source_fetches (
			fetch_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			source_id BIGINT UNSIGNED NOT NULL, fetch_state VARCHAR(32) NOT NULL,
			parse_state VARCHAR(32) NOT NULL,
			KEY idx_source_fetches_source (source_id)
		)`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatal(err)
		}
	}
	result, err := db.ExecContext(ctx, `INSERT INTO sources
		(kind, name, url, canonical_url, format_hint, state, depth, enabled, priority,
		 fetch_interval_seconds, next_fetch_at)
		VALUES ('github-raw','backpressure','https://example.invalid/source',
		'https://example.invalid/source','auto','active',0,TRUE,100,900,UTC_TIMESTAMP(6))`)
	if err != nil {
		t.Fatal(err)
	}
	sourceID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO source_fetches
		(source_id, fetch_state, parse_state) VALUES (?, 'success', 'pending')`, sourceID); err != nil {
		t.Fatal(err)
	}

	claimed, err := ClaimDueSources(ctx, db, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(claimed) != 0 {
		t.Fatalf("claimed=%d, want 0 while source has a pending parse", len(claimed))
	}
	if _, err := db.ExecContext(ctx, `UPDATE source_fetches SET parse_state='success' WHERE source_id=?`, sourceID); err != nil {
		t.Fatal(err)
	}
	claimed, err = ClaimDueSources(ctx, db, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(claimed) != 1 || claimed[0].ID != uint64(sourceID) {
		t.Fatalf("claimed=%v, want source_id=%d after parse completion", claimed, sourceID)
	}
}
