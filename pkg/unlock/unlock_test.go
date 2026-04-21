package unlock

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// redirectClient routes every request to a single httptest server by
// rewriting req.URL.Host / Scheme. This sidesteps real DNS + TLS while
// keeping the per-path routing logic inside each handler.
func redirectClient(srv *httptest.Server) *http.Client {
	return &http.Client{
		Transport: &rewriteRT{target: srv.URL, inner: srv.Client().Transport},
	}
}

type rewriteRT struct {
	target string
	inner  http.RoundTripper
}

func (r *rewriteRT) RoundTrip(req *http.Request) (*http.Response, error) {
	// Preserve the original host in a header so the mock handler can
	// dispatch by "what would have been hit".
	origHost := req.URL.Host
	origPath := req.URL.Path
	req = req.Clone(req.Context())
	req.Header.Set("X-Orig-Host", origHost)
	// Rewrite to the mock server.
	newURL := r.target + origPath
	if req.URL.RawQuery != "" {
		newURL += "?" + req.URL.RawQuery
	}
	parsed, err := req.URL.Parse(newURL)
	if err != nil {
		return nil, err
	}
	req.URL = parsed
	req.Host = ""
	inner := r.inner
	if inner == nil {
		inner = http.DefaultTransport
	}
	return inner.RoundTrip(req)
}

// mockRouter lets each test define a per-(host,path) response.
type mockRouter struct {
	routes map[string]http.HandlerFunc
}

func newMockRouter() *mockRouter { return &mockRouter{routes: map[string]http.HandlerFunc{}} }

func (m *mockRouter) on(host, path string, h http.HandlerFunc) {
	m.routes[host+path] = h
}

func (m *mockRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("X-Orig-Host") + r.URL.Path
	if h, ok := m.routes[key]; ok {
		h(w, r)
		return
	}
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintln(w, "unmatched:", key)
}

// --- Netflix -----------------------------------------------------------

func TestCheckNetflix(t *testing.T) {
	cases := []struct {
		name                       string
		originalCode, licensedCode int
		wantStatus                 Status
	}{
		{"full-unlock", 200, 200, Unlocked},
		{"originals-only", 200, 404, Partial},
		{"blocked", 404, 404, Blocked},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := newMockRouter()
			router.on("www.netflix.com", "/title/81280792", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(c.originalCode) })
			router.on("www.netflix.com", "/title/70143836", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(c.licensedCode) })
			srv := httptest.NewServer(router)
			defer srv.Close()

			r := CheckNetflix(context.Background(), redirectClient(srv))
			if r.Status != c.wantStatus {
				t.Errorf("got %s, want %s (detail=%q)", r.Status, c.wantStatus, r.Detail)
			}
		})
	}
}

// --- Disney+ -----------------------------------------------------------

func TestCheckDisney(t *testing.T) {
	cases := []struct {
		name       string
		handler    http.HandlerFunc
		wantStatus Status
		wantRegion string
	}{
		{
			name: "unlocked-with-lang",
			handler: func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, `<html lang="en-US"><body>Welcome</body></html>`)
			},
			wantStatus: Unlocked, wantRegion: "US",
		},
		{
			name: "blocked-via-body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, `<html><body>Disney+ is not available in your region yet.</body></html>`)
			},
			wantStatus: Blocked,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := newMockRouter()
			router.on("www.disneyplus.com", "/", c.handler)
			srv := httptest.NewServer(router)
			defer srv.Close()

			r := CheckDisney(context.Background(), redirectClient(srv))
			if r.Status != c.wantStatus {
				t.Errorf("status=%s want=%s detail=%q", r.Status, c.wantStatus, r.Detail)
			}
			if r.Region != c.wantRegion {
				t.Errorf("region=%q want=%q", r.Region, c.wantRegion)
			}
		})
	}
}

// --- YouTube Premium --------------------------------------------------

func TestCheckYouTubePremium(t *testing.T) {
	cases := []struct {
		name       string
		body       string
		wantStatus Status
		wantRegion string
	}{
		{"unlocked-US", `...,"INNERTUBE_CONTEXT":{"gl":"US","hl":"en"},...`, Unlocked, "US"},
		{"unlocked-JP", `..."countryCode":"JP"...`, Unlocked, "JP"},
		{"blocked-body", "Premium is not available in your country.", Blocked, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := newMockRouter()
			router.on("www.youtube.com", "/premium", func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, c.body)
			})
			srv := httptest.NewServer(router)
			defer srv.Close()

			r := CheckYouTubePremium(context.Background(), redirectClient(srv))
			if r.Status != c.wantStatus {
				t.Errorf("status=%s want=%s detail=%q body=%q", r.Status, c.wantStatus, r.Detail, c.body)
			}
			if r.Region != c.wantRegion {
				t.Errorf("region=%q want=%q", r.Region, c.wantRegion)
			}
		})
	}
}

// --- ChatGPT ----------------------------------------------------------

func TestCheckChatGPT(t *testing.T) {
	cases := []struct {
		name         string
		complianceOK bool
		complBody    string
		traceLoc     string
		wantStatus   Status
		wantRegion   string
	}{
		{"unlocked-us", true, `{"cookie_requirements":{"ok":true}}`, "US", Unlocked, "US"},
		{"blocked-hk", false, `{"error":"unsupported_country"}`, "HK", Blocked, "HK"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			router := newMockRouter()
			router.on("chat.openai.com", "/cdn-cgi/trace", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "fl=abc\nh=chat.openai.com\nip=1.2.3.4\nloc=%s\n", strings.ToLower(c.traceLoc))
			})
			router.on("api.openai.com", "/compliance/cookie_requirements", func(w http.ResponseWriter, r *http.Request) {
				if c.complianceOK {
					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
				io.WriteString(w, c.complBody)
			})
			srv := httptest.NewServer(router)
			defer srv.Close()

			r := CheckChatGPT(context.Background(), redirectClient(srv))
			if r.Status != c.wantStatus {
				t.Errorf("status=%s want=%s detail=%q", r.Status, c.wantStatus, r.Detail)
			}
			if r.Region != c.wantRegion {
				t.Errorf("region=%q want=%q", r.Region, c.wantRegion)
			}
		})
	}
}

// --- Run --------------------------------------------------------------

func TestRun_AllTargetsGetCalled(t *testing.T) {
	var calls []string
	mkCheck := func(name string, status Status) CheckFunc {
		return func(ctx context.Context, client *http.Client) Result {
			calls = append(calls, name)
			return Result{Target: name, Status: status}
		}
	}
	targets := []Target{
		{Name: "a", Check: mkCheck("a", Unlocked)},
		{Name: "b", Check: mkCheck("b", Blocked)},
		{Name: "c", Check: mkCheck("c", Partial)},
	}
	results := Run(context.Background(), http.DefaultClient, targets, 100*time.Millisecond)
	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}
	wantCalls := []string{"a", "b", "c"}
	if strings.Join(calls, ",") != strings.Join(wantCalls, ",") {
		t.Errorf("calls = %v, want %v", calls, wantCalls)
	}
	for i, r := range results {
		if r.Target != wantCalls[i] {
			t.Errorf("result[%d].Target = %q, want %q", i, r.Target, wantCalls[i])
		}
		if r.StatusText != r.Status.String() {
			t.Errorf("StatusText not populated: %+v", r)
		}
	}
}

func TestStatusString(t *testing.T) {
	cases := map[Status]string{
		Unknown: "unknown", Blocked: "blocked", Partial: "partial", Unlocked: "unlocked",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("Status(%d).String() = %q, want %q", s, got, want)
		}
	}
}
