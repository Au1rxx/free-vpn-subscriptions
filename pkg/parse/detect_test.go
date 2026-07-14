package parse

import (
	"encoding/base64"
	"testing"
)

func TestParseAutoMixedKeepsGoodEntriesErrorsAndURLs(t *testing.T) {
	body := []byte("vless://11111111-1111-1111-1111-111111111111@example.com:443?security=tls\nnot-a-node\nhttps://example.test/sub\n")
	result := Parse(body, FormatAuto)
	if result.Format != FormatURIList || len(result.Nodes) != 1 || len(result.Errors) != 1 || len(result.DiscoveredURLs) != 1 {
		t.Fatalf("unexpected result: format=%s nodes=%d errors=%d urls=%d", result.Format, len(result.Nodes), len(result.Errors), len(result.DiscoveredURLs))
	}
	if result.Errors[0].Code != "unsupported_scheme" || result.Errors[0].SampleHash == "" {
		t.Fatalf("unexpected error: %+v", result.Errors[0])
	}
}

func TestParseURIListKeepsHTTPProxiesAsNodes(t *testing.T) {
	body := []byte("http://001.224.003.122:3888\nhttps://user:pw@192.0.2.2:8443\nhttps://example.test:8443/subscription\n")
	result := Parse(body, FormatURIList)
	if len(result.Nodes) != 2 || len(result.DiscoveredURLs) != 1 || len(result.Errors) != 0 {
		t.Fatalf("unexpected result: nodes=%d urls=%d errors=%d", len(result.Nodes), len(result.DiscoveredURLs), len(result.Errors))
	}
	if result.Nodes[0].Protocol != "http" || result.Nodes[0].Server != "1.224.3.122" {
		t.Fatalf("unexpected HTTP proxy: %+v", result.Nodes[0])
	}
}

func TestParseAutoDetectsBase64AndClash(t *testing.T) {
	uri := "trojan://pw@example.com:443?sni=example.com"
	encoded := base64.StdEncoding.EncodeToString([]byte(uri))
	if result := Parse([]byte(encoded), FormatAuto); result.Format != FormatBase64 || len(result.Nodes) != 1 {
		t.Fatalf("base64 result: %+v", result)
	}
	clash := []byte("proxies:\n  - name: one\n    type: trojan\n    server: example.com\n    port: 443\n    password: pw\n")
	if result := Parse(clash, FormatAuto); result.Format != FormatClash || len(result.Nodes) != 1 {
		t.Fatalf("clash result: %+v", result)
	}
}

func TestParseErrorSampleIsBounded(t *testing.T) {
	result := Parse([]byte("bad://"+string(make([]byte, 10000))), FormatURIList)
	if len(result.Errors) != 1 || len(result.Errors[0].SampleHash) != 16 || len(result.Errors[0].Message) > 256 {
		t.Fatalf("unbounded error: %+v", result.Errors)
	}
}
