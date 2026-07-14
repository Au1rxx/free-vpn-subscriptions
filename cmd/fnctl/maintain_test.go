package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/maintain"
)

func TestWriteMaintenanceReportIncludesStorageMetricRows(t *testing.T) {
	report := maintain.Report{
		BeforeBytes:       100,
		AfterBytes:        90,
		Rows:              map[string]uint64{"raw_payloads": 0},
		Policy:            maintain.Policy{RawPayloadDays: 30},
		DryRun:            false,
		StorageMetricRows: 22,
	}
	var output bytes.Buffer
	if err := writeMaintenanceReport(&output, report); err != nil {
		t.Fatal(err)
	}
	got := output.String()
	for _, expected := range []string{
		"dry_run=false", "before_bytes=100", "after_bytes=90",
		"raw_ttl_days=30", `rows={"raw_payloads":0}`, "storage_metric_rows=22",
	} {
		if !strings.Contains(got, expected) {
			t.Fatalf("output %q missing %q", got, expected)
		}
	}
}
