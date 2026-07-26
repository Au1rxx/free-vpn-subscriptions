package emit

import (
	"fmt"
	"strings"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// nameAllocator hands out display names that stay stable from one run to the
// next.
//
// Names deliberately exclude the node's position in the list. Ranking shifts
// every run as latencies are re-measured, so an index-prefixed name changes
// for almost every node even when the published set barely moves: one hourly
// publish rewrote 54% of clash.yaml while only 13 of ~1350 nodes had actually
// changed. That churn inflates the repository and, because clients key their
// per-node latency history and manual selection off the name, it also throws
// away everything the user's client had learned.
type nameAllocator struct{ taken map[string]bool }

func newNameAllocator() *nameAllocator {
	return &nameAllocator{taken: make(map[string]bool)}
}

// name returns a unique display name for n. Uniqueness is required: Clash
// rejects duplicate proxy names and sing-box rejects duplicate outbound tags.
// Collisions are resolved with a numeric suffix, which only affects the
// genuinely duplicated names rather than the whole list.
func (a *nameAllocator) name(n *node.Node) string {
	base := safe(n.Name)
	protocol := safe(string(n.Protocol))
	switch {
	case base == "" || base == "node":
		base = protocol
	case protocol != "" && protocol != "node" && !strings.HasPrefix(base, protocol+"-"):
		base = protocol + "-" + base
	}
	if base == "" {
		base = "node"
	}
	candidate := base
	for suffix := 2; a.taken[candidate]; suffix++ {
		candidate = fmt.Sprintf("%s-%d", base, suffix)
	}
	a.taken[candidate] = true
	return candidate
}
