package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
)

func TestValidationStatusQueriesUseOneAggregatePerTable(t *testing.T) {
	if len(validationStatusQueries) != 4 {
		t.Fatalf("query count=%d, want 4", len(validationStatusQueries))
	}
	combined := strings.ToLower(strings.Join(validationStatusQueries, "\n"))
	if strings.Contains(combined, "(select") {
		t.Fatal("validation status queries must not use scalar subqueries")
	}
	for _, table := range []string{"validation_batches", "validation_queue", "validation_attempts", "node_current_status"} {
		if count := strings.Count(combined, "from "+table); count != 1 {
			t.Errorf("table %s is read %d times, want 1", table, count)
		}
	}
	if !strings.Contains(combined, "from validation_attempts force index (idx_validation_attempts_status)") {
		t.Fatal("validation attempt aggregate must use its covering index")
	}
}

func TestReadValidationStatusAggregatesEachTable(t *testing.T) {
	db, ctx := openValidationStatusTestDatabase(t)
	statements := []string{
		`INSERT INTO validation_batches VALUES (1), (2)`,
		`INSERT INTO validation_queue (job_state,next_attempt_at,leased_until) VALUES
			('pending', UTC_TIMESTAMP(6)-INTERVAL 2 MINUTE, NULL),
			('pending', UTC_TIMESTAMP(6)+INTERVAL 1 HOUR, NULL),
			('leased', UTC_TIMESTAMP(6), UTC_TIMESTAMP(6)-INTERVAL 1 MINUTE),
			('leased', UTC_TIMESTAMP(6), UTC_TIMESTAMP(6)+INTERVAL 1 HOUR)`,
		`INSERT INTO validation_attempts
			(passed,partial_success,performance_bytes,performance_error_code,bytes_per_second) VALUES
			(TRUE,FALSE,1000,NULL,100),
			(FALSE,TRUE,0,'speed_failed',0),
			(FALSE,FALSE,NULL,NULL,NULL)`,
		`INSERT INTO node_current_status
			(availability_state,last_validation_at,quality_score,quality_grade) VALUES
			('available',UTC_TIMESTAMP(6),85,'A'),
			('available',UTC_TIMESTAMP(6)-INTERVAL 25 HOUR,95,'S'),
			('degraded',UTC_TIMESTAMP(6),55,'C'),
			('unavailable',NULL,10,'D')`,
	}
	for _, statement := range statements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatal(err)
		}
	}

	status, err := ReadValidationStatus(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if status.Batches != 2 || status.Attempts != 3 || status.CurrentStatuses != 4 {
		t.Fatalf("top-level counts: %+v", status)
	}
	if status.PendingJobs != 2 || status.EligiblePendingJobs != 1 || status.LeasedJobs != 2 || status.ExpiredLeases != 1 {
		t.Fatalf("queue counts: %+v", status)
	}
	if status.OldestPendingAgeSeconds < 120 || status.OldestPendingAgeSeconds > 300 {
		t.Fatalf("oldest pending age=%d, want 120..300", status.OldestPendingAgeSeconds)
	}
	if status.Passed != 1 || status.Partial != 1 || status.Failed != 1 {
		t.Fatalf("attempt outcomes: %+v", status)
	}
	if status.PerformanceAttempts != 2 || status.PerformanceSuccesses != 1 || status.AverageBytesPerSecond != 100 {
		t.Fatalf("performance counts: %+v", status)
	}
	if status.Available != 2 || status.Available24H != 1 || status.Degraded != 1 || status.Unavailable != 1 {
		t.Fatalf("availability counts: %+v", status)
	}
	if status.ScoredNodes != 3 || status.AverageQualityScore != 78 {
		t.Fatalf("score counts: %+v", status)
	}
	for _, grade := range []string{"S", "A", "C", "D"} {
		if status.ByGrade[grade] != 1 {
			t.Errorf("grade %s=%d, want 1", grade, status.ByGrade[grade])
		}
	}
	for _, grade := range []string{"B", "U"} {
		if status.ByGrade[grade] != 0 {
			t.Errorf("grade %s=%d, want 0", grade, status.ByGrade[grade])
		}
	}
}

func TestReadValidationStatusRunsIndependentAggregatesConcurrently(t *testing.T) {
	db, ctx := openValidationStatusTestDatabase(t)
	original := validationStatusQueries
	validationStatusQueries = []string{
		`SELECT 0 FROM (SELECT SLEEP(0.5)) AS delay_row`,
		`SELECT 0,0,0,0,0 FROM (SELECT SLEEP(0.5)) AS delay_row`,
		`SELECT 0,0,0,0,0,0,0 FROM (SELECT SLEEP(0.5)) AS delay_row`,
		`SELECT 0,0,0,0,0,0,0,0,0,0,0,0,0 FROM (SELECT SLEEP(0.5)) AS delay_row`,
	}
	t.Cleanup(func() { validationStatusQueries = original })

	started := time.Now()
	if _, err := ReadValidationStatus(ctx, db); err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(started); elapsed >= 1500*time.Millisecond {
		t.Fatalf("independent aggregates took %s, want less than 1.5s", elapsed)
	}
}

func openValidationStatusTestDatabase(t *testing.T) (*sql.DB, context.Context) {
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
	database := fmt.Sprintf("vpn_nodes_validation_status_%d", time.Now().UnixNano())
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
		`CREATE TABLE validation_batches (validation_batch_id BIGINT UNSIGNED PRIMARY KEY)`,
		`CREATE TABLE validation_queue (
			validation_job_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			job_state VARCHAR(32) NOT NULL,
			next_attempt_at DATETIME(6) NOT NULL,
			leased_until DATETIME(6) NULL
		)`,
		`CREATE TABLE validation_attempts (
			validation_attempt_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			passed BOOLEAN NOT NULL,
			partial_success BOOLEAN NOT NULL,
			performance_bytes BIGINT UNSIGNED NULL,
			performance_error_code VARCHAR(64) NULL,
			bytes_per_second BIGINT UNSIGNED NULL,
			KEY idx_validation_attempts_status (passed,partial_success,performance_bytes,performance_error_code,bytes_per_second)
		)`,
		`CREATE TABLE node_current_status (
			node_config_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			availability_state VARCHAR(32) NOT NULL,
			last_validation_at DATETIME(6) NULL,
			quality_score TINYINT UNSIGNED NOT NULL,
			quality_grade CHAR(1) NOT NULL
		)`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatal(err)
		}
	}
	return db, ctx
}
