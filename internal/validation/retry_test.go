package validation

import (
	"testing"
	"time"
)

func TestNextAttemptBackoff(t *testing.T) {
	now := time.Date(2026, 7, 13, 0, 0, 0, 0, time.UTC)
	wants := []time.Duration{15 * time.Minute, time.Hour, 6 * time.Hour, 24 * time.Hour, 3 * 24 * time.Hour, 7 * 24 * time.Hour}
	for index, want := range wants {
		got := NextAttempt("failed", index+1, now).Sub(now)
		if got != want {
			t.Fatalf("failures=%d got=%s want=%s", index+1, got, want)
		}
	}
	if got := NextAttempt("success", 99, now).Sub(now); got != 6*time.Hour {
		t.Fatalf("success retry=%s", got)
	}
}
