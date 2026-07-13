package node

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ParseURI dispatches a proxy URI to the protocol-specific parser.
func ParseURI(uri string) (*Node, error) {
	uri = strings.TrimSpace(uri)
	switch {
	case strings.HasPrefix(uri, "vless://"):
		return parseVLESS(uri)
	case strings.HasPrefix(uri, "vmess://"):
		return parseVMess(uri)
	case strings.HasPrefix(uri, "trojan://"):
		return parseTrojan(uri)
	case strings.HasPrefix(uri, "ss://"):
		return parseShadowsocks(uri)
	case strings.HasPrefix(uri, "ssr://"):
		return parseShadowsocksR(uri)
	case strings.HasPrefix(uri, "hysteria://"):
		return parseHysteria(uri)
	case strings.HasPrefix(uri, "hysteria2://"), strings.HasPrefix(uri, "hy2://"):
		return parseHysteria2(uri)
	case strings.HasPrefix(uri, "tuic://"):
		return parseTUIC(uri)
	case strings.HasPrefix(uri, "wireguard://"), strings.HasPrefix(uri, "wg://"):
		return parseWireGuard(uri)
	case strings.HasPrefix(uri, "socks4://"):
		return parseUserProxy(uri, ProtoSOCKS4)
	case strings.HasPrefix(uri, "socks5://"), strings.HasPrefix(uri, "socks://"):
		return parseUserProxy(uri, ProtoSOCKS5)
	case strings.HasPrefix(uri, "http://"):
		return parseUserProxy(uri, ProtoHTTP)
	case strings.HasPrefix(uri, "https://"):
		return parseUserProxy(uri, ProtoHTTPS)
	}
	return nil, fmt.Errorf("unsupported scheme: %.40q", uri)
}

// parseVLESS parses vless://UUID@host:port?params#name
func parseVLESS(raw string) (*Node, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return nil, fmt.Errorf("vless port: %w", err)
	}
	q := u.Query()
	n := &Node{
		Protocol:    ProtoVLESS,
		Name:        decodeFragment(u.Fragment),
		Server:      u.Hostname(),
		Port:        port,
		UUID:        userName(u),
		Network:     q.Get("type"),
		Security:    q.Get("security"),
		SNI:         q.Get("sni"),
		ALPN:        q.Get("alpn"),
		Fingerprint: q.Get("fp"),
		PublicKey:   q.Get("pbk"),
		ShortID:     q.Get("sid"),
		SpiderX:     q.Get("spx"),
		Flow:        q.Get("flow"),
		Path:        q.Get("path"),
		Host:        q.Get("host"),
		ServiceName: q.Get("serviceName"),
	}
	if n.Network == "" {
		n.Network = "tcp"
	}
	if n.Security == "" {
		n.Security = "none"
	}
	return n, nil
}

