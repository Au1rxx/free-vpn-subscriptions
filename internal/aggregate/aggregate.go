// Package aggregate deduplicates, filters, and ranks probed nodes into a
// final shortlist ready for output.
package aggregate

import (
	"sort"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/node"
)

// Summary captures aggregated statistics for downstream consumers
// (status.json, README).
type Summary struct {
	TotalFetched    int            `json:"total_fetched"`
	TotalAlive      int            `json:"total_alive"`
	TotalVerified   int            `json:"total_verified,omitempty"`
	TotalSelected   int            `json:"total_selected"`
	BySource        map[string]int `json:"by_source"`
	ByProtocol      map[string]int `json:"by_protocol"`
	ByCountry       map[string]int `json:"by_country,omitempty"`
	MedianLatencyMS int            `json:"median_latency_ms"`
	MinLatencyMS    int            `json:"min_latency_ms"`
	GeneratedAtUnix int64          `json:"generated_at_unix"`
}

// Run applies the aggregation pipeline: protocol filter → dedup → RTT filter
// → sort by latency → top-N. Input `alive` is the set of probe-passing nodes.
func Run(alive []*node.Node, cfg config.AggregateConfig) ([]*node.Node, Summary) {
	proto := allowedProtocols(cfg.Protocols)

	// 1. protocol + RTT filter
	filtered := make([]*node.Node, 0, len(alive))
	for _, n := range alive {
		if !proto[n.Protocol] {
			continue
		}
		if cfg.MaxRTTMS > 0 && n.LatencyMS > cfg.MaxRTTMS {
			continue
		}
		filtered = append(filtered, n)
	}

	// 2. dedup by endpoint key (keep fastest)
	best := make(map[string]*node.Node, len(filtered))
	for _, n := range filtered {
		k := n.Key()
		if prev, ok := best[k]; !ok || n.LatencyMS < prev.LatencyMS {
			best[k] = n
		}
	}
	deduped := make([]*node.Node, 0, len(best))
	for _, n := range best {
		deduped = append(deduped, n)
	}

	// 3. sort by latency ascending
	sort.Slice(deduped, func(i, j int) bool {
		return deduped[i].LatencyMS < deduped[j].LatencyMS
	})

	// 4. top-N
	selected := deduped
	if cfg.TopN > 0 && len(selected) > cfg.TopN {
		selected = selected[:cfg.TopN]
	}

	return selected, buildSummary(alive, selected)
}

func allowedProtocols(list []string) map[string]bool {
	if len(list) == 0 {
		return map[string]bool{
			node.ProtoVLESS: true, node.ProtoVMess: true,
			node.ProtoTrojan: true, node.ProtoSS: true,
			node.ProtoHysteria2: true,
		}
	}
	out := make(map[string]bool, len(list))
	for _, p := range list {
		// accept short aliases
		switch p {
		case "ss":
			out[node.ProtoSS] = true
		case "hy2":
			out[node.ProtoHysteria2] = true
		default:
			out[p] = true
		}
	}
	return out
}

func buildSummary(alive, selected []*node.Node) Summary {
	s := Summary{
		TotalAlive:    len(alive),
		TotalSelected: len(selected),
		BySource:      map[string]int{},
		ByProtocol:    map[string]int{},
	}
	for _, n := range selected {
		s.BySource[n.SourceName]++
		s.ByProtocol[n.Protocol]++
	}
	if len(selected) > 0 {
		s.MinLatencyMS = selected[0].LatencyMS
		s.MedianLatencyMS = medianLatency(selected)
	}
	return s
}

// medianLatency returns the true median of a latency-ascending slice:
// middle element for odd length, mean of the two middle elements for even.
func medianLatency(sorted []*node.Node) int {
	n := len(sorted)
	mid := n / 2
	if n%2 == 1 {
		return sorted[mid].LatencyMS
	}
	return (sorted[mid-1].LatencyMS + sorted[mid].LatencyMS) / 2
}
