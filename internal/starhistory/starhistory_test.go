package starhistory

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchFollowsPaginationAndSendsGitHubHeaders(t *testing.T) {
	var calls int
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if got := r.Header.Get("Accept"); got != "application/vnd.github.star+json" {
			t.Errorf("Accept=%q", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization=%q", got)
		}
		if r.URL.Path != "/repos/acme/tool/stargazers" {
			t.Errorf("path=%q", r.URL.Path)
		}
		switch r.URL.Query().Get("page") {
		case "":
			w.Header().Set("Link", fmt.Sprintf(`<%s/repos/acme/tool/stargazers?per_page=100&page=2>; rel="next"`, server.URL))
			fmt.Fprint(w, `[{"starred_at":"2026-07-01T10:00:00Z"}]`)
		case "2":
			fmt.Fprint(w, `[{"starred_at":"2026-07-03T12:30:00Z"}]`)
		default:
			http.Error(w, "unexpected page", http.StatusBadRequest)
		}
	}))
	defer server.Close()

	got, err := Fetch(context.Background(), server.Client(), server.URL, "acme/tool", "test-token")
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("calls=%d", calls)
	}
	want := []time.Time{
		time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC),
		time.Date(2026, 7, 3, 12, 30, 0, 0, time.UTC),
	}
	if len(got) != len(want) {
		t.Fatalf("len=%d want=%d", len(got), len(want))
	}
	for i := range want {
		if !got[i].Equal(want[i]) {
			t.Errorf("got[%d]=%s want=%s", i, got[i], want[i])
		}
	}
}

func TestRenderSVGIsDeterministicAndAccumulatesByUTCDate(t *testing.T) {
	stars := []time.Time{
		time.Date(2026, 7, 3, 23, 0, 0, 0, time.FixedZone("west", -7*60*60)),
		time.Date(2026, 7, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2026, 7, 3, 2, 0, 0, 0, time.UTC),
	}
	first, err := RenderSVG("acme/tool", stars)
	if err != nil {
		t.Fatal(err)
	}
	second, err := RenderSVG("acme/tool", stars)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("SVG output is not deterministic")
	}
	if err := xml.Unmarshal(first, new(any)); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}
	text := string(first)
	for _, want := range []string{"<svg", "acme/tool", "3 stars", "2026-07-01", "2026-07-04"} {
		if !strings.Contains(text, want) {
			t.Errorf("SVG missing %q", want)
		}
	}
}

func TestFetchRejectsInvalidRepositoryBeforeRequest(t *testing.T) {
	_, err := Fetch(context.Background(), http.DefaultClient, "https://api.github.com", "not-a-slug", "token")
	if err == nil || !strings.Contains(err.Error(), "owner/repo") {
		t.Fatalf("err=%v", err)
	}
}

func TestFetchRejectsPaginationToAnotherOrigin(t *testing.T) {
	var unexpectedCalls int
	other := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		unexpectedCalls++
	}))
	defer other.Close()
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Link", fmt.Sprintf(`<%s/stolen>; rel="next"`, other.URL))
		fmt.Fprint(w, `[{"starred_at":"2026-07-01T10:00:00Z"}]`)
	}))
	defer api.Close()

	_, err := Fetch(context.Background(), api.Client(), api.URL, "acme/tool", "test-token")
	if err == nil || !strings.Contains(err.Error(), "pagination origin") {
		t.Fatalf("err=%v", err)
	}
	if unexpectedCalls != 0 {
		t.Fatalf("token-bearing request reached another origin")
	}
}

func TestFetchReportsGitHubHTTPError(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "rate limited", http.StatusForbidden)
	}))
	defer api.Close()

	_, err := Fetch(context.Background(), api.Client(), api.URL, "acme/tool", "test-token")
	if err == nil || !strings.Contains(err.Error(), "HTTP 403: rate limited") {
		t.Fatalf("err=%v", err)
	}
}
