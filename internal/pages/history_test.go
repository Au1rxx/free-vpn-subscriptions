package pages

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestUpdateHistorySortsDeduplicatesAndKeepsLatest720Points(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "network-history.json")
	start := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	points := make([]HistoryPoint, 0, 721)
	for i := 0; i < 721; i++ {
		points = append(points, HistoryPoint{
			GeneratedAt:     start.Add(time.Duration(i) * time.Hour),
			Selected:        i,
			Verified:        i + 100,
			MedianLatencyMS: 200,
			Countries:       10,
		})
	}
	writeHistoryDocument(t, path, 1, append(points, points[100]))

	current := points[len(points)-1]
	current.Selected = 2000
	got, err := UpdateHistory(path, current)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 720 {
		t.Fatalf("points = %d, want 720", len(got))
	}
	if !got[0].GeneratedAt.Equal(start.Add(time.Hour)) {
		t.Fatalf("oldest point = %s, want %s", got[0].GeneratedAt, start.Add(time.Hour))
	}
	if got[len(got)-1].Selected != 2000 {
		t.Fatalf("current point did not replace duplicate timestamp: %+v", got[len(got)-1])
	}

	var stored struct {
		SchemaVersion int            `json:"schema_version"`
		Points        []HistoryPoint `json:"points"`
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(body, &stored); err != nil {
		t.Fatal(err)
	}
	if stored.SchemaVersion != 1 || len(stored.Points) != 720 {
		t.Fatalf("stored history = version %d, points %d", stored.SchemaVersion, len(stored.Points))
	}
}

func TestUpdateHistoryRejectsCorruptUnknownAndFutureHistoryWithoutOverwriting(t *testing.T) {
	now := time.Date(2026, 7, 26, 4, 0, 0, 0, time.UTC)
	current := HistoryPoint{GeneratedAt: now, Selected: 100, Verified: 200, MedianLatencyMS: 300, Countries: 4}
	tests := []struct {
		name string
		body string
	}{
		{name: "corrupt", body: `{"schema_version":`},
		{name: "unknown schema", body: `{"schema_version":2,"points":[]}`},
		{name: "future point", body: `{"schema_version":1,"points":[{"generated_at":"2026-07-26T05:00:00Z","selected":1,"verified":1,"median_latency_ms":1,"countries":1}]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "network-history.json")
			if err := os.WriteFile(path, []byte(tt.body), 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := UpdateHistory(path, current); err == nil {
				t.Fatal("UpdateHistory() error = nil, want non-nil")
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.body {
				t.Fatalf("history was overwritten on error: %q", got)
			}
		})
	}
}

func TestBuildTrendsUsesLatestPointAtOrBeforeEachWindow(t *testing.T) {
	now := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)
	points := []HistoryPoint{
		{GeneratedAt: now.Add(-31 * 24 * time.Hour), Selected: 100, Verified: 500, MedianLatencyMS: 500},
		{GeneratedAt: now.Add(-30 * 24 * time.Hour), Selected: 200, Verified: 600, MedianLatencyMS: 450},
		{GeneratedAt: now.Add(-7 * 24 * time.Hour), Selected: 300, Verified: 700, MedianLatencyMS: 400},
		{GeneratedAt: now.Add(-25 * time.Hour), Selected: 350, Verified: 750, MedianLatencyMS: 390},
		{GeneratedAt: now.Add(-24 * time.Hour), Selected: 400, Verified: 800, MedianLatencyMS: 350},
		{GeneratedAt: now, Selected: 450, Verified: 900, MedianLatencyMS: 300},
	}

	got := BuildTrends(points, now)
	if !got.Hours24.Available || got.Hours24.Selected != 50 || got.Hours24.Verified != 100 || got.Hours24.MedianLatencyMS != -50 {
		t.Fatalf("24h trend = %+v", got.Hours24)
	}
	if !got.Days7.Available || got.Days7.Selected != 150 {
		t.Fatalf("7d trend = %+v", got.Days7)
	}
	if !got.Days30.Available || got.Days30.Selected != 250 {
		t.Fatalf("30d trend = %+v", got.Days30)
	}
}

func TestBuildTrendsMarksUnavailableWindows(t *testing.T) {
	now := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)
	got := BuildTrends([]HistoryPoint{
		{GeneratedAt: now.Add(-2 * time.Hour), Selected: 100},
		{GeneratedAt: now, Selected: 120},
	}, now)
	if got.Hours24.Available || got.Days7.Available || got.Days30.Available {
		t.Fatalf("trends should be unavailable: %+v", got)
	}
}

func writeHistoryDocument(t *testing.T, path string, version int, points []HistoryPoint) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(struct {
		SchemaVersion int            `json:"schema_version"`
		Points        []HistoryPoint `json:"points"`
	}{SchemaVersion: version, Points: points})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}
}
