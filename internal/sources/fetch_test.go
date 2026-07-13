package sources

import (
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchRawConditionalAndHash(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("If-None-Match") == `"same"` {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("ETag", `"same"`)
		_, _ = w.Write([]byte("payload"))
	}))
	defer server.Close()
	response, err := FetchRaw(context.Background(), Request{URL: server.URL, Timeout: time.Second})
	if err != nil || string(response.Body) != "payload" || response.SHA256 == ([32]byte{}) {
		t.Fatalf("response=%+v err=%v", response, err)
	}
	response, err = FetchRaw(context.Background(), Request{URL: server.URL, ETag: `"same"`, Timeout: time.Second})
	if err != nil || response.StatusCode != http.StatusNotModified || len(response.Body) != 0 {
		t.Fatalf("304 response=%+v err=%v", response, err)
	}
}

func TestFetchRawGzipAndLimits(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		writer := gzip.NewWriter(w)
		_, _ = writer.Write([]byte(strings.Repeat("x", 1000)))
		_ = writer.Close()
	}))
	defer server.Close()
	response, err := FetchRaw(context.Background(), Request{URL: server.URL, MaxBodyBytes: 512, MaxDecodedBytes: 2000})
	if err != nil || len(response.Body) != 1000 {
		t.Fatalf("gzip response bytes=%d err=%v", len(response.Body), err)
	}
	_, err = FetchRaw(context.Background(), Request{URL: server.URL, MaxBodyBytes: 512, MaxDecodedBytes: 100})
	if fetchErrorCode(err) != "decoded_body_too_large" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFetchRawBodyLimitAndRedirectLoop(t *testing.T) {
	large := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = fmt.Fprint(w, strings.Repeat("x", 100)) }))
	defer large.Close()
	_, err := FetchRaw(context.Background(), Request{URL: large.URL, MaxBodyBytes: 10, MaxDecodedBytes: 100})
	if fetchErrorCode(err) != "body_too_large" {
		t.Fatalf("unexpected limit error: %v", err)
	}

	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, r.URL.String(), http.StatusFound) }))
	defer redirect.Close()
	_, err = FetchRaw(context.Background(), Request{URL: redirect.URL, MaxRedirects: 2})
	if fetchErrorCode(err) != "too_many_redirects" {
		t.Fatalf("unexpected redirect error: %v", err)
	}
}
