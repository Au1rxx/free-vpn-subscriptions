// Package httpfallback provides narrowly scoped HTTP transport fallbacks.
package httpfallback

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const telegramFallbackHost = "telegram.me"

type lookupHostFunc func(context.Context, string) ([]string, error)
type dialContextFunc func(context.Context, string, string) (net.Conn, error)

// Do executes request and, only when t.me does not exist in DNS, retries the
// same URL through addresses resolved for telegram.me. Keeping the original
// URL and TLS server name means Telegram redirects back to t.me remain usable.
func Do(client *http.Client, request *http.Request) (*http.Response, error) {
	dialer := &net.Dialer{}
	return doWithNetwork(client, request, net.DefaultResolver.LookupHost, dialer.DialContext)
}

func doWithNetwork(client *http.Client, request *http.Request, lookup lookupHostFunc, aliasDial dialContextFunc) (*http.Response, error) {
	if client == nil || request == nil {
		return nil, fmt.Errorf("HTTP client and request are required")
	}
	response, err := client.Do(request)
	if err == nil || !eligible(request, err) {
		return response, err
	}
	addresses, lookupErr := lookup(request.Context(), telegramFallbackHost)
	if lookupErr != nil {
		return nil, fmt.Errorf("t.me DNS lookup failed (%v); telegram.me fallback failed: %w", err, lookupErr)
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("t.me DNS lookup failed (%v); telegram.me fallback failed: no addresses", err)
	}
	fallbackClient, fallbackErr := clientWithTelegramAlias(client, addresses, aliasDial)
	if fallbackErr != nil {
		return nil, fmt.Errorf("t.me DNS lookup failed (%v); telegram.me fallback failed: %w", err, fallbackErr)
	}
	retry := request.Clone(request.Context())
	response, retryErr := fallbackClient.Do(retry)
	if retryErr != nil {
		return nil, fmt.Errorf("t.me DNS lookup failed (%v); telegram.me fallback failed: %w", err, retryErr)
	}
	return response, nil
}

func clientWithTelegramAlias(client *http.Client, addresses []string, aliasDial dialContextFunc) (*http.Client, error) {
	var transport *http.Transport
	switch current := client.Transport.(type) {
	case nil:
		transport = http.DefaultTransport.(*http.Transport).Clone()
	case *http.Transport:
		transport = current.Clone()
	default:
		return nil, fmt.Errorf("HTTP transport type %T cannot install DNS alias", client.Transport)
	}
	regularDial := transport.DialContext
	if regularDial == nil {
		regularDial = (&net.Dialer{}).DialContext
	}
	transport.DialTLSContext = nil
	transport.DisableKeepAlives = true
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, splitErr := net.SplitHostPort(address)
		if splitErr != nil || !strings.EqualFold(host, "t.me") {
			return regularDial(ctx, network, address)
		}
		var lastErr error
		for _, aliasAddress := range addresses {
			connection, dialErr := aliasDial(ctx, network, net.JoinHostPort(aliasAddress, port))
			if dialErr == nil {
				return connection, nil
			}
			lastErr = dialErr
		}
		return nil, lastErr
	}
	fallback := *client
	fallback.Transport = transport
	return &fallback, nil
}

func eligible(request *http.Request, err error) bool {
	if request.Method != http.MethodGet || request.Body != nil || request.URL == nil || request.URL.Port() != "" ||
		!strings.EqualFold(request.URL.Hostname(), "t.me") {
		return false
	}
	var dnsError *net.DNSError
	return errors.As(err, &dnsError) && dnsError.IsNotFound
}
