package node

import (
	"encoding/base64"
	"testing"
)

func TestParseURIExtendedProtocols(t *testing.T) {
	ssrPayload := "example.com:443:auth_sha1_v4:aes-256-cfb:tls1.2_ticket_auth:" + base64.RawURLEncoding.EncodeToString([]byte("pw")) + "/"
	tests := []struct {
		name, uri, protocol string
	}{
		{"ssr", "ssr://" + base64.RawURLEncoding.EncodeToString([]byte(ssrPayload)), ProtoSSR},
		{"hysteria", "hysteria://token@example.com:443?sni=cdn.example.com", ProtoHysteria},
		{"tuic", "tuic://11111111-1111-1111-1111-111111111111:pw@example.com:443?sni=cdn.example.com", ProtoTUIC},
		{"wireguard", "wireguard://private@example.com:51820?publickey=peer&address=10.0.0.2%2F32&reserved=1%2C2%2C3", ProtoWireGuard},
		{"socks4", "socks4://user:pw@example.com:1080", ProtoSOCKS4},
		{"socks5", "socks5://user:pw@example.com:1080", ProtoSOCKS5},
		{"http", "http://user:pw@example.com:8080", ProtoHTTP},
		{"https", "https://user:pw@example.com:8443", ProtoHTTPS},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := ParseURI(tt.uri)
			if err != nil {
				t.Fatal(err)
			}
			if n.Protocol != tt.protocol || n.Server != "example.com" || n.Port <= 0 || !n.Valid() {
				t.Fatalf("unexpected node: %+v", n)
			}
		})
	}
}

func TestParseURIExtendedRejectsMissingCredentials(t *testing.T) {
	for _, uri := range []string{
		"hysteria://example.com:443", "tuic://example.com:443",
		"wireguard://example.com:51820?publickey=peer",
		"socks5://user@example.com:1080", "http://user@example.com:8080",
	} {
		n, err := ParseURI(uri)
		if err == nil && n.Valid() {
			t.Fatalf("expected missing credential rejection for %s: %+v", uri, n)
		}
	}
}

func TestParseURIExtendedAcceptsAnonymousUserProxies(t *testing.T) {
	for _, uri := range []string{
		"socks4://192.0.2.1:1080", "socks5://192.0.2.2:1080",
		"http://192.0.2.3:8080", "https://192.0.2.4:8443",
	} {
		n, err := ParseURI(uri)
		if err != nil || !n.Valid() {
			t.Fatalf("expected anonymous proxy to be valid for %s: node=%+v err=%v", uri, n, err)
		}
	}
}

func TestParseUserProxyNormalizesZeroPaddedIPv4(t *testing.T) {
	n, err := ParseURI("http://001.224.003.122:3888")
	if err != nil {
		t.Fatal(err)
	}
	if n.Server != "1.224.3.122" {
		t.Fatalf("server=%q, want canonical IPv4", n.Server)
	}
}
