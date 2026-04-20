package verify

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"golang.org/x/net/proxy"
)

// outcome summarizes one node's performance for a single round.
type outcome struct {
	passed    bool // every configured target returned OK within timeout
	successes int  // number of targets that returned OK
	medianMS  int  // median of successful request latencies, 0 if none
}

// probeViaSocks dials the given local SOCKS5 port (provided by our sing-box)
// and runs the configured targets through it. A target counts as OK when the
// HTTP response status is in [200, 400) (covers 200/204).
func probeViaSocks(ctx context.Context, port int, cfg Config) outcome {
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	dialer, err := proxy.SOCKS5("tcp", addr, nil, &net.Dialer{Timeout: time.Duration(cfg.TimeoutMS) * time.Millisecond})
	if err != nil {
		return outcome{}
	}
	contextDialer, ok := dialer.(proxy.ContextDialer)
	if !ok {
		return outcome{}
	}

	transport := &http.Transport{
		DialContext:           contextDialer.DialContext,
		TLSHandshakeTimeout:   time.Duration(cfg.TimeoutMS) * time.Millisecond,
		ResponseHeaderTimeout: time.Duration(cfg.TimeoutMS) * time.Millisecond,
		DisableKeepAlives:     true,
		MaxIdleConns:          1,
	}
	defer transport.CloseIdleConnections()

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(cfg.TimeoutMS) * time.Millisecond,
		// Do not follow redirects — a redirect through a free node often
		// masks a captive portal; fail closed.
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var latencies []int
	successes := 0
	for _, target := range cfg.Targets {
		if ctx.Err() != nil {
			break
		}
		if _, err := url.Parse(target); err != nil {
			continue
		}
		ok, ms := hit(ctx, client, target)
		if ok {
			successes++
			latencies = append(latencies, ms)
		}
	}

	sort.Ints(latencies)
	return outcome{
		passed:    successes == len(cfg.Targets),
		successes: successes,
		medianMS:  median(latencies),
	}
}

func hit(ctx context.Context, client *http.Client, target string) (bool, int) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return false, 0
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (fvs-probe)")

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return false, 0
	}
	defer resp.Body.Close()
	// Drain a small amount so the connection is cleanly reusable / torn down.
	_, _ = io.CopyN(io.Discard, resp.Body, 4096)

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return false, 0
	}
	return true, int(time.Since(start) / time.Millisecond)
}

func median(sorted []int) int {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	mid := n / 2
	if n%2 == 1 {
		return sorted[mid]
	}
	return (sorted[mid-1] + sorted[mid]) / 2
}
