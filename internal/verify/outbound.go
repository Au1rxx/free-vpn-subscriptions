package verify

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// Outbound renders a node into a sing-box outbound JSON object.
// Returns nil for protocols we cannot translate (keeps the batch clean).
// Tag is assigned by the caller (e.g. "out-0").
//
// Deliberately a superset of internal/subscribe.Singbox — that function is
// shaped for end-user clients and skips WS/gRPC/Reality details we need
// here to actually forward traffic through the probe pipeline.
func BuildOutbound(n *node.Node, tag string) (map[string]any, error) {
	if n == nil || !n.Valid() {
		return nil, fmt.Errorf("invalid %s outbound configuration", protocolName(n))
	}
	if tag == "" {
		return nil, fmt.Errorf("outbound tag is required")
	}
	ob := map[string]any{
		"tag":         tag,
		"server":      n.Server,
		"server_port": n.Port,
	}
	switch n.Protocol {
	case node.ProtoVLESS:
		ob["type"] = "vless"
		ob["uuid"] = n.UUID
		if n.Flow != "" {
			ob["flow"] = n.Flow
		}
		if tls := buildTLS(n); tls != nil {
			ob["tls"] = tls
		}
		if tr := buildTransport(n); tr != nil {
			ob["transport"] = tr
		}
	case node.ProtoVMess:
		ob["type"] = "vmess"
		ob["uuid"] = n.UUID
		ob["alter_id"] = n.AlterID
		ob["security"] = orDefault(n.Cipher, "auto")
		if n.Security == "tls" {
			if tls := buildTLS(n); tls != nil {
				ob["tls"] = tls
			}
		}
		if tr := buildTransport(n); tr != nil {
			ob["transport"] = tr
		}
	case node.ProtoTrojan:
		ob["type"] = "trojan"
		ob["password"] = n.Password
		if tls := buildTLS(n); tls != nil {
			ob["tls"] = tls
		} else {
			ob["tls"] = map[string]any{"enabled": true, "server_name": orDefault(n.SNI, n.Server), "insecure": n.Insecure}
		}
		if tr := buildTransport(n); tr != nil {
			ob["transport"] = tr
		}
	case node.ProtoSS:
		ob["type"] = "shadowsocks"
		ob["method"] = n.Cipher
		ob["password"] = n.Password
	case node.ProtoHysteria2:
		ob["type"] = "hysteria2"
		ob["password"] = n.Password
		if tls := buildTLS(n); tls != nil {
			ob["tls"] = tls
		} else {
			ob["tls"] = map[string]any{"enabled": true, "server_name": orDefault(n.SNI, n.Server), "insecure": n.Insecure}
		}
		if kind := n.Extra["obfs"]; kind != "" {
			ob["obfs"] = map[string]any{"type": kind, "password": n.Extra["obfs_password"]}
		}
		if ports := splitCSV(n.Extra["server_ports"]); len(ports) > 0 {
			ob["server_ports"] = ports
		}
	case node.ProtoTUIC:
		ob["type"], ob["uuid"], ob["password"] = "tuic", n.UUID, n.Password
		ob["congestion_control"] = orDefault(n.Extra["congestion_control"], "bbr")
		ob["tls"] = orTLS(n)
	case node.ProtoWireGuard:
		addresses := splitCSV(n.Extra["address"])
		if len(addresses) == 0 {
			return nil, fmt.Errorf("wireguard local address is required")
		}
		peer := map[string]any{"address": n.Server, "port": n.Port, "public_key": n.PublicKey, "allowed_ips": []string{"0.0.0.0/0", "::/0"}}
		if key := n.Extra["pre_shared_key"]; key != "" {
			peer["pre_shared_key"] = key
		}
		if reserved := parseReserved(n.Extra["reserved"]); len(reserved) > 0 {
			peer["reserved"] = reserved
		}
		ob = map[string]any{"type": "wireguard", "tag": tag, "address": addresses, "private_key": n.Password, "peers": []map[string]any{peer}}
	case node.ProtoSOCKS4, node.ProtoSOCKS5:
		ob["type"], ob["version"], ob["username"], ob["password"] = "socks", "5", n.Username, n.Password
		if n.Protocol == node.ProtoSOCKS4 {
			ob["version"] = "4"
		}
	case node.ProtoHTTP, node.ProtoHTTPS:
		ob["type"], ob["username"], ob["password"] = "http", n.Username, n.Password
		if n.Protocol == node.ProtoHTTPS || n.Security == "tls" {
			ob["tls"] = orTLS(n)
		}
	default:
		return nil, fmt.Errorf("unsupported outbound protocol %q", n.Protocol)
	}
	return ob, nil
}

// IsEndpoint reports protocols represented in the sing-box endpoints section
// since the legacy outbound form was removed in sing-box 1.13.
func IsEndpoint(n *node.Node) bool { return n != nil && n.Protocol == node.ProtoWireGuard }

func buildOutbound(n *node.Node, tag string) map[string]any {
	outbound, err := BuildOutbound(n, tag)
	if err != nil {
		return nil
	}
	return outbound
}

func buildTLS(n *node.Node) map[string]any {
	if n.Security != "tls" && n.Security != "reality" {
		return nil
	}
	sni := orDefault(n.SNI, n.Server)
	tls := map[string]any{
		"enabled":     true,
		"server_name": sni,
		"insecure":    n.Insecure,
	}
	if n.ALPN != "" {
		tls["alpn"] = splitCSV(n.ALPN)
	}
	if n.Fingerprint != "" {
		tls["utls"] = map[string]any{"enabled": true, "fingerprint": n.Fingerprint}
	}
	if n.Security == "reality" {
		reality := map[string]any{"enabled": true}
		if n.PublicKey != "" {
			reality["public_key"] = n.PublicKey
		}
		if n.ShortID != "" {
			reality["short_id"] = n.ShortID
		}
		tls["reality"] = reality
		// Reality strictly requires utls; default to chrome when upstream omits.
		if _, ok := tls["utls"]; !ok {
			tls["utls"] = map[string]any{"enabled": true, "fingerprint": "chrome"}
		}
	}
	return tls
}

func orTLS(n *node.Node) map[string]any {
	if tls := buildTLS(n); tls != nil {
		return tls
	}
	return map[string]any{"enabled": true, "server_name": orDefault(n.SNI, n.Server), "insecure": n.Insecure}
}

func splitCSV(value string) []string {
	var result []string
	for _, item := range strings.Split(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			result = append(result, item)
		}
	}
	return result
}

func parseReserved(value string) []uint8 {
	var result []uint8
	for _, item := range splitCSV(value) {
		number, err := strconv.ParseUint(item, 10, 8)
		if err != nil {
			return nil
		}
		result = append(result, uint8(number))
	}
	return result
}

func protocolName(n *node.Node) string {
	if n == nil {
		return "nil"
	}
	return n.Protocol
}

func buildTransport(n *node.Node) map[string]any {
	switch n.Network {
	case "ws":
		t := map[string]any{"type": "ws"}
		if n.Path != "" {
			t["path"] = n.Path
		}
		if n.Host != "" {
			t["headers"] = map[string]any{"Host": n.Host}
		}
		return t
	case "grpc":
		if n.ServiceName == "" {
			return nil
		}
		return map[string]any{"type": "grpc", "service_name": n.ServiceName}
	}
	return nil
}

func orDefault(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
