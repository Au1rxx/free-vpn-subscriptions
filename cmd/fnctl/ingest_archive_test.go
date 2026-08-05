package main

import (
	"strings"
	"testing"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/ingest"
)

func TestIngestSpoolCompatibilityLimitUsesFetchCeiling(t *testing.T) {
	cfg := &config.Config{}
	cfg.Fetch.MaxBodyBytes = 192 << 20
	cfg.Fetch.ArchiveThresholdBytes = 50 << 20
	if got := ingestSpoolBodyLimit(cfg); got != 192<<20 {
		t.Fatalf("spool body limit=%d", got)
	}
}

func TestFormatArchiveStatusIsBounded(t *testing.T) {
	output := formatArchiveStatus(ingest.ArchiveStatus{
		Files: 10, Referenced: 8, Missing: 1, Corrupt: 1, Unreferenced: 2, Bytes: 4096,
		QuarantineFiles: 158, QuarantineBytes: 4300000000,
	})
	for _, want := range []string{
		"files=10", "referenced=8", "missing=1", "corrupt=1", "unreferenced=2", "bytes=4096",
		"quarantine_files=158", "quarantine_bytes=4300000000",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q: %s", want, output)
		}
	}
	if strings.Count(output, "\n") != 1 || strings.Contains(output, "/var/") || strings.Contains(output, "http") {
		t.Fatalf("unbounded or sensitive output: %q", output)
	}
}

func TestRecoverLegacySpoolCommandDefaultsToOneNewestFile(t *testing.T) {
	command := newRecoverLegacySpoolCmd()
	limit, err := command.Flags().GetInt("limit")
	if err != nil {
		t.Fatal(err)
	}
	newest, err := command.Flags().GetBool("newest-first")
	if err != nil {
		t.Fatal(err)
	}
	if limit != 1 || !newest {
		t.Fatalf("limit=%d newest=%t", limit, newest)
	}
}
