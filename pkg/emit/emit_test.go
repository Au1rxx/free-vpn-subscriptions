package emit

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

func sampleNodes() []*node.Node {
	return []*node.Node{
		{
			Name:     "demo-trojan",
			Protocol: node.ProtoTrojan,
			Server:   "trojan.example.com",
			Port:     443,
			Password: "pw",
			SNI:      "trojan.example.com",
		},
		{
			Name:     "demo-ss",
			Protocol: node.ProtoSS,
			Server:   "ss.example.com",
			Port:     8388,
			Cipher:   "aes-256-gcm",
			Password: "pw2",
		},
	}
}

func TestClash_ContainsProxies(t *testing.T) {
	out, err := Clash(sampleNodes())
	if err != nil {
		t.Fatalf("Clash error: %v", err)
	}
	if !strings.Contains(out, "proxies:") {
		t.Error("Clash output missing proxies block")
	}
	if !strings.Contains(out, "trojan.example.com") {
		t.Error("Clash output missing trojan server")
	}
}

func TestSingbox_ValidJSON(t *testing.T) {
	out, err := Singbox(sampleNodes())
	if err != nil {
		t.Fatalf("Singbox error: %v", err)
	}
	var cfg map[string]any
	if err := json.Unmarshal([]byte(out), &cfg); err != nil {
		t.Fatalf("Singbox output is not valid JSON: %v", err)
	}
	if _, ok := cfg["outbounds"]; !ok {
		t.Error("Singbox output missing outbounds")
	}
}

func TestV2RayBase64_Decodable(t *testing.T) {
	out := V2RayBase64(sampleNodes())
	if out == "" {
		t.Fatal("V2RayBase64 returned empty string")
	}
	// Should be pure base64 — no line breaks, only base64 alphabet.
	for _, r := range out {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') || r == '+' || r == '/' || r == '=') {
			t.Fatalf("V2RayBase64 contains non-base64 char %q", r)
		}
	}
}
