package classify

import (
	"testing"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

func TestLatencyPercentilesAreProtocolLocal(t *testing.T) {
	values := latencyPercentiles([]store.ClassificationCandidate{{NodeConfigID: 1, Protocol: "vless", LatencyMS: 10}, {NodeConfigID: 2, Protocol: "vless", LatencyMS: 20}, {NodeConfigID: 3, Protocol: "trojan", LatencyMS: 100}})
	if values[1] != 0 || values[2] != 1 || values[3] != 0 {
		t.Fatalf("percentiles=%v", values)
	}
}
