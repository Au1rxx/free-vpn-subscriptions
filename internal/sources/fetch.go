// Package sources performs bounded HTTP collection of subscription feeds.
package sources

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/parse"
)

const (
	defaultMaxBodyBytes    = int64(20 << 20)
	defaultMaxDecodedBytes = int64(64 << 20)
)

// Request contains conditional request metadata and hard resource limits.
type Request struct {
	URL, ETag, LastModified, UserAgent string
	Timeout                            time.Duration
	MaxBodyBytes, MaxDecodedBytes      int64
	MaxRedirects                       int
}

// Response is the decoded response plus metadata needed for persistence.
type Response struct {
	StatusCode                                                 int
	FinalURL, ETag, LastModified, ContentType, ContentEncoding string
	Body                                                       []byte
	SHA256                                                     [32]byte
	FetchedAt                                                  time.Time
	Duration                                                   time.Duration
}

// FetchError exposes a stable code without including response bodies.
type FetchError struct {
	Code string
	Err  error
}

func (e *FetchError) Error() string { return e.Code + ": " + e.Err.Error() }
func (e *FetchError) Unwrap() error { return e.Err }

func fetchErrorCode(err error) string {
	if typed, ok := err.(*FetchError); ok {
		return typed.Code
	}
	return ""
}

// FetchRaw retrieves one response with redirect, compressed and decoded size
// limits. HTTP 304 is a successful metadata-only response.
func FetchRaw(ctx context.Context, request Request) (Response, error) {
	started := time.Now()
	request = requestDefaults(request)
	if err := validateHTTPURL(request.URL); err != nil {
		return Response{}, &FetchError{Code: "invalid_url", Err: err}
	}
	redirects := 0
	client := &http.Client{
		Timeout: request.Timeout,
		CheckRedirect: func(next *http.Request, _ []*http.Request) error {
			redirects++
			if redirects > request.MaxRedirects {
				return &FetchError{Code: "too_many_redirects", Err: fmt.Errorf("redirect limit %d exceeded", request.MaxRedirects)}
			}
			if err := validateHTTPURL(next.URL.String()); err != nil {
				return &FetchError{Code: "invalid_redirect", Err: err}
			}
			return nil
		},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, request.URL, nil)
	if err != nil {
		return Response{}, &FetchError{Code: "build_request", Err: err}
	}
	req.Header.Set("User-Agent", request.UserAgent)
	req.Header.Set("Accept-Encoding", "gzip")
	if request.ETag != "" {
		req.Header.Set("If-None-Match", request.ETag)
	}
	if request.LastModified != "" {
		req.Header.Set("If-Modified-Since", request.LastModified)
	}
	resp, err := client.Do(req)
	if err != nil {
		if typed := findFetchError(err); typed != nil {
			return Response{}, typed
		}
		return Response{}, &FetchError{Code: "http_failed", Err: err}
	}
	defer resp.Body.Close()
	result := Response{
		StatusCode: resp.StatusCode, FinalURL: resp.Request.URL.String(),
		ETag: resp.Header.Get("ETag"), LastModified: resp.Header.Get("Last-Modified"),
		ContentType: resp.Header.Get("Content-Type"), ContentEncoding: resp.Header.Get("Content-Encoding"),
		FetchedAt: time.Now().UTC(), Duration: time.Since(started),
	}
	if resp.StatusCode == http.StatusNotModified {
		return result, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return result, &FetchError{Code: "http_status", Err: fmt.Errorf("HTTP status %d", resp.StatusCode)}
	}
	raw, exceeded, err := readLimited(resp.Body, request.MaxBodyBytes)
	if err != nil {
		return result, &FetchError{Code: "read_failed", Err: err}
	}
	if exceeded {
		return result, &FetchError{Code: "body_too_large", Err: fmt.Errorf("compressed response exceeds %d bytes", request.MaxBodyBytes)}
	}
	decoded := raw
	if strings.EqualFold(strings.TrimSpace(result.ContentEncoding), "gzip") {
		reader, err := gzip.NewReader(bytes.NewReader(raw))
		if err != nil {
			return result, &FetchError{Code: "invalid_gzip", Err: err}
		}
		defer reader.Close()
		decoded, exceeded, err = readLimited(reader, request.MaxDecodedBytes)
		if err != nil {
			return result, &FetchError{Code: "decode_failed", Err: err}
		}
	} else if int64(len(decoded)) > request.MaxDecodedBytes {
		exceeded = true
	}
	if exceeded {
		return result, &FetchError{Code: "decoded_body_too_large", Err: fmt.Errorf("decoded response exceeds %d bytes", request.MaxDecodedBytes)}
	}
	result.Body, result.SHA256 = decoded, sha256.Sum256(decoded)
	return result, nil
}

func requestDefaults(request Request) Request {
	if request.Timeout <= 0 {
		request.Timeout = 30 * time.Second
	}
	if request.MaxBodyBytes <= 0 {
		request.MaxBodyBytes = defaultMaxBodyBytes
	}
	if request.MaxDecodedBytes <= 0 {
		request.MaxDecodedBytes = defaultMaxDecodedBytes
	}
	if request.MaxRedirects <= 0 {
		request.MaxRedirects = 5
	}
	if request.UserAgent == "" {
		request.UserAgent = "free-vpn-subscriptions/2.0 (+https://github.com/Au1rxx/free-vpn-subscriptions)"
	}
	return request
}

func validateHTTPURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return fmt.Errorf("only absolute HTTP/HTTPS URLs are allowed")
	}
	return nil
}

func readLimited(reader io.Reader, maximum int64) ([]byte, bool, error) {
	body, err := io.ReadAll(io.LimitReader(reader, maximum+1))
	if err != nil {
		return nil, false, err
	}
	if int64(len(body)) > maximum {
		return body[:maximum], true, nil
	}
	return body, false, nil
}

func findFetchError(err error) *FetchError {
	for current := err; current != nil; {
		if typed, ok := current.(*FetchError); ok {
			return typed
		}
		type unwrapper interface{ Unwrap() error }
		next, ok := current.(unwrapper)
		if !ok {
			break
		}
		current = next.Unwrap()
	}
	return nil
}

// Fetch retains the original aggregate API while using the bounded fetcher.
func Fetch(ctx context.Context, src config.Source, timeout time.Duration) ([]*node.Node, error) {
	if !src.Enabled {
		return nil, nil
	}
	response, err := FetchRaw(ctx, Request{URL: src.URL, Timeout: timeout})
	if err != nil {
		return nil, fmt.Errorf("source %q: %w", src.Name, err)
	}
	format := parse.Format(src.Format)
	result := parse.Parse(response.Body, format)
	if len(result.Nodes) == 0 && len(result.Errors) > 0 {
		return nil, fmt.Errorf("source %q: parse %s: %s", src.Name, result.Format, result.Errors[0].Code)
	}
	for _, n := range result.Nodes {
		n.SourceName = src.Name
	}
	return result.Nodes, nil
}
