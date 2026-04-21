package parse

import "testing"

func TestURIList_SkipsEmptyAndComments(t *testing.T) {
	body := `
# a comment
trojan://pw@example.com:443?sni=example.com#ok

not-a-uri
vless://11111111-1111-1111-1111-111111111111@example.com:443?security=tls&sni=example.com#vl
`
	got := URIList(body)
	if len(got) != 2 {
		t.Fatalf("got %d nodes, want 2", len(got))
	}
	if got[0].Protocol != "trojan" {
		t.Errorf("first node protocol = %q, want trojan", got[0].Protocol)
	}
	if got[1].Protocol != "vless" {
		t.Errorf("second node protocol = %q, want vless", got[1].Protocol)
	}
}

func TestBase64List_RoundTrip(t *testing.T) {
	// Base64 of "trojan://pw@example.com:443?sni=example.com#ok"
	body := []byte("dHJvamFuOi8vcHdAZXhhbXBsZS5jb206NDQzP3NuaT1leGFtcGxlLmNvbSNvaw==")
	got, err := Base64List(body)
	if err != nil {
		t.Fatalf("Base64List error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d nodes, want 1", len(got))
	}
	if got[0].Server != "example.com" || got[0].Port != 443 {
		t.Errorf("decoded server:port = %s:%d, want example.com:443", got[0].Server, got[0].Port)
	}
}
