package node

import "testing"

func TestIsSupportedProtocol(t *testing.T) {
	supported := []string{
		ProtoVLESS, ProtoVMess, ProtoTrojan, ProtoSS, ProtoSSR,
		ProtoHysteria, ProtoHysteria2, ProtoTUIC, ProtoWireGuard,
		ProtoSOCKS4, ProtoSOCKS5, ProtoHTTP, ProtoHTTPS,
	}
	for _, protocol := range supported {
		if !IsSupportedProtocol(protocol) {
			t.Errorf("protocol %q should be supported", protocol)
		}
	}
	for _, protocol := range []string{"", "invalid", "VLESS", "http "} {
		if IsSupportedProtocol(protocol) {
			t.Errorf("protocol %q should not be supported", protocol)
		}
	}
}
