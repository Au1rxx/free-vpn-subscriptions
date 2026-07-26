package emit

import (
	"strings"
	"testing"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

func TestNameAllocatorIsIndependentOfPosition(t *testing.T) {
	first := &node.Node{Protocol: node.ProtoVLESS, Name: "vless-100"}
	second := &node.Node{Protocol: node.ProtoTrojan, Name: "trojan-200"}

	forward := newNameAllocator()
	a1, b1 := forward.name(first), forward.name(second)

	// Same nodes, opposite order — a re-rank must not rename anything.
	reverse := newNameAllocator()
	b2, a2 := reverse.name(second), reverse.name(first)

	if a1 != a2 || b1 != b2 {
		t.Fatalf("names changed with position: %q/%q vs %q/%q", a1, b1, a2, b2)
	}
}

func TestNameAllocatorDropsIndexPrefix(t *testing.T) {
	alloc := newNameAllocator()
	got := alloc.name(&node.Node{Protocol: node.ProtoVLESS, Name: "vless-153601537"})
	if got != "vless-153601537" {
		t.Fatalf("name=%q", got)
	}
}

// A name that does not already carry its protocol keeps the protocol prefix,
// which is how nodes parsed from upstream subscriptions stay identifiable.
func TestNameAllocatorPrefixesProtocolWhenAbsent(t *testing.T) {
	alloc := newNameAllocator()
	got := alloc.name(&node.Node{Protocol: node.ProtoTrojan, Name: "Tokyo 01"})
	if got != "trojan-Tokyo 01" {
		t.Fatalf("name=%q", got)
	}
}

func TestNameAllocatorResolvesCollisions(t *testing.T) {
	alloc := newNameAllocator()
	var got []string
	for i := 0; i < 3; i++ {
		got = append(got, alloc.name(&node.Node{Protocol: node.ProtoVMess, Name: "vmess-dup"}))
	}
	want := []string{"vmess-dup", "vmess-dup-2", "vmess-dup-3"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("collision names=%v want %v", got, want)
		}
	}
}

func TestNameAllocatorHandlesEmptyName(t *testing.T) {
	alloc := newNameAllocator()
	if got := alloc.name(&node.Node{Protocol: node.ProtoVLESS}); got != "vless" {
		t.Fatalf("name=%q", got)
	}
}

// Clash rejects duplicate proxy names and sing-box rejects duplicate tags, so
// the emitters must never produce a repeat even when the inputs collide.
func TestEmittersProduceUniqueStableNames(t *testing.T) {
	nodes := []*node.Node{
		{Protocol: node.ProtoVLESS, Server: "a.example", Port: 443, UUID: "u1", Name: "dup"},
		{Protocol: node.ProtoVLESS, Server: "b.example", Port: 443, UUID: "u2", Name: "dup"},
		{Protocol: node.ProtoTrojan, Server: "c.example", Port: 443, Password: "p", Name: "trojan-9"},
	}
	profile, err := Clash(nodes)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"vless-dup", "vless-dup-2", "trojan-9"} {
		if !strings.Contains(profile, want) {
			t.Errorf("clash profile missing %q", want)
		}
	}
	// The old scheme prefixed an ordinal; make sure it is gone.
	if strings.Contains(profile, "01-vless") || strings.Contains(profile, "name: 01-") {
		t.Error("clash profile still carries an index-prefixed name")
	}
}
