// Package httpfallback provides narrowly scoped HTTP transport fallbacks.
package httpfallback

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const telegramFallbackHost = "telegram.me"

// Do executes request and retries a bodyless t.me GET through telegram.me
// exactly once when the original domain does not exist in DNS.
func Do(client *http.Client, request *http.Request) (*http.Response, error) {
	if client == nil || request == nil {
		return nil, fmt.Errorf("HTTP client and request are required")
	}
	response, err := client.Do(request)
	if err == nil || !eligible(request, err) {
		return response, err
	}
	retry := request.Clone(request.Context())
	copiedURL := *request.URL
	copiedURL.Host = telegramFallbackHost
	retry.URL = &copiedURL
	retry.Host = ""
	response, retryErr := client.Do(retry)
	if retryErr != nil {
		return nil, fmt.Errorf("t.me DNS lookup failed (%v); telegram.me fallback failed: %w", err, retryErr)
	}
	return response, nil
}

func eligible(request *http.Request, err error) bool {
	if request.Method != http.MethodGet || request.Body != nil || request.URL == nil || request.URL.Port() != "" ||
		!strings.EqualFold(request.URL.Hostname(), "t.me") {
		return false
	}
	var dnsError *net.DNSError
	return errors.As(err, &dnsError) && dnsError.IsNotFound
}
