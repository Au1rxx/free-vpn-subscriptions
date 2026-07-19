package httpfallback

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

func TestDoFallsBackFromTMeNXDomainWithoutBreakingRedirects(t *testing.T) {
	var paths []string
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		paths = append(paths, request.URL.RequestURI())
		if request.Host != "t.me" {
			t.Errorf("request host=%q, want t.me", request.Host)
		}
		if request.Header.Get("If-None-Match") != `"same"` {
			t.Errorf("fallback dropped conditional request header: %v", request.Header)
		}
		if request.URL.Path == "/s/channel" {
			http.Redirect(writer, request, "https://t.me/channel?q=1", http.StatusFound)
			return
		}
		_, _ = io.WriteString(writer, "preview")
	}))
	defer server.Close()

	initialDials := 0
	transport := server.Client().Transport.(*http.Transport).Clone()
	transport.TLSClientConfig = transport.TLSClientConfig.Clone()
	transport.TLSClientConfig.InsecureSkipVerify = true // Test server certificate is not issued for t.me.
	transport.DialContext = func(context.Context, string, string) (net.Conn, error) {
		initialDials++
		return nil, &net.DNSError{Name: "t.me", Err: "no such host", IsNotFound: true}
	}
	client := &http.Client{Transport: transport}
	request, err := http.NewRequest(http.MethodGet, "https://t.me/s/channel?q=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("If-None-Match", `"same"`)
	aliasDials := 0
	response, err := doWithNetwork(client, request,
		func(_ context.Context, host string) ([]string, error) {
			if host != telegramFallbackHost {
				t.Errorf("lookup host=%q", host)
			}
			return []string{"203.0.113.10"}, nil
		},
		func(ctx context.Context, network, address string) (net.Conn, error) {
			aliasDials++
			if address != "203.0.113.10:443" {
				t.Errorf("alias address=%q", address)
				return nil, errors.New("unexpected alias address")
			}
			return (&net.Dialer{}).DialContext(ctx, network, strings.TrimPrefix(server.URL, "https://"))
		})
	if err != nil || response.StatusCode != http.StatusOK {
		t.Fatalf("response=%v err=%v", response, err)
	}
	response.Body.Close()
	if initialDials != 1 || aliasDials != 2 {
		t.Fatalf("initial dials=%d alias dials=%d", initialDials, aliasDials)
	}
	if strings.Join(paths, ",") != "/s/channel?q=1,/channel?q=1" {
		t.Fatalf("paths=%v", paths)
	}
	if response.Request.URL.Hostname() != "t.me" {
		t.Fatalf("final URL host=%q", response.Request.URL.Hostname())
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
	fallbackErr := errors.New("fallback unavailable")
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(context.Context, string, string) (net.Conn, error) {
		return nil, &net.DNSError{Name: "t.me", Err: "no such host", IsNotFound: true}
	}
	client := &http.Client{Transport: transport}
	request, _ := http.NewRequest(http.MethodGet, "https://t.me/s/channel", nil)
	_, err := doWithNetwork(client, request,
		func(context.Context, string) ([]string, error) { return []string{"203.0.113.10"}, nil },
		func(context.Context, string, string) (net.Conn, error) { return nil, fallbackErr })
	if !errors.Is(err, fallbackErr) || !strings.Contains(err.Error(), "telegram.me fallback failed") {
		t.Fatalf("unexpected fallback error: %v", err)
	}
}
