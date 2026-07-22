package maintain

import (
	"errors"
	"testing"
	"time"
)

func TestRunRequiredRollupsFailClosedBeforeMaintenance(t *testing.T) {
	want := errors.New("rollup failed")
	dates := []time.Time{time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC), time.Date(2026, 7, 8, 0, 0, 0, 0, time.UTC)}
	var called []time.Time
	err := runRequiredRollups(false, dates, func(date time.Time) error {
		called = append(called, date)
		if date.Equal(dates[1]) {
			return want
		}
		return nil
	})
	if len(called) != 2 || !errors.Is(err, want) {
		t.Fatalf("called=%v error=%v, want second rollup failure", called, err)
	}
}

func TestRunRequiredRollupsSkipDryRunWrites(t *testing.T) {
	called := false
	if err := runRequiredRollups(true, []time.Time{{}}, func(time.Time) error {
		called = true
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if called {
		t.Fatal("dry-run must not write daily rollups")
	}
}
