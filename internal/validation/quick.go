package validation

import (
	"context"
	"crypto/tls"
	"net"
	"strconv"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

type TLSRequest struct {
	Address, ServerName string
	Insecure            bool
	Timeout             time.Duration
}

type NetworkProber interface {
	TCP(context.Context, string, time.Duration) (time.Duration, error)
	TLS(context.Context, TLSRequest) (time.Duration, error)
	ProxyHandshake(context.Context, *node.Node, time.Duration) (time.Duration, error)
}

type HostResolver interface {
	LookupHost(context.Context, string) ([]string, error)
}

type StageResult struct {
	Passed                  bool
	Route                   Route
	DNSMS, ConnectMS, TLSMS int
	ErrorCode, ErrorSummary string
}

type QuickChecker struct {
	Prober   NetworkProber
	Resolver HostResolver
	Timeout  time.Duration
}

func (c QuickChecker) Check(ctx context.Context, n *node.Node) StageResult {
	result := StageResult{Route: RouteFor(n)}
	if !n.Valid() {
		result.ErrorCode, result.ErrorSummary = "invalid_config", "configuration is missing required fields"
		return result
	}
	if c.Timeout <= 0 {
		c.Timeout = 5 * time.Second
	}
	if c.Prober == nil {
		c.Prober = standardNetworkProber{}
	}
	if c.Resolver == nil {
		c.Resolver = net.DefaultResolver
	}
	if net.ParseIP(n.Server) == nil {
		started := time.Now()
		if _, err := c.Resolver.LookupHost(ctx, n.Server); err != nil {
			result.ErrorCode, result.ErrorSummary = "dns_failed", boundedValidationError(err)
			return result
		}
		result.DNSMS = elapsedMilliseconds(started)
	}
	address := net.JoinHostPort(n.Server, strconv.Itoa(n.Port))
	switch result.Route.Kind {
	case TransportQUIC, TransportWireGuard:
		result.Passed = true
		return result
	case TransportProxyHandshake:
		duration, err := c.Prober.ProxyHandshake(ctx, n, c.Timeout)
		result.ConnectMS = durationMS(duration)
		if err != nil {
			result.ErrorCode, result.ErrorSummary = "proxy_handshake_failed", boundedValidationError(err)
			return result
		}
	default:
		duration, err := c.Prober.TCP(ctx, address, c.Timeout)
		result.ConnectMS = durationMS(duration)
		if err != nil {
			result.ErrorCode, result.ErrorSummary = "tcp_failed", boundedValidationError(err)
			return result
		}
	}
	if result.Route.PrecheckTLS {
		serverName := n.SNI
		if serverName == "" {
			serverName = n.Server
		}
		duration, err := c.Prober.TLS(ctx, TLSRequest{Address: address, ServerName: serverName, Insecure: n.Insecure, Timeout: c.Timeout})
		result.TLSMS = durationMS(duration)
		if err != nil {
			result.ErrorCode, result.ErrorSummary = "tls_failed", boundedValidationError(err)
			return result
		}
	}
	result.Passed = true
	return result
}

type standardNetworkProber struct{}

func (standardNetworkProber) TCP(ctx context.Context, address string, timeout time.Duration) (time.Duration, error) {
	started := time.Now()
	connection, err := (&net.Dialer{Timeout: timeout}).DialContext(ctx, "tcp", address)
	if err == nil {
		_ = connection.Close()
	}
	return time.Since(started), err
}

func (standardNetworkProber) TLS(ctx context.Context, request TLSRequest) (time.Duration, error) {
	started := time.Now()
	dialer := tls.Dialer{NetDialer: &net.Dialer{Timeout: request.Timeout}, Config: &tls.Config{ServerName: request.ServerName, InsecureSkipVerify: request.Insecure, MinVersion: tls.VersionTLS12}}
	connection, err := dialer.DialContext(ctx, "tcp", request.Address)
	if err == nil {
		_ = connection.Close()
	}
	return time.Since(started), err
}

func (standardNetworkProber) ProxyHandshake(ctx context.Context, n *node.Node, timeout time.Duration) (time.Duration, error) {
	address := net.JoinHostPort(n.Server, strconv.Itoa(n.Port))
	duration, err := (standardNetworkProber{}).TCP(ctx, address, timeout)
	if err != nil {
		return duration, err
	}
	return duration, nil
}

func boundedValidationError(err error) string {
	message := err.Error()
	if len(message) > 1024 {
		message = message[:1024]
	}
	return message
}

func elapsedMilliseconds(start time.Time) int { return durationMS(time.Since(start)) }
func durationMS(value time.Duration) int {
	if value <= 0 {
		return 0
	}
	return int(value / time.Millisecond)
}