// parseVMess parses vmess://base64(JSON).
type vmessJSON struct {
	V    any    `json:"v"`
	PS   string `json:"ps"`
	Add  string `json:"add"`
	Port any    `json:"port"`
	ID   string `json:"id"`
	Aid  any    `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	TLS  string `json:"tls"`
	SNI  string `json:"sni"`
	ALPN string `json:"alpn"`
	FP   string `json:"fp"`
	SCY  string `json:"scy"`
}

func parseVMess(raw string) (*Node, error) {
	payload := strings.TrimPrefix(raw, "vmess://")
	decoded, err := B64Decode(payload)
	if err != nil {
		return nil, fmt.Errorf("vmess base64: %w", err)
	}
	var vj vmessJSON
	if err := json.Unmarshal(decoded, &vj); err != nil {
		return nil, fmt.Errorf("vmess json: %w", err)
	}
	port, err := anyToInt(vj.Port)
	if err != nil {
		return nil, fmt.Errorf("vmess port: %w", err)
	}
	aid, _ := anyToInt(vj.Aid)
	security := "none"
	if vj.TLS == "tls" || vj.TLS == "reality" {
		security = vj.TLS
	}
	if vj.Net == "" {
		vj.Net = "tcp"
	}
	return &Node{
		Protocol:    ProtoVMess,
		Name:        vj.PS,
		Server:      vj.Add,
		Port:        port,
		UUID:        vj.ID,
		AlterID:     aid,
		Network:     vj.Net,
		Security:    security,
		SNI:         firstNonEmpty(vj.SNI, vj.Host),
		ALPN:        vj.ALPN,
		Fingerprint: vj.FP,
		Path:        vj.Path,
		Host:        vj.Host,
	}, nil
}

// parseTrojan parses trojan://password@host:port?params#name
func parseTrojan(raw string) (*Node, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return nil, fmt.Errorf("trojan port: %w", err)
	}
	q := u.Query()
	return &Node{
		Protocol:    ProtoTrojan,
		Name:        decodeFragment(u.Fragment),
		Server:      u.Hostname(),
		Port:        port,
		Password:    userName(u),
		Network:     firstNonEmpty(q.Get("type"), "tcp"),
		Security:    "tls",
		SNI:         q.Get("sni"),
		ALPN:        q.Get("alpn"),
		Fingerprint: q.Get("fp"),
		Path:        q.Get("path"),
		Host:        q.Get("host"),
		Insecure:    q.Get("allowInsecure") == "1",
	}, nil
}

// parseShadowsocks handles both legacy (base64 of method:pass@host:port) and
// modern (ss://base64(method:pass)@host:port) forms.
func parseShadowsocks(raw string) (*Node, error) {
	raw = strings.TrimPrefix(raw, "ss://")
	var name string
	if i := strings.Index(raw, "#"); i >= 0 {
		name = decodeFragment(raw[i+1:])
		raw = raw[:i]
	}
	// Strip query string if any (rare for ss).
	if i := strings.Index(raw, "?"); i >= 0 {
		raw = raw[:i]
	}

	var method, password, hostport string

	if strings.Contains(raw, "@") {
		// modern form: base64(method:pass)@host:port
		parts := strings.SplitN(raw, "@", 2)
		creds, err := B64Decode(parts[0])
		if err != nil {
			return nil, fmt.Errorf("ss creds base64: %w", err)
		}
		mp := strings.SplitN(string(creds), ":", 2)
		if len(mp) != 2 {
			return nil, errors.New("ss creds malformed")
		}
		method, password = mp[0], mp[1]
		hostport = parts[1]
	} else {
		// legacy: whole thing base64
		decoded, err := B64Decode(raw)
		if err != nil {
			return nil, fmt.Errorf("ss legacy base64: %w", err)
		}
		s := string(decoded)
		at := strings.LastIndex(s, "@")
		if at < 0 {
			return nil, errors.New("ss legacy: no @")
		}
		mp := strings.SplitN(s[:at], ":", 2)
		if len(mp) != 2 {
			return nil, errors.New("ss legacy creds malformed")
		}
		method, password = mp[0], mp[1]
		hostport = s[at+1:]
	}

	host, portStr, err := splitHostPort(hostport)
	if err != nil {
		return nil, fmt.Errorf("ss hostport: %w", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("ss port: %w", err)
	}
	return &Node{
		Protocol: ProtoSS,
		Name:     name,
		Server:   host,
		Port:     port,
		Cipher:   method,
		Password: password,
	}, nil
}

// parseHysteria2 parses hy2://password@host:port?params#name
func parseHysteria2(raw string) (*Node, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return nil, fmt.Errorf("hy2 port: %w", err)
	}
	q := u.Query()
	return &Node{
		Protocol: ProtoHysteria2,
		Name:     decodeFragment(u.Fragment),
		Server:   u.Hostname(),
		Port:     port,
		Password: userName(u),
		Security: "tls",
		SNI:      q.Get("sni"),
		Insecure: q.Get("insecure") == "1",
	}, nil
}

func parseShadowsocksR(raw string) (*Node, error) {
	decoded, err := B64Decode(strings.TrimPrefix(raw, "ssr://"))
	if err != nil {
		return nil, fmt.Errorf("ssr base64: %w", err)
	}
	main := strings.SplitN(string(decoded), "/?", 2)[0]
	parts := strings.Split(main, ":")
	if len(parts) < 6 {
		return nil, errors.New("ssr fields missing")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("ssr port: %w", err)
	}
	password, err := B64Decode(parts[5])
	if err != nil {
		return nil, fmt.Errorf("ssr password: %w", err)
	}
	return &Node{Protocol: ProtoSSR, Server: parts[0], Port: port, Cipher: parts[3], Password: string(password), Extra: map[string]string{"protocol": parts[2], "obfs": parts[4]}}, nil
}

func parseHysteria(raw string) (*Node, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return nil, fmt.Errorf("hysteria port: %w", err)
	}
	password := userName(u)
	if password == "" {
		password = u.Query().Get("auth")
	}
	return &Node{Protocol: ProtoHysteria, Name: decodeFragment(u.Fragment), Server: u.Hostname(), Port: port, Password: password, Security: "tls", SNI: firstNonEmpty(u.Query().Get("sni"), u.Query().Get("peer")), Insecure: u.Query().Get("insecure") == "1"}, nil
}

func parseTUIC(raw string) (*Node, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return nil, fmt.Errorf("tuic port: %w", err)
	}
	username, password := userCredentials(u)
	return &Node{Protocol: ProtoTUIC, Name: decodeFragment(u.Fragment), Server: u.Hostname(), Port: port, UUID: username, Password: password, Security: "tls", SNI: u.Query().Get("sni"), ALPN: u.Query().Get("alpn"), Insecure: u.Query().Get("allow_insecure") == "1"}, nil
}

func parseWireGuard(raw string) (*Node, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return nil, fmt.Errorf("wireguard port: %w", err)
	}
	q := u.Query()
	return &Node{Protocol: ProtoWireGuard, Name: decodeFragment(u.Fragment), Server: u.Hostname(), Port: port, Password: userName(u), PublicKey: firstNonEmpty(q.Get("publickey"), q.Get("public_key")), Extra: map[string]string{"address": q.Get("address"), "reserved": q.Get("reserved"), "pre_shared_key": q.Get("presharedkey")}}, nil
}

func parseUserProxy(raw, protocol string) (*Node, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return nil, fmt.Errorf("%s port: %w", protocol, err)
	}
	username, password := userCredentials(u)
	security := "none"
	if protocol == ProtoHTTPS {
		security = "tls"
	}
	return &Node{Protocol: protocol, Name: decodeFragment(u.Fragment), Server: u.Hostname(), Port: port, Username: username, Password: password, Security: security}, nil
}

// ---- helpers ----

// B64Decode accepts both standard and URL-safe base64, padded or unpadded.
func B64Decode(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	// pad to multiple of 4
	if m := len(s) % 4; m != 0 {
		s += strings.Repeat("=", 4-m)
	}
	if data, err := base64.StdEncoding.DecodeString(s); err == nil {
		return data, nil
	}
	return base64.URLEncoding.DecodeString(s)
}

func anyToInt(v any) (int, error) {
	switch x := v.(type) {
	case float64:
		return int(x), nil
	case int:
		return x, nil
	case string:
		return strconv.Atoi(x)
	}
	return 0, fmt.Errorf("not an int: %v", v)
}

func decodeFragment(s string) string {
	if s == "" {
		return ""
	}
	if decoded, err := url.QueryUnescape(s); err == nil {
		return decoded
	}
	return s
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func userName(u *url.URL) string {
	if u.User == nil {
		return ""
	}
	return u.User.Username()
}

func userCredentials(u *url.URL) (string, string) {
	if u.User == nil {
		return "", ""
	}
	password, _ := u.User.Password()
	return u.User.Username(), password
}

func splitHostPort(s string) (string, string, error) {
	// simple IPv4/hostname; IPv6 support can be added later
	i := strings.LastIndex(s, ":")
	if i < 0 {
		return "", "", errors.New("no port")
	}
	return s[:i], s[i+1:], nil
}
