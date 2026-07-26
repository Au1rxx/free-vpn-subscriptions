package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/aggregate"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/pages"
)

func TestPreparePageHistoryDegradesToCurrentSnapshotAndPreservesCorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "network-history.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	const corrupt = `{"schema_version":`
	if err := os.WriteFile(path, []byte(corrupt), 0o644); err != nil {
		t.Fatal(err)
	}
	summary := aggregate.Summary{
		TotalSelected:   2000,
		TotalVerified:   7866,
		MedianLatencyMS: 584,
		GeneratedAtUnix: time.Date(2026, 7, 26, 4, 0, 0, 0, time.UTC).Unix(),
		ByCountry:       map[string]int{"US": 20, "JP": 10, "SG": 2, "XX": 50},
	}
	var warnings bytes.Buffer

	points := preparePageHistory(path, summary, 3, &warnings)
	if len(points) != 1 {
		t.Fatalf("points = %d, want current snapshot only", len(points))
	}
	if points[0].Selected != 2000 || points[0].Verified != 7866 || points[0].Countries != 2 {
		t.Fatalf("current snapshot = %+v", points[0])
	}
	if count := strings.Count(warnings.String(), "warning: page history:"); count != 1 {
		t.Fatalf("warning count = %d, output %q", count, warnings.String())
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != corrupt {
		t.Fatalf("corrupt file was overwritten: %q", body)
	}

	siteDir := filepath.Join(t.TempDir(), "docs")
	if err := pages.Generate(pages.Input{
		Title:         "Free VPN Subscriptions",
		RepoURL:       "https://github.com/example/repo",
		SiteURL:       "https://example.github.io/repo",
		Summary:       summary,
		MinPerCountry: 3,
		History:       points,
	}, siteDir); err != nil {
		t.Fatalf("Generate() after history degradation: %v", err)
	}
	if _, err := os.Stat(filepath.Join(siteDir, "index.html")); err != nil {
		t.Fatalf("generated index: %v", err)
	}
}
