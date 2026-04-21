package emit

import (
	"strings"
	"testing"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

var partialNodes = []*node.Node{
	{Name: "ss-a", Protocol: node.ProtoSS, Server: "1.1.1.1", Port: 8388, Cipher: "aes-256-gcm", Password: "pw1"},
	{Name: "trj-b", Protocol: node.ProtoTrojan, Server: "2.2.2.2", Port: 443, Password: "pw2", SNI: "b.example"},
	{Name: "vm-c", Protocol: node.ProtoVMess, Server: "3.3.3.3", Port: 8443, UUID: "uuid-c", Cipher: "auto", Network: "ws", Path: "/c", Host: "c.example", Security: "tls", SNI: "c.example"},
	// VLESS + Hysteria2 should be silently dropped by the partial emitters.
	{Name: "vl-d", Protocol: node.ProtoVLESS, Server: "4.4.4.4", Port: 443, UUID: "uuid-d", Security: "tls", SNI: "d.example"},
	{Name: "hy-e", Protocol: node.ProtoHysteria2, Server: "5.5.5.5", Port: 443, Password: "pw5", SNI: "e.example"},
}

func TestSurge_StructureAndCoverage(t *testing.T) {
	out, err := Surge(partialNodes)
	if err != nil {
		t.Fatalf("Surge: %v", err)
	}
	for _, want := range []string{
		"[Proxy]", "[Proxy Group]",
		"= ss, 1.1.1.1, 8388", "encrypt-method=aes-256-gcm",
		"= trojan, 2.2.2.2, 443", "sni=b.example",
		"= vmess, 3.3.3.3, 8443", "username=uuid-c", "ws=true", "tls=true",
		"auto = url-test", "select = select, auto",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Surge output missing %q:\n%s", want, out)
		}
	}
	// Partial coverage: VLESS + Hysteria2 should be absent.
	if strings.Contains(out, "4.4.4.4") || strings.Contains(out, "5.5.5.5") {
		t.Errorf("Surge should drop VLESS/Hy2 nodes, got:\n%s", out)
	}
}

func TestQuantumultX_StructureAndCoverage(t *testing.T) {
	out, err := QuantumultX(partialNodes)
	if err != nil {
		t.Fatalf("QuantumultX: %v", err)
	}
	for _, want := range []string{
		"[server_local]",
		"shadowsocks=1.1.1.1:8388", "method=aes-256-gcm",
		"trojan=2.2.2.2:443", "over-tls=true", "tls-host=b.example",
		"vmess=3.3.3.3:8443", "password=uuid-c", "obfs=wss", "obfs-host=c.example",
		"tag=",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("QuantumultX output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "4.4.4.4") || strings.Contains(out, "5.5.5.5") {
		t.Errorf("QuantumultX should drop VLESS/Hy2, got:\n%s", out)
	}
}

func TestLoon_StructureAndCoverage(t *testing.T) {
	out, err := Loon(partialNodes)
	if err != nil {
		t.Fatalf("Loon: %v", err)
	}
	for _, want := range []string{
		"[Proxy]", "[Proxy Group]",
		"= Shadowsocks,1.1.1.1,8388,aes-256-gcm",
		"= trojan,2.2.2.2,443,", "tls-name=b.example",
		"= vmess,3.3.3.3,8443,auto,", "transport:ws", "path=/c", "over-tls=true",
		"auto = url-test", "select = select,auto,",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Loon output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "4.4.4.4") || strings.Contains(out, "5.5.5.5") {
		t.Errorf("Loon should drop VLESS/Hy2, got:\n%s", out)
	}
}
