package validation

import "time"

func NextAttempt(state string, failures int, now time.Time) time.Time {
	if state == "success" {
		return now.Add(6 * time.Hour)
	}
	if state == "partial" {
		return now.Add(time.Hour)
	}
	intervals := []time.Duration{15 * time.Minute, time.Hour, 6 * time.Hour, 24 * time.Hour, 3 * 24 * time.Hour, 7 * 24 * time.Hour}
	if failures < 1 {
		failures = 1
	}
	if failures > len(intervals) {
		failures = len(intervals)
	}
	return now.Add(intervals[failures-1])
}
