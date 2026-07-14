package main

import (
	"strings"
	"testing"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

func TestRootContainsDatabaseCommands(t *testing.T) {
	root := newRootCmd()
	want := map[string]bool{
		"aggregate": false, "migrate": false, "db-status": false,
		"import-seeds": false, "fetch": false, "parse": false, "discover": false, "ingest-status": false,
		"validate-worker":   false,
		"validation-status": false,
		"classify":          false,
		"maintain":          false,
		"export-db":         false,
		"prune-discovery":   false,
		"requeue-parses":    false,
	}
	for _, command := range root.Commands() {
		if _, ok := want[command.Name()]; ok {
			want[command.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("missing command %s", name)
		}
	}
}

func TestFormatDatabaseStatusIsBoundedAndContainsNoDSN(t *testing.T) {
	status := store.DatabaseStatus{
		Server: store.ServerInfo{
			Version: "9.7.1-cloud", Cipher: "TLS_AES_128_GCM_SHA256",
			TimeZone: "UTC", Charset: "utf8mb4", Collation: "utf8mb4_0900_ai_ci",
		},
		AppliedMigrations:   6,
		BusinessTables:      22,
		EmptyTableComments:  0,
		EmptyColumnComments: 0,
		EnabledPolicies:     6,
		DataBytes:           1024,
		IndexBytes:          2048,
		AllocatedBytes:      4096,
	}
	out := formatDatabaseStatus(status)
	for _, want := range []string{
		"version=9.7.1-cloud", "tls=TLS_AES_128_GCM_SHA256", "migrations=6", "tables=22",
		"empty_table_comments=0", "empty_column_comments=0", "enabled_policies=6",
		"allocated_bytes=4096", "total_bytes=4096",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("status missing %q: %s", want, out)
		}
	}
	if strings.Contains(out, "@tcp(") || strings.Contains(out, "password") {
		t.Fatalf("status leaked connection details: %s", out)
	}
	if lines := strings.Count(out, "\n"); lines > 10 {
		t.Fatalf("status output too long: %d lines", lines)
	}
}

func TestFormatValidationStatusIncludesPerformanceCoverage(t *testing.T) {
	status := store.ValidationStatus{Batches: 10, Attempts: 20, CurrentStatuses: 18,
		PendingJobs: 5, LeasedJobs: 2, ExpiredLeases: 1, Passed: 8, Partial: 2, Failed: 10,
		Available: 8, Available24H: 6, Degraded: 2, Unavailable: 10,
		EligiblePendingJobs: 4, OldestPendingAgeSeconds: 3600, PerformanceAttempts: 7,
		PerformanceSuccesses: 6, AverageBytesPerSecond: 1048576,
		ScoredNodes: 10, AverageQualityScore: 82, ByGrade: map[string]uint64{"A": 6, "B": 4}}
	out := formatValidationStatus(status)
	for _, want := range []string{"attempts=20", "available_24h=6", "eligible_pending_jobs=4", "oldest_pending_age_seconds=3600", "performance_attempts=7", "performance_successes=6", "average_bytes_per_second=1048576", "scored_nodes=10", "average_quality_score=82", "grade_A=6", "grade_B=4", "grade_U=0"} {
		if !strings.Contains(out, want) {
			t.Errorf("validation status missing %q: %s", want, out)
		}
	}
	if lines := strings.Count(out, "\n"); lines > 5 {
		t.Fatalf("validation status output too long: %d lines", lines)
	}
}

func TestFormatIngestStatusIncludesSourceHealth(t *testing.T) {
	status := store.IngestStatus{Sources: 65, EnabledSources: 63, Fetches24H: 60,
		SuccessfulFetches24H: 57, FailedFetches24H: 3, ByProtocol: map[string]uint64{},
		ScoredSources: 61, AverageSourceQuality: 72, MaximumSourceQuality: 96,
		SourceKinds:         map[string]store.SourceKindStatus{"github-raw": {Total: 40, Enabled: 39}},
		FetchErrorCounts24H: []store.NamedCount{{Name: "http_status_http_404", Count: 2}},
		ParseErrorCounts:    []store.NamedCount{{Name: "invalid_uri", Count: 7}}}
	out := formatIngestStatus(status)
	for _, want := range []string{"sources=65", "enabled_sources=63", "fetches_24h=60", "successful_fetches_24h=57", "failed_fetches_24h=3", "source_quality scored_sources=61 average=72 maximum=96", "source_kind=github-raw total=40 enabled=39", "fetch_error=http_status_http_404 count=2", "parse_error=invalid_uri count=7"} {
		if !strings.Contains(out, want) {
			t.Errorf("ingest status missing %q: %s", want, out)
		}
	}
}
