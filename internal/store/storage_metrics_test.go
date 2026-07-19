package store

import (
	"context"
	"math"
	"os"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
)

func TestRecordStorageMetricsIntegration(t *testing.T) {
	configPath := os.Getenv("VPN_NODE_TEST_CONFIG")
	if configPath == "" {
		t.Skip("VPN_NODE_TEST_CONFIG is not set")
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	db, err := Open(ctx, cfg.Database, cfg.Database.Name)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sampledAt := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Microsecond)
	var rowsBefore int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM storage_metrics`).Scan(&rowsBefore); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `DELETE FROM storage_metrics WHERE sampled_at=?`, sampledAt); err != nil {
		t.Fatal(err)
	}
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if _, cleanupErr := db.ExecContext(cleanupCtx, `DELETE FROM storage_metrics WHERE sampled_at=?`, sampledAt); cleanupErr != nil {
			t.Errorf("cleanup storage metrics: %v", cleanupErr)
			return
		}
		var rowsAfter int
		if cleanupErr := db.QueryRowContext(cleanupCtx, `SELECT COUNT(*) FROM storage_metrics`).Scan(&rowsAfter); cleanupErr != nil {
			t.Errorf("count storage metrics after cleanup: %v", cleanupErr)
		} else if rowsAfter != rowsBefore {
			t.Errorf("storage metric rows after cleanup=%d, want %d", rowsAfter, rowsBefore)
		}
	}()

	var tableCount int64
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema=DATABASE()`).Scan(&tableCount); err != nil {
		t.Fatal(err)
	}
	written, err := RecordStorageMetrics(ctx, db, sampledAt, DatabaseCapacityBytes, 5<<30)
	if err != nil {
		t.Fatal(err)
	}
	if written != tableCount {
		t.Fatalf("written=%d, want table count %d", written, tableCount)
	}

	var persisted int64
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM storage_metrics WHERE sampled_at=?`, sampledAt).Scan(&persisted); err != nil {
		t.Fatal(err)
	}
	if persisted != tableCount {
		t.Fatalf("persisted=%d, want table count %d", persisted, tableCount)
	}

	var schema, table string
	var dataBytes, totalBytes, capacity uint64
	var usage float64
	if err := db.QueryRowContext(ctx, `SELECT table_schema,table_name,data_bytes,total_bytes,capacity_bytes,usage_percent
		FROM storage_metrics WHERE sampled_at=? AND table_name='node_configs'`, sampledAt).Scan(
		&schema, &table, &dataBytes, &totalBytes, &capacity, &usage); err != nil {
		t.Fatal(err)
	}
	if schema != cfg.Database.Name || table != "node_configs" {
		t.Fatalf("sample identity=%s.%s, want %s.node_configs", schema, table, cfg.Database.Name)
	}
	if dataBytes == 0 || totalBytes == 0 {
		t.Fatalf("node_configs bytes data=%d total=%d, want non-zero", dataBytes, totalBytes)
	}
	if capacity != DatabaseCapacityBytes || math.Abs(usage-10) > 0.0001 {
		t.Fatalf("capacity=%d usage=%v, want %d and 10.000", capacity, usage, DatabaseCapacityBytes)
	}

	if _, err := RecordStorageMetrics(ctx, db, sampledAt, DatabaseCapacityBytes, 10<<30); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*),MAX(usage_percent) FROM storage_metrics WHERE sampled_at=?`, sampledAt).Scan(&persisted, &usage); err != nil {
		t.Fatal(err)
	}
	if persisted != tableCount || math.Abs(usage-20) > 0.0001 {
		t.Fatalf("idempotent rows=%d usage=%v, want %d and 20.000", persisted, usage, tableCount)
	}

	if _, err := RecordStorageMetrics(ctx, nil, sampledAt, 0, 1); err == nil {
		t.Fatal("zero capacity was accepted")
	}
}
