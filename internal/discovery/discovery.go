// Package discovery finds public subscription source candidates without
// writing them to persistent storage.
package discovery

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/httpfallback"
)

const maxDiscoveryBody = int64(8 << 20)

// Seed is an adapter input. URL is used by page/API adapters and Query by
// code-search adapters. TokenFile is optional and never included in output.
type Seed struct {
	URL, Query, Name, TokenFile string
	Depth                       int
}

// Candidate is evidence that a public URL may contain subscriptions.
type Candidate struct {
	URL, Kind, Name, ParentURL, Evidence string
	Depth                                int
}

// Discoverer returns candidates and does not mutate the database.
type Discoverer interface {
	Discover(context.Context, Seed) ([]Candidate, error)
}

// Error carries stable retry information for scheduler backoff.
type Error struct {
	Code              string
	RetryAfterSeconds int
	Err               error
}

func (e *Error) Error() string { return e.Code + ": " + e.Err.Error() }
func (e *Error) Unwrap() error { return e.Err }

func fetch(ctx context.Context, client *http.Client, rawURL, token string) ([]byte, http.Header, error) {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, nil, &Error{Code: "invalid_url", Err: err}
	}
	req.Header.Set("User-Agent", "free-vpn-subscriptions-discovery/2.0")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	response, err := httpfallback.Do(client, req)
	if err != nil {
		return nil, nil, &Error{Code: "http_failed", Err: err}
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusTooManyRequests || response.StatusCode == http.StatusForbidden {
		retry, _ := strconv.Atoi(response.Header.Get("Retry-After"))
		return nil, response.Header, &Error{Code: "rate_limited", RetryAfterSeconds: retry, Err: fmt.Errorf("HTTP status %d", response.StatusCode)}
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, response.Header, &Error{Code: "http_status", Err: fmt.Errorf("HTTP status %d", response.StatusCode)}
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, maxDiscoveryBody+1))
	if err != nil {
		return nil, response.Header, &Error{Code: "read_failed", Err: err}
	}
	if int64(len(body)) > maxDiscoveryBody {
		return nil, response.Header, &Error{Code: "body_too_large", Err: fmt.Errorf("response exceeds %d bytes", maxDiscoveryBody)}
	}
	return body, response.Header, nil
}

func readOptionalCredential(path string) string {
	if path == "" {
		return ""
	}
	body, err := os.ReadFile(path)
	if err != nil || len(body) > 16<<10 {
		return ""
	}
	return strings.TrimSpace(string(body))
}

func deduplicate(candidates []Candidate) []Candidate {
	seen := make(map[string]bool, len(candidates))
	result := make([]Candidate, 0, len(candidates))
	for _, candidate := range candidates {
		parsed, err := url.Parse(strings.TrimSpace(candidate.URL))
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			continue
		}
		parsed.Fragment = ""
		candidate.URL = parsed.String()
		candidate.Evidence = boundedEvidence(candidate.Evidence)
		if !seen[candidate.URL] {
			seen[candidate.URL] = true
			result = append(result, candidate)
		}
	}
	return result
}

func boundedEvidence(value string) string {
	if len(value) > 512 {
		return value[:512]
	}
	return value
}
