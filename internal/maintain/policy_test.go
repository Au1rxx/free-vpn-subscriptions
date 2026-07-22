package maintain

import "testing"

func TestCapacityPoliciesAtWatermarks(t *testing.T) {
	for _, test := range []struct {
		usage                                     float64
		raw, parse, attempt, batch, fetch, export int
		pause                                     bool
	}{
		{69.999, 30, 30, 14, 14, 90, 30, false},
		{70, 14, 14, 7, 7, 60, 14, false},
		{80, 7, 7, 3, 3, 30, 7, false},
		{90, 3, 3, 2, 2, 14, 3, true},
		{94, 1, 1, 1, 1, 7, 1, true},
	} {
		policy := PolicyForUsage(test.usage)
		if policy.RawPayloadDays != test.raw || policy.ParseErrorDays != test.parse ||
			policy.AttemptDays != test.attempt || policy.BatchDays != test.batch ||
			policy.FetchDays != test.fetch || policy.ExportDays != test.export ||
			policy.PauseColdSources != test.pause {
			t.Fatalf("usage=%v policy=%+v", test.usage, policy)
		}
	}
}
