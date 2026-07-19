package store

import (
	"encoding/json"
	"testing"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

func TestExportNormalizedConfigRoundTrip(t *testing.T) {
	original := &node.Node{Protocol: node.ProtoVLESS, Server: "example.com", Port: 443, UUID: "id", Security: "tls", Extra: map[string]string{"x": "y"}}
	body, err := original.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	var decoded node.Node
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.ConfigFingerprint() != original.ConfigFingerprint() {
		t.Fatalf("round trip differs: %+v", decoded)
	}
}
