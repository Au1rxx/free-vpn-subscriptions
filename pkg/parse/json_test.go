package parse

import "testing"

func TestParseSingBoxAndXrayOutbounds(t *testing.T) {
	singbox := []byte(`{"outbounds":[{"type":"vless","tag":"one","server":"example.com","server_port":443,"uuid":"11111111-1111-1111-1111-111111111111","tls":{"enabled":true,"server_name":"cdn.example.com"}}]}`)
	result := Parse(singbox, FormatSingBox)
	if len(result.Nodes) != 1 || result.Nodes[0].Protocol != "vless" || result.Nodes[0].Security != "tls" {
		t.Fatalf("sing-box result: %+v", result)
	}

	xray := []byte(`{"outbounds":[{"protocol":"trojan","settings":{"servers":[{"address":"example.com","port":443,"password":"pw"}]}}]}`)
	result = Parse(xray, FormatXray)
	if len(result.Nodes) != 1 || result.Nodes[0].Protocol != "trojan" {
		t.Fatalf("xray result: %+v", result)
	}
}
