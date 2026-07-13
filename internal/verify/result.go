package verify

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/proxy"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

type TargetResult struct {
	URLHash, ErrorCode                                   string
	OK                                                   bool
	StatusCode, DNSMS, ConnectMS, TLSMS, TTFBMS, TotalMS int
}

type Request struct {
	Targets        []string
	Timeout        time.Duration
	BasePort       int
	SingBoxBin     string
	StartupTimeout time.Duration
}

type Engine struct{}

func (Engine) Verify(ctx context.Context, n *node.Node, request Request) Result {
	result := Result{Node: n, Protocol: n.Protocol, Attempts: len(request.Targets)}
	if _, err := BuildOutbound(n, "out-0"); err != nil {
		result.ErrorCode, result.ErrorSummary = "config_invalid", boundedResultError(err)
		return result
	}
	if request.Timeout <= 0 {
		request.Timeout = 8 * time.Second
	}
	if len(request.Targets) == 0 {
		request.Targets = []string{"http://www.gstatic.com/generate_204", "https://www.cloudflare.com/cdn-cgi/trace"}
		result.Attempts = len(request.Targets)
	}
	if request.BasePort == 0 {
		request.BasePort = availablePort()
	}
	config := applyDefaults(Config{BasePort: request.BasePort, BatchSize: 1, Concurrency: 1,
		TimeoutMS: int(request.Timeout / time.Millisecond), Targets: request.Targets,
		SingBoxBin: request.SingBoxBin, StartupTimeout: request.StartupTimeout})
	batch, err := buildBatchConfig([]*node.Node{n}, request.BasePort)
	if err != nil {
		result.ErrorCode, result.ErrorSummary = "config_build_failed", boundedResultError(err)
		return result
	}
	started := time.Now()
	process, err := startSingbox(ctx, batch, config)
	result.StartMS = durationToMS(time.Since(started))
	if err != nil {
		result.ErrorCode, result.ErrorSummary = "singbox_rejected", boundedResultError(err)
		return result
	}
	result.ConfigAccepted, result.ProxyStarted = true, true
	defer process.stop()
	client, err := socksHTTPClient(request.BasePort, request.Timeout)
	if err != nil {
		result.ErrorCode, result.ErrorSummary = "proxy_client_failed", boundedResultError(err)
		return result
	}
	for _, target := range request.Targets {
		targetResult, body := probeTarget(ctx, client, target)
		result.Targets = append(result.Targets, targetResult)
		if targetResult.OK {
			parseExitTrace(body, &result)
		}
	}
	finalizeEngineResult(&result)
	return result
}

func socksHTTPClient(port int, timeout time.Duration) (*http.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)), nil, &net.Dialer{Timeout: timeout})
	if err != nil {
		return nil, err
	}
	contextDialer, ok := dialer.(proxy.ContextDialer)
	if !ok {
		return nil, errors.New("SOCKS dialer does not support contexts")
	}
	transport := &http.Transport{DialContext: contextDialer.DialContext, TLSHandshakeTimeout: timeout,
		ResponseHeaderTimeout: timeout, DisableKeepAlives: true, MaxIdleConns: 1}
	return &http.Client{Transport: transport, Timeout: timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}, nil
}

func probeTarget(ctx context.Context, client *http.Client, target string) (TargetResult, []byte) {
	result := TargetResult{URLHash: targetURLHash(target)}
	if parsed, err := url.Parse(target); err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		result.ErrorCode = "invalid_target"
		return result, nil
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		result.ErrorCode = "invalid_target"
		return result, nil
	}
	var dnsStart, connectStart, tlsStart, requestStart time.Time
	trace := &httptrace.ClientTrace{
		DNSStart:             func(httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone:              func(httptrace.DNSDoneInfo) { result.DNSMS = elapsedSince(dnsStart) },
		ConnectStart:         func(_, _ string) { connectStart = time.Now() },
		ConnectDone:          func(_, _ string, _ error) { result.ConnectMS = elapsedSince(connectStart) },
		TLSHandshakeStart:    func() { tlsStart = time.Now() },
		TLSHandshakeDone:     func(tls.ConnectionState, error) { result.TLSMS = elapsedSince(tlsStart) },
		GotFirstResponseByte: func() { result.TTFBMS = elapsedSince(requestStart) },
	}
	request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))
	request.Header.Set("User-Agent", "free-vpn-subscriptions-validator/2.0")
	requestStart = time.Now()
	response, err := client.Do(request)
	result.TotalMS = elapsedSince(requestStart)
	if err != nil {
		result.ErrorCode = classifyProbeError(err)
		return result, nil
	}
	defer response.Body.Close()
	result.StatusCode = response.StatusCode
	body, readErr := io.ReadAll(io.LimitReader(response.Body, 32<<10))
	if readErr != nil {
		result.ErrorCode = "response_read_failed"
		return result, nil
	}
	if response.StatusCode < 200 || response.StatusCode >= 400 {
		result.ErrorCode = "http_status"
		return result, body
	}
	result.OK = true
	return result, body
}

func finalizeEngineResult(result *Result) {
	var latencies []int
	for _, target := range result.Targets {
		if target.OK {
			result.Successes++
			latencies = append(latencies, target.TotalMS)
		}
	}
	sort.Ints(latencies)
	result.HTTPMedianMS = median(latencies)
	result.Passed = result.Attempts > 0 && result.Successes == result.Attempts
	result.PartialSuccess = result.Successes > 0 && !result.Passed
	if result.Successes == 0 && result.ErrorCode == "" {
		result.ErrorCode = "all_targets_failed"
	}
}

func parseExitTrace(body []byte, result *Result) {
	for _, line := range strings.Split(string(body), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch key {
		case "ip":
			if net.ParseIP(value) != nil {
				result.ExitIP = value
			}
		case "loc":
			if len(value) == 2 {
				result.ExitCountry = strings.ToUpper(value)
			}
		}
	}
}

func classifyProbeError(err error) string {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return "timeout"
	}
	var dnsError *net.DNSError
	if errors.As(err, &dnsError) {
		return "dns_failed"
	}
	if strings.Contains(strings.ToLower(err.Error()), "refused") {
		return "connection_refused"
	}
	return "request_failed"
}

func targetURLHash(target string) string {
	digest := sha256.Sum256([]byte(target))
	return hex.EncodeToString(digest[:8])
}

func availablePort() int {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 20000
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func boundedResultError(err error) string {
	message := err.Error()
	if len(message) > 1024 {
		return message[:1024]
	}
	return message
}

func elapsedSince(start time.Time) int {
	if start.IsZero() {
		return 0
	}
	return durationToMS(time.Since(start))
}

func durationToMS(duration time.Duration) int {
	if duration <= 0 {
		return 0
	}
	return int(duration / time.Millisecond)
}
