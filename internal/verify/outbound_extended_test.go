package verify

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

func TestBuildOutboundExtendedProtocols(t *testing.T) {
	nodes := []*node.Node{
		{Protocol: node.ProtoVLESS, Server: "example.com", Port: 443, UUID: "id"},
		{Protocol: node.ProtoVMess, Server: "example.com", Port: 443, UUID: "id"},
		{Protocol: node.ProtoTrojan, Server: "example.com", Port: 443, Password: "pw", Security: "tls"},
		{Protocol: node.ProtoSS, Server: "example.com", Port: 443, Cipher: "aes-128-gcm", Password: "pw"},
		{Protocol: node.ProtoHysteria2, Server: "example.com", Port: 443, Password: "pw", Security: "tls", Extra: map[string]string{"obfs": "salamander", "obfs_password": "mask"}},
		{Protocol: node.ProtoTUIC, Server: "example.com", Port: 443, UUID: "id", Password: "pw", Security: "tls", Extra: map[string]string{"congestion_control": "bbr"}},
		{Protocol: node.ProtoWireGuard, Server: "example.com", Port: 51820, Password: "private", PublicKey: "peer", Extra: map[string]string{"address": "10.0.0.2/32", "reserved": "1,2,3"}},
		{Protocol: node.ProtoSOCKS5, Server: "example.com", Port: 1080, Username: "user", Password: "pw"},
		{Protocol: node.ProtoHTTP, Server: "example.com", Port: 8080, Username: "user", Password: "pw"},
	}
	for _, n := range nodes {
		outbound, err := BuildOutbound(n, "proxy")
		if err != nil || outbound["tag"] != "proxy" || outbound["type"] == "" {
			t.Fatalf("protocol=%s outbound=%+v err=%v", n.Protocol, outbound, err)
		}
	}
}

func TestOutboundAcceptedBySingBox(t *testing.T) {
	bin, err := exec.LookPath("sing-box")
	if err != nil {
		t.Skip("sing-box is not installed")
	}
	nodes := []*node.Node{
		{Protocol: node.ProtoVLESS, Server: "example.com", Port: 443, UUID: "11111111-1111-1111-1111-111111111111"},
		{Protocol: node.ProtoTrojan, Server: "example.com", Port: 443, Password: "pw", Security: "tls"},
		{Protocol: node.ProtoSS, Server: "example.com", Port: 443, Cipher: "aes-128-gcm", Password: "pw"},
		{Protocol: node.ProtoHysteria2, Server: "example.com", Port: 443, Password: "pw", Security: "tls"},
		{Protocol: node.ProtoTUIC, Server: "example.com", Port: 443, UUID: "11111111-1111-1111-1111-111111111111", Password: "pw", Security: "tls"},
		{Protocol: node.ProtoWireGuard, Server: "example.com", Port: 51820, Password: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", PublicKey: "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB=", Extra: map[string]string{"address": "10.0.0.2/32"}},
		{Protocol: node.ProtoSOCKS5, Server: "example.com", Port: 1080, Username: "user", Password: "pw"},
		{Protocol: node.ProtoHTTP, Server: "example.com", Port: 8080, Username: "user", Password: "pw"},
	}
	for _, n := range nodes {
		t.Run(n.Protocol, func(t *testing.T) {
			outbound, err := BuildOutbound(n, "proxy")
			if err != nil {
				t.Fatal(err)
			}
			config := map[string]any{"outbounds": []map[string]any{outbound}}
			if IsEndpoint(n) {
				config = map[string]any{"endpoints": []map[string]any{outbound}}
			}
			body, err := json.Marshal(config)
			if err != nil {
				t.Fatal(err)
			}
			path := t.TempDir() + "/config.json"
			if err := os.WriteFile(path, body, 0o600); err != nil {
				t.Fatal(err)
			}
			if output, err := exec.Command(bin, "check", "-c", path).CombinedOutput(); err != nil {
				t.Fatalf("sing-box rejected outbound: %v: %s", err, output)
			}
		})
	}
}

func TestBuildOutboundRejectsUnsupported(t *testing.T) {
	if _, err := BuildOutbound(&node.Node{Protocol: "unknown", Server: "example.com", Port: 1}, "proxy"); err == nil {
		t.Fatal("unsupported protocol was accepted")
	}
}

func TestBuildOutboundOmitsEmptyProxyCredentials(t *testing.T) {
	for _, n := range []*node.Node{
		{Protocol: node.ProtoSOCKS5, Server: "192.0.2.1", Port: 1080},
		{Protocol: node.ProtoHTTP, Server: "192.0.2.2", Port: 8080},
	} {
		outbound, err := BuildOutbound(n, "proxy")
		if err != nil {
			t.Fatal(err)
		}
		if _, present := outbound["username"]; present {
			t.Fatalf("anonymous %s outbound contains username: %+v", n.Protocol, outbound)
		}
		if _, present := outbound["password"]; present {
			t.Fatalf("anonymous %s outbound contains password: %+v", n.Protocol, outbound)
		}
	}
}
