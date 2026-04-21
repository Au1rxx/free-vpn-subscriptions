package emit

import (
	"testing"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

func TestSingboxOutbound_ReturnsTaggedOutbound(t *testing.T) {
	n := &node.Node{
		Name: "t1", Protocol: node.ProtoTrojan,
		Server: "t.example.com", Port: 443,
		Password: "pw", SNI: "t.example.com",
	}
	ob := SingboxOutbound(n, "proxy")
	if ob == nil {
		t.Fatal("got nil outbound")
	}
	if ob["tag"] != "proxy" {
		t.Errorf("tag=%v, want proxy", ob["tag"])
	}
	if ob["type"] != "trojan" {
		t.Errorf("type=%v, want trojan", ob["type"])
	}
	if ob["server"] != "t.example.com" {
		t.Errorf("server=%v", ob["server"])
	}
	if ob["server_port"] != 443 {
		t.Errorf("port=%v", ob["server_port"])
	}
}

func TestSingboxOutbound_UnknownProtocolReturnsNil(t *testing.T) {
	n := &node.Node{Name: "x", Protocol: "made-up", Server: "x", Port: 1}
	if ob := SingboxOutbound(n, "p"); ob != nil {
		t.Errorf("want nil for unknown protocol, got %v", ob)
	}
}
