package lifecycle

import "time"

type Input struct {
	State                              string
	Now, FirstSeenAt, LastSeenAt       time.Time
	LastSuccessAt, ArchivedAt          *time.Time
	ActiveSources, ConsecutiveFailures int
	EverSucceeded                      bool
}

type Decision struct {
	State                           string
	NextCheckAt, ArchiveAt, PurgeAt *time.Time
	Reason                          string
}

func Decide(input Input) Decision {
	now := input.Now.UTC()
	if input.ActiveSources > 0 && (input.State == "archived" || input.State == "dead" || input.State == "purged") {
		return Decision{State: "pending", NextCheckAt: pointer(now), Reason: "source_reappeared"}
	}
	if input.State == "archived" {
		archive := input.ArchivedAt
		if archive != nil && input.ActiveSources == 0 && now.Sub(*archive) >= 90*24*time.Hour {
			return Decision{State: "purged", PurgeAt: pointer(now), Reason: "archive_ttl_elapsed"}
		}
		return Decision{State: "archived", PurgeAt: pointer(archiveValue(archive, now).Add(90 * 24 * time.Hour)), Reason: "archive_retained"}
	}
	if input.LastSuccessAt != nil {
		age := now.Sub(*input.LastSuccessAt)
		if age < 7*24*time.Hour {
			state, reason := "active", "recent_success"
			if input.ConsecutiveFailures > 0 {
				state, reason = "degraded", "recent_failures"
			}
			return Decision{State: state, NextCheckAt: pointer(now.Add(time.Hour)), Reason: reason}
		}
		if input.ActiveSources > 0 {
			return Decision{State: "stale", NextCheckAt: pointer(now.Add(6 * time.Hour)), Reason: "success_older_than_7d"}
		}
		if age >= 90*24*time.Hour {
			return Decision{State: "dead", ArchiveAt: pointer(now), Reason: "successful_node_absent_90d"}
		}
		return Decision{State: "stale", NextCheckAt: pointer(now.Add(24 * time.Hour)), Reason: "source_absent"}
	}
	if input.ActiveSources > 0 && now.Sub(input.FirstSeenAt) < 30*24*time.Hour {
		return Decision{State: "pending", NextCheckAt: pointer(now.Add(15 * time.Minute)), Reason: "awaiting_first_success"}
	}
	if input.ActiveSources == 0 && now.Sub(input.FirstSeenAt) >= 30*24*time.Hour {
		return Decision{State: "dead", ArchiveAt: pointer(now), Reason: "never_succeeded_30d"}
	}
	return Decision{State: "pending", NextCheckAt: pointer(now.Add(time.Hour)), Reason: "unverified"}
}

func pointer(value time.Time) *time.Time { return &value }
func archiveValue(value *time.Time, fallback time.Time) time.Time {
	if value == nil {
		return fallback
	}
	return *value
}
