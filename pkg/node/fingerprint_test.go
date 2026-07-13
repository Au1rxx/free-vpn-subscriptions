package node

import (
	"bytes"
	"testing"
)

func TestConfigFingerprintKeepsCredentialsDistinct(t *testing.T) {
	a := &Node{Protocol: ProtoVLESS, Server: "EXAMPLE.com.", Port: 443, UUID: "a", Name: "one"}
	b := &Node{Protocol: ProtoVLESS, Server: "example.com", Port: 443, UUID: "b", Name: "two"}
	if a.EndpointFingerprint() != b.EndpointFingerprint() {
		t.Fatal("same endpoint must match")
	}
	if a.ConfigFingerprint() == b.ConfigFingerprint() {
		t.Fatal("different credentials must differ")
	}
}

func TestConfigFingerprintIgnoresDisplayAndRuntimeFields(t *testing.T) {
	a := &Node{Protocol: ProtoVLESS, Server: "example.com", Port: 443, UUID: "a", ALPN: "h2,http/1.1"}
	b := *a
	b.Name, b.Country, b.SourceName, b.LatencyMS, b.TCPLatencyMS = "other", "US", "source", 99, 50
	b.ALPN = "http/1.1,h2"
	if a.ConfigFingerprint() != b.ConfigFingerprint() {
		t.Fatal("display/runtime fields or ALPN order changed fingerprint")
	}
}

func TestCanonicalJSONIsStableAndDoesNotContainDisplayFields(t *testing.T) {
	n := &Node{Protocol: ProtoTrojan, Server: "EXAMPLE.com.", Port: 443, Password: "secret", Name: "display", Extra: map[string]string{"z": "2", "a": "1"}}
	first, err := n.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	second, err := n.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatalf("canonical JSON changed: %s != %s", first, second)
	}
	if bytes.Contains(first, []byte("display")) {
		t.Fatalf("canonical JSON contains display field: %s", first)
	}
}
