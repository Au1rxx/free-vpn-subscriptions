package pages

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	historySchemaVersion = 1
	maxHistoryPoints     = 30 * 24
)

// HistoryPoint is one hourly public-network snapshot.
type HistoryPoint struct {
	GeneratedAt     time.Time `json:"generated_at"`
	Selected        int       `json:"selected"`
	Verified        int       `json:"verified"`
	MedianLatencyMS int       `json:"median_latency_ms"`
	Countries       int       `json:"countries"`
}

// TrendDelta compares the current snapshot with a historical baseline.
type TrendDelta struct {
	Available       bool
	Selected        int
	Verified        int
	MedianLatencyMS int
}

// TrendSummary contains the rolling windows displayed on the site.
type TrendSummary struct {
	Hours24 TrendDelta
	Days7   TrendDelta
	Days30  TrendDelta
}

type historyDocument struct {
	SchemaVersion int            `json:"schema_version"`
	Points        []HistoryPoint `json:"points"`
}

// UpdateHistory adds current to the rolling history and atomically persists
// the newest 30 days of hourly snapshots.
func UpdateHistory(path string, current HistoryPoint) ([]HistoryPoint, error) {
	current.GeneratedAt = current.GeneratedAt.UTC()
	if err := validateHistoryPoint(current); err != nil {
		return nil, fmt.Errorf("current point: %w", err)
	}

	points, err := readHistory(path, current.GeneratedAt)
	if err != nil {
		return nil, err
	}
	points = append(points, current)

	byTime := make(map[int64]HistoryPoint, len(points))
	for _, point := range points {
		point.GeneratedAt = point.GeneratedAt.UTC()
		byTime[point.GeneratedAt.UnixNano()] = point
	}
	points = points[:0]
	for _, point := range byTime {
		points = append(points, point)
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].GeneratedAt.Before(points[j].GeneratedAt)
	})
	if len(points) > maxHistoryPoints {
		points = points[len(points)-maxHistoryPoints:]
	}

	document := historyDocument{SchemaVersion: historySchemaVersion, Points: points}
	body, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	body = append(body, '\n')
	if err := writeHistoryAtomically(path, body); err != nil {
		return nil, err
	}
	return points, nil
}

func readHistory(path string, now time.Time) ([]HistoryPoint, error) {
	body, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	var document historyDocument
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			err = errors.New("multiple JSON values")
		}
		return nil, fmt.Errorf("decode: %w", err)
	}
	if document.SchemaVersion != historySchemaVersion {
		return nil, fmt.Errorf("schema version %d is unsupported", document.SchemaVersion)
	}
	for i := range document.Points {
		document.Points[i].GeneratedAt = document.Points[i].GeneratedAt.UTC()
		if err := validateHistoryPoint(document.Points[i]); err != nil {
			return nil, fmt.Errorf("point %d: %w", i, err)
		}
		if document.Points[i].GeneratedAt.After(now) {
			return nil, fmt.Errorf("point %d is in the future", i)
		}
	}
	return document.Points, nil
}

func validateHistoryPoint(point HistoryPoint) error {
	if point.GeneratedAt.IsZero() {
		return errors.New("generated_at is required")
	}
	if point.Selected < 0 || point.Verified < 0 || point.MedianLatencyMS < 0 || point.Countries < 0 {
		return errors.New("metrics must not be negative")
	}
	return nil
}

func writeHistoryAtomically(path string, body []byte) (returnErr error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}
	file, err := os.CreateTemp(dir, ".network-history-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}
	tempPath := file.Name()
	defer func() {
		if err := os.Remove(tempPath); err != nil && !errors.Is(err, os.ErrNotExist) && returnErr == nil {
			returnErr = fmt.Errorf("remove temporary file: %w", err)
		}
	}()
	if err := file.Chmod(0o644); err != nil {
		_ = file.Close()
		return fmt.Errorf("set permissions: %w", err)
	}
	if _, err := file.Write(body); err != nil {
		_ = file.Close()
		return fmt.Errorf("write temporary file: %w", err)
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return fmt.Errorf("sync temporary file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace history: %w", err)
	}
	return nil
}

// BuildTrends compares the latest point at or before now with the latest
// baseline at or before each rolling-window boundary.
func BuildTrends(points []HistoryPoint, now time.Time) TrendSummary {
	now = now.UTC()
	sorted := append([]HistoryPoint(nil), points...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].GeneratedAt.Before(sorted[j].GeneratedAt)
	})

	current, ok := latestAtOrBefore(sorted, now)
	if !ok {
		return TrendSummary{}
	}
	return TrendSummary{
		Hours24: buildTrendDelta(sorted, current, now.Add(-24*time.Hour)),
		Days7:   buildTrendDelta(sorted, current, now.Add(-7*24*time.Hour)),
		Days30:  buildTrendDelta(sorted, current, now.Add(-30*24*time.Hour)),
	}
}

func buildTrendDelta(points []HistoryPoint, current HistoryPoint, boundary time.Time) TrendDelta {
	baseline, ok := latestAtOrBefore(points, boundary)
	if !ok {
		return TrendDelta{}
	}
	return TrendDelta{
		Available:       true,
		Selected:        current.Selected - baseline.Selected,
		Verified:        current.Verified - baseline.Verified,
		MedianLatencyMS: current.MedianLatencyMS - baseline.MedianLatencyMS,
	}
}

func latestAtOrBefore(points []HistoryPoint, boundary time.Time) (HistoryPoint, bool) {
	index := sort.Search(len(points), func(i int) bool {
		return points[i].GeneratedAt.After(boundary)
	})
	if index == 0 {
		return HistoryPoint{}, false
	}
	return points[index-1], true
}
