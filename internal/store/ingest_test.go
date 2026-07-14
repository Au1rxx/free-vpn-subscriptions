package store

import (
	"bytes"
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/parse"
)

func TestCanonicalizeSourceURLIsStable(t *testing.T) {
	a, hashA, err := CanonicalizeSourceURL("HTTPS://Example.COM/path?b=2&a=1#fragment")
	if err != nil {
		t.Fatal(err)
	}
	b, hashB, err := CanonicalizeSourceURL("https://example.com/path?a=1&b=2")
	if err != nil {
		t.Fatal(err)
	}
	if a != b || hashA != hashB || a != "https://example.com/path?a=1&b=2" {
		t.Fatalf("canonical mismatch: %q %q", a, b)
	}
}

func TestPersistParseResultIntegration(t *testing.T) {
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
	db, err := Open(ctx, cfg.Database, cfg.Database.Name)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	suffix := strconv.FormatInt(time.Now().UnixNano(), 10)
	source, err := UpsertSource(ctx, db, SourceRecord{Name: "integration", URL: "https://example.invalid/" + suffix, Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	fetch, err := FinishFetch(ctx, db, FetchWrite{SourceID: source.ID, StatusCode: 200, Body: []byte("integration"), StartedAt: time.Now().UTC(), FinishedAt: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}
	nodes := []*node.Node{
		{Protocol: node.ProtoVLESS, Server: "integration-" + suffix + ".invalid", Port: 443, UUID: "a"},
		{Protocol: node.ProtoVLESS, Server: "integration-" + suffix + ".invalid", Port: 443, UUID: "b"},
	}
	defer func() {
		_, _ = db.ExecContext(context.Background(), `DELETE FROM sources WHERE source_id=?`, source.ID)
		_, _ = db.ExecContext(context.Background(), `DELETE FROM raw_payloads WHERE content_sha256=?`, fetch.PayloadHash[:])
		for _, n := range nodes {
			fingerprint := n.ConfigFingerprint()
			_, _ = db.ExecContext(context.Background(), `DELETE FROM node_configs WHERE config_fingerprint=?`, fingerprint[:])
		}
		host, fingerprint := endpointIdentity(nodes[0])
		_, _ = db.ExecContext(context.Background(), `DELETE FROM endpoints WHERE host=? AND host_hash=? AND port=?`, host, fingerprint[:], nodes[0].Port)
	}()
	report, err := PersistParseResult(ctx, db, source.ID, fetch.ID, parse.Result{Format: parse.FormatURIList, Nodes: nodes}, "integration-test")
	if err != nil {
		t.Fatal(err)
	}
	if report.NewEndpoints != 1 || report.NewConfigs != 2 || report.QueueJobs != 2 {
		t.Fatalf("unexpected persistence report: %+v", report)
	}
	firstFingerprint, secondFingerprint := nodes[0].ConfigFingerprint(), nodes[1].ConfigFingerprint()
	var missingExpiry int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM node_configs
		WHERE config_fingerprint IN (?,?) AND expires_at IS NULL`, firstFingerprint[:], secondFingerprint[:]).Scan(&missingExpiry); err != nil {
		t.Fatal(err)
	}
	if missingExpiry != 0 {
		t.Fatalf("persisted configs without expiry: %d", missingExpiry)
	}
	var seenBefore int
	if err := db.QueryRowContext(ctx, `SELECT seen_count FROM node_source_stats
		WHERE source_id=? AND node_config_id=(SELECT node_config_id FROM node_configs WHERE config_fingerprint=?)`,
		source.ID, firstFingerprint[:]).Scan(&seenBefore); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `UPDATE parse_runs SET parse_state='running', finished_at=NULL
		WHERE fetch_id=? AND parser_version='integration-test'`, fetch.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `UPDATE source_fetches SET parse_state='pending' WHERE fetch_id=?`, fetch.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := PersistParseResult(ctx, db, source.ID, fetch.ID, parse.Result{Format: parse.FormatURIList, Nodes: nodes}, "integration-test"); err != nil {
		t.Fatal(err)
	}
	var seenAfter int
	if err := db.QueryRowContext(ctx, `SELECT seen_count FROM node_source_stats
		WHERE source_id=? AND node_config_id=(SELECT node_config_id FROM node_configs WHERE config_fingerprint=?)`,
		source.ID, firstFingerprint[:]).Scan(&seenAfter); err != nil {
		t.Fatal(err)
	}
	if seenAfter != seenBefore {
		t.Fatalf("same fetch replay changed seen_count from %d to %d", seenBefore, seenAfter)
	}
}

func TestPersistParseResultCommitsNodeBatchesIndependentlyIntegration(t *testing.T) {
	configPath := os.Getenv("VPN_NODE_TEST_CONFIG")
	if configPath == "" {
		t.Skip("VPN_NODE_TEST_CONFIG is not set")
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	db, err := Open(ctx, cfg.Database, cfg.Database.Name)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	suffix := strconv.FormatInt(time.Now().UnixNano(), 10)
	hostPrefix := "batch-commit-" + suffix
	seedVersion := "batch-seed-" + suffix
	targetVersion := "batch-target-" + suffix
	seedSource, err := UpsertSource(ctx, db, SourceRecord{Name: seedVersion, URL: "https://example.invalid/" + seedVersion, Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	seedFetch, err := FinishFetch(ctx, db, FetchWrite{SourceID: seedSource.ID, StatusCode: 200, Body: []byte(seedVersion), StartedAt: time.Now().UTC(), FinishedAt: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}
	lockedNode := &node.Node{Protocol: node.ProtoVLESS, Server: hostPrefix + "-locked.invalid", Port: 443, UUID: "locked"}
	if _, err := PersistParseResult(ctx, db, seedSource.ID, seedFetch.ID, parse.Result{Format: parse.FormatURIList, Nodes: []*node.Node{lockedNode}}, seedVersion); err != nil {
		t.Fatal(err)
	}

	targetSource, err := UpsertSource(ctx, db, SourceRecord{Name: targetVersion, URL: "https://example.invalid/" + targetVersion, Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	targetFetch, err := FinishFetch(ctx, db, FetchWrite{SourceID: targetSource.ID, StatusCode: 200, Body: []byte(targetVersion), StartedAt: time.Now().UTC(), FinishedAt: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}
	nodes := make([]*node.Node, 0, nodePersistBatchSize+1)
	for i := 0; i < nodePersistBatchSize; i++ {
		nodes = append(nodes, &node.Node{Protocol: node.ProtoVLESS, Server: hostPrefix + "-new.invalid", Port: 443, UUID: strconv.Itoa(i)})
	}
	nodes = append(nodes, lockedNode)

	defer func() {
		_, _ = db.ExecContext(context.Background(), `DELETE FROM sources WHERE source_id IN (?,?)`, seedSource.ID, targetSource.ID)
		_, _ = db.ExecContext(context.Background(), `DELETE FROM node_configs WHERE parser_version IN (?,?)`, seedVersion, targetVersion)
		_, _ = db.ExecContext(context.Background(), `DELETE FROM endpoints WHERE host LIKE ?`, hostPrefix+"%")
	}()

	_, lockedHash := endpointIdentity(lockedNode)
	lockTx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	var lockedEndpointID uint64
	if err := lockTx.QueryRowContext(ctx, `SELECT endpoint_id FROM endpoints WHERE host_hash=? AND port=? FOR UPDATE`, lockedHash[:], lockedNode.Port).Scan(&lockedEndpointID); err != nil {
		_ = lockTx.Rollback()
		t.Fatal(err)
	}

	result := make(chan error, 1)
	go func() {
		_, persistErr := PersistParseResult(ctx, db, targetSource.ID, targetFetch.ID,
			parse.Result{Format: parse.FormatURIList, Nodes: nodes}, targetVersion)
		result <- persistErr
	}()

	visible := false
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		var count int
		if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM node_configs WHERE parser_version=?`, targetVersion).Scan(&count); err != nil {
			_ = lockTx.Rollback()
			t.Fatal(err)
		}
		if count == nodePersistBatchSize {
			visible = true
			break
		}
		select {
		case err := <-result:
			_ = lockTx.Rollback()
			t.Fatalf("parse finished before the locked second batch: %v", err)
		case <-time.After(100 * time.Millisecond):
		}
	}
	if err := lockTx.Rollback(); err != nil {
		t.Fatal(err)
	}
	if err := <-result; err != nil {
		t.Fatal(err)
	}
	if !visible {
		t.Fatal("first node batch was not visible while the second batch was locked")
	}
}

func TestPayloadCompressionRoundTrip(t *testing.T) {
	body := bytes.Repeat([]byte("nodes\n"), 1000)
	compressed, err := compressPayload(body)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := decompressPayload(compressed, int64(len(body)))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decoded, body) || len(compressed) >= len(body) {
		t.Fatalf("roundtrip bytes=%d compressed=%d", len(decoded), len(compressed))
	}
}

func TestNodeIdentityKeepsSameEndpointDifferentConfigs(t *testing.T) {
	a := &node.Node{Protocol: node.ProtoVLESS, Server: "example.com", Port: 443, UUID: "a"}
	b := &node.Node{Protocol: node.ProtoVLESS, Server: "EXAMPLE.com.", Port: 443, UUID: "b"}
	hostA, endpointA := endpointIdentity(a)
	hostB, endpointB := endpointIdentity(b)
	if hostA != hostB || endpointA != endpointB {
		t.Fatal("same endpoint identity differs")
	}
	if a.ConfigFingerprint() == b.ConfigFingerprint() {
		t.Fatal("different configurations collapsed")
	}
}

func TestNodeExpiryIsThirtyDaysAfterLastObservation(t *testing.T) {
	seen := time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC)
	if got, want := nodeExpiresAt(seen), seen.Add(30*24*time.Hour); !got.Equal(want) {
		t.Fatalf("node expiry=%s, want %s", got, want)
	}
}
