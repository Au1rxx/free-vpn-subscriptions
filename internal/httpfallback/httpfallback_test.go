package httpfallback

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

func TestDoFallsBackFromTMeNXDomain(t *testing.T) {
	var hosts []string
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		hosts = append(hosts, request.URL.Hostname())
		if request.URL.Hostname() == "t.me" {
			return nil, &net.DNSError{Name: "t.me", Err: "no such host", IsNotFound: true}
		}
		if request.URL.Path != "/s/channel" || request.URL.RawQuery != "q=1" {
			t.Fatalf("fallback changed URL path or query: %s", request.URL)
		}
		if request.Header.Get("If-None-Match") != `"same"` {
			t.Fatalf("fallback dropped conditional request header: %v", request.Header)
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header),
			Body: io.NopCloser(strings.NewReader("preview")), Request: request}, nil
	})}
	request, err := http.NewRequest(http.MethodGet, "https://t.me/s/channel?q=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("If-None-Match", `"same"`)
	response, err := Do(client, request)
	if err != nil || response.StatusCode != http.StatusOK {
		t.Fatalf("response=%v err=%v", response, err)
	}
	response.Body.Close()
	if !reflect.DeepEqual(hosts, []string{"t.me", "telegram.me"}) {
		t.Fatalf("request hosts=%v", hosts)
	}
}

func TestDoRejectsIneligibleFallbacks(t *testing.T) {
	tests := []struct {
		name, method, rawURL string
		body                 io.Reader
		err                  error
	}{
		{name: "other host", method: http.MethodGet, rawURL: "https://example.test/s/channel",
			err: &net.DNSError{Name: "example.test", Err: "no such host", IsNotFound: true}},
		{name: "DNS timeout", method: http.MethodGet, rawURL: "https://t.me/s/channel",
			err: &net.DNSError{Name: "t.me", Err: "timeout", IsTimeout: true}},
		{name: "request body", method: http.MethodPost, rawURL: "https://t.me/s/channel",
			body: strings.NewReader("body"), err: &net.DNSError{Name: "t.me", Err: "no such host", IsNotFound: true}},
		{name: "nonstandard port", method: http.MethodGet, rawURL: "https://t.me:8443/s/channel",
			err: &net.DNSError{Name: "t.me", Err: "no such host", IsNotFound: true}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				calls++
				return nil, test.err
			})}
			request, err := http.NewRequestWithContext(context.Background(), test.method, test.rawURL, test.body)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := Do(client, request); err == nil {
				t.Fatal("ineligible request unexpectedly succeeded")
			}
			if calls != 1 {
				t.Fatalf("calls=%d, want 1", calls)
			}
		})
	}
}

func TestDoDoesNotRetryHTTPStatus(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{StatusCode: http.StatusNotFound, Header: make(http.Header),
			Body: io.NopCloser(strings.NewReader("missing")), Request: request}, nil
	})}
	request, _ := http.NewRequest(http.MethodGet, "https://t.me/s/channel", nil)
	response, err := Do(client, request)
	if err != nil || response.StatusCode != http.StatusNotFound {
		t.Fatalf("response=%v err=%v", response, err)
	}
	response.Body.Close()
	if calls != 1 {
		t.Fatalf("calls=%d, want 1", calls)
	}
}

func TestDoReportsFallbackFailure(t *testing.T) {
	calls := 0
	fallbackErr := errors.New("fallback unavailable")
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		calls++
		if request.URL.Hostname() == "t.me" {
			return nil, &net.DNSError{Name: "t.me", Err: "no such host", IsNotFound: true}
		}
		return nil, fallbackErr
	})}
	request, _ := http.NewRequest(http.MethodGet, "https://t.me/s/channel", nil)
	_, err := Do(client, request)
	if !errors.Is(err, fallbackErr) || !strings.Contains(err.Error(), "telegram.me fallback failed") {
		t.Fatalf("unexpected fallback error: %v", err)
	}
	if calls != 2 {
		t.Fatalf("calls=%d, want 2", calls)
	}
}
