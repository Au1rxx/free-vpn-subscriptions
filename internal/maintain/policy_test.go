package maintain

import "testing"

func TestCapacityPoliciesAtWatermarks(t *testing.T) {
	for _, test := range []struct {
		usage        float64
		raw          int
		pause, store bool
	}{
		{69, 30, false, true}, {70, 14, false, true}, {80, 7, false, true}, {90, 3, true, false}, {94, 1, true, false},
	} {
		policy := PolicyForUsage(test.usage)
		if policy.RawPayloadDays != test.raw || policy.PauseColdSources != test.pause || policy.StoreRawBodies != test.store {
			t.Fatalf("usage=%v policy=%+v", test.usage, policy)
		}
	}
}
