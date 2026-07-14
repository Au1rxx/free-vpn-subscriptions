package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

func TestValidationQueueRejectsInvalidClaim(t *testing.T) {
	for _, test := range []struct {
		owner string
		limit int
		lease time.Duration
	}{{"", 1, time.Minute}, {"worker", 0, time.Minute}, {"worker", 1, 0}} {
		if _, err := ClaimValidationJobs(context.Background(), nil, test.owner, test.limit, test.lease); err == nil {
			t.Fatalf("expected validation error for %+v", test)
		}
	}
}

func TestValidationQueueOwnershipErrorIsStable(t *testing.T) {
	if !errors.Is(ErrLeaseOwnership, ErrLeaseOwnership) {
		t.Fatal("lease ownership sentinel is not stable")
	}
}

func TestClaimInitialValidationJobsByProtocolFiltersEligibility(t *testing.T) {
	db, ctx := openValidationQueueTestDatabase(t)
	eligibleA := insertValidationQueueTestJob(t, ctx, db, node.ProtoHTTP, false, "pending", 0, time.Now().Add(-time.Minute))
	eligibleB := insertValidationQueueTestJob(t, ctx, db, node.ProtoHTTP, false, "pending", 0, time.Now().Add(-time.Minute))
	insertValidationQueueTestJob(t, ctx, db, node.ProtoHTTPS, false, "pending", 0, time.Now().Add(-time.Minute))
	insertValidationQueueTestJob(t, ctx, db, node.ProtoHTTP, false, "pending", 1, time.Now().Add(-time.Minute))
	insertValidationQueueTestJob(t, ctx, db, node.ProtoHTTP, false, "pending", 0, time.Now().Add(time.Hour))
	insertValidationQueueTestJob(t, ctx, db, node.ProtoHTTP, false, "leased", 0, time.Now().Add(-time.Minute))
	insertValidationQueueTestJob(t, ctx, db, node.ProtoHTTP, true, "pending", 0, time.Now().Add(-time.Minute))

	jobs, err := ClaimInitialValidationJobsByProtocol(ctx, db, "initial-http-a", node.ProtoHTTP, 10, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 2 {
		t.Fatalf("claimed=%d, want 2", len(jobs))
	}
	want := map[uint64]bool{eligibleA: true, eligibleB: true}
	for _, job := range jobs {
		if !want[job.ID] {
			t.Errorf("claimed ineligible job_id=%d", job.ID)
		}
		if job.Protocol != node.ProtoHTTP || job.State != "leased" || job.LeaseOwner != "initial-http-a" || job.Attempts != 1 {
			t.Errorf("unexpected claimed job: %+v", job)
		}
	}
	var leased int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM validation_queue
		WHERE lease_owner='initial-http-a' AND job_state='leased' AND attempts=1`).Scan(&leased); err != nil {
		t.Fatal(err)
	}
	if leased != 2 {
		t.Fatalf("persisted leased jobs=%d, want 2", leased)
	}
}

func TestClaimInitialValidationJobsByProtocolDoesNotDuplicateConcurrentClaims(t *testing.T) {
	db, ctx := openValidationQueueTestDatabase(t)
	for range 4 {
		insertValidationQueueTestJob(t, ctx, db, node.ProtoSOCKS5, false, "pending", 0, time.Now().Add(-time.Minute))
	}

	start := make(chan struct{})
	results := make(chan []ValidationJob, 2)
	errors := make(chan error, 2)
	var wait sync.WaitGroup
	for _, owner := range []string{"initial-socks5-a", "initial-socks5-b"} {
		owner := owner
		wait.Add(1)
		go func() {
			defer wait.Done()
			<-start
			jobs, err := ClaimInitialValidationJobsByProtocol(ctx, db, owner, node.ProtoSOCKS5, 2, time.Minute)
			if err != nil {
				errors <- err
				return
			}
			results <- jobs
		}()
	}
	close(start)
	wait.Wait()
	close(results)
	close(errors)
	for err := range errors {
		t.Fatal(err)
	}

	claimed := make(map[uint64]bool)
	for jobs := range results {
		if len(jobs) != 2 {
			t.Fatalf("one owner claimed=%d, want 2", len(jobs))
		}
		for _, job := range jobs {
			if claimed[job.ID] {
				t.Fatalf("job_id=%d was claimed twice", job.ID)
			}
			claimed[job.ID] = true
		}
	}
	if len(claimed) != 4 {
		t.Fatalf("unique claimed jobs=%d, want 4", len(claimed))
	}
}

func openValidationQueueTestDatabase(t *testing.T) (*sql.DB, context.Context) {
	t.Helper()
	configPath := os.Getenv("VPN_NODE_TEST_CONFIG")
	if configPath == "" {
		t.Skip("VPN_NODE_TEST_CONFIG is not set")
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	admin, err := Open(ctx, cfg.Database, "")
	if err != nil {
		t.Fatal(err)
	}
	database := fmt.Sprintf("vpn_nodes_validation_queue_%d", time.Now().UnixNano())
	if _, err := admin.ExecContext(ctx, "CREATE DATABASE `"+database+"` CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci"); err != nil {
		admin.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		_, _ = admin.ExecContext(cleanupCtx, "DROP DATABASE IF EXISTS `"+database+"`")
		_ = admin.Close()
	})
	db, err := Open(ctx, cfg.Database, database)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	for _, statement := range []string{
		`CREATE TABLE node_configs (
			node_config_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			protocol VARCHAR(32) NOT NULL,
			normalized_config JSON NOT NULL,
			is_exportable BOOLEAN NOT NULL DEFAULT FALSE,
			last_success_at DATETIME(6) NULL,
			KEY idx_node_configs_export (is_exportable, protocol, last_success_at)
		)`,
		`CREATE TABLE validation_queue (
			validation_job_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			node_config_id BIGINT UNSIGNED NOT NULL,
			stage VARCHAR(32) NOT NULL,
			priority SMALLINT NOT NULL,
			job_state VARCHAR(32) NOT NULL,
			attempts INT UNSIGNED NOT NULL,
			lease_owner VARCHAR(255) NULL,
			next_attempt_at DATETIME(6) NOT NULL,
			leased_until DATETIME(6) NULL
		)`,
		`CREATE TABLE validation_attempts (
			validation_attempt_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			node_config_id BIGINT UNSIGNED NOT NULL,
			performance_bytes BIGINT UNSIGNED NULL,
			started_at DATETIME(6) NOT NULL
		)`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatal(err)
		}
	}
	return db, ctx
}

func insertValidationQueueTestJob(t *testing.T, ctx context.Context, db *sql.DB, protocol string, exportable bool, state string, attempts int, nextAttempt time.Time) uint64 {
	t.Helper()
	result, err := db.ExecContext(ctx, `INSERT INTO node_configs
		(protocol, normalized_config, is_exportable)
		VALUES (?, JSON_OBJECT('protocol', ?, 'server', 'test.invalid', 'port', 443), ?)`, protocol, protocol, exportable)
	if err != nil {
		t.Fatal(err)
	}
	nodeConfigID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	result, err = db.ExecContext(ctx, `INSERT INTO validation_queue
		(node_config_id, stage, priority, job_state, attempts, next_attempt_at, leased_until)
		VALUES (?, 'connectivity', 100, ?, ?, ?, CASE WHEN ?='leased' THEN DATE_ADD(UTC_TIMESTAMP(6), INTERVAL 1 HOUR) ELSE NULL END)`,
		nodeConfigID, state, attempts, nextAttempt.UTC(), state)
	if err != nil {
		t.Fatal(err)
	}
	jobID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	return uint64(jobID)
}
