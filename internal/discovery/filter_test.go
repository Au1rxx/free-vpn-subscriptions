package discovery

import "testing"

func TestLikelySubscriptionURL(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"https://raw.githubusercontent.com/owner/repo/main/sub.txt", true},
		{"https://raw.githubusercontent.com/wiki/gfpcom/free-proxy-list/lists/vless.txt", true},
		{"https://example.com/api/subscription?id=1", true},
		{"https://proxypool.link/clash/proxies", true},
		{"https://github.com/owner/repo", false},
		{"https://github.com/owner/repo/issues", false},
		{"https://github.githubassets.com/assets/site.css", false},
		{"https://example.com/favicon.ico", false},
	}
	for _, test := range tests {
		if got := LikelySubscriptionURL(test.url); got != test.want {
			t.Errorf("LikelySubscriptionURL(%q)=%t want %t", test.url, got, test.want)
		}
	}
}
