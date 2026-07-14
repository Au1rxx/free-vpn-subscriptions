package parse

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// Parse detects or honors a subscription format while retaining individual
// failures instead of discarding them silently.
func Parse(body []byte, hint Format) Result {
	format := hint
	if format == "" || format == FormatAuto {
		format = detectFormat(body)
	}
	result := Result{Format: format}
	switch format {
	case FormatURIList:
		return parseURIText(body, format)
	case FormatBase64:
		decoded, err := node.B64Decode(strings.TrimSpace(string(body)))
		if err != nil {
			result.Errors = append(result.Errors, newEntryError(0, "invalid_base64", "", body, err))
			return result
		}
		return parseURIText(decoded, format)
	case FormatClash:
		nodes, err := Clash(body)
		if err != nil {
			result.Errors = append(result.Errors, newEntryError(0, "invalid_clash", "", body, err))
			return result
		}
		result.Nodes = nodes
		return result
	case FormatSingBox, FormatXray:
		if format == FormatSingBox {
			return SingBox(body)
		}
		return Xray(body)
	default:
		result.Errors = append(result.Errors, newEntryError(0, "unsupported_format", "", body, fmt.Errorf("unsupported format %q", format)))
		return result
	}
}

func detectFormat(body []byte) Format {
	trimmed := strings.TrimSpace(string(body))
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		var root map[string]any
		if json.Unmarshal(body, &root) == nil {
			if outbounds, ok := root["outbounds"].([]any); ok && len(outbounds) > 0 {
				if first, ok := outbounds[0].(map[string]any); ok {
					if _, ok := first["protocol"]; ok {
						return FormatXray
					}
				}
				return FormatSingBox
			}
		}
	}
	if strings.Contains(trimmed, "proxies:") {
		return FormatClash
	}
	if decoded, err := node.B64Decode(trimmed); err == nil && containsProxyScheme(string(decoded)) {
		return FormatBase64
	}
	return FormatURIList
}

func parseURIText(body []byte, format Format) Result {
	result := Result{Format: format}
	for index, raw := range strings.Split(string(body), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if discoveredURL(line) && !looksLikeHTTPProxy(line) {
			result.DiscoveredURLs = append(result.DiscoveredURLs, line)
			continue
		}
		n, err := node.ParseURI(line)
		if err != nil {
			result.Errors = append(result.Errors, newEntryError(index+1, errorCode(err), schemeOf(line), []byte(line), err))
			continue
		}
		if !n.Valid() {
			result.Errors = append(result.Errors, newEntryError(index+1, "invalid_node", schemeOf(line), []byte(line), fmt.Errorf("required fields are missing")))
			continue
		}
		result.Nodes = append(result.Nodes, n)
	}
	return result
}

func discoveredURL(value string) bool {
	u, err := url.Parse(value)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

func looksLikeHTTPProxy(value string) bool {
	u, err := url.Parse(value)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Port() == "" {
		return false
	}
	return (u.EscapedPath() == "" || u.EscapedPath() == "/") && u.RawQuery == ""
}

func containsProxyScheme(value string) bool {
	for _, prefix := range []string{"vless://", "vmess://", "trojan://", "ss://", "hysteria2://", "hy2://"} {
		if strings.Contains(value, prefix) {
			return true
		}
	}
	return false
}

func schemeOf(value string) string {
	if index := strings.Index(value, "://"); index > 0 && index <= 24 {
		return strings.ToLower(value[:index])
	}
	return ""
}

func errorCode(err error) string {
	if strings.Contains(err.Error(), "unsupported scheme") {
		return "unsupported_scheme"
	}
	return "parse_failed"
}

func newEntryError(line int, code, scheme string, sample []byte, err error) EntryError {
	digest := sha256.Sum256(sample)
	message := err.Error()
	for len(message) > 256 {
		_, size := utf8.DecodeLastRuneInString(message)
		message = message[:len(message)-size]
	}
	return EntryError{Line: line, Code: code, Scheme: scheme, SampleHash: hex.EncodeToString(digest[:8]), Message: message}
}
