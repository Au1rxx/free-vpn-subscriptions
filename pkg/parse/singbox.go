package parse

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// SingBox converts supported client outbounds from a sing-box JSON config.
func SingBox(body []byte) Result {
	result := Result{Format: FormatSingBox}
	var root struct {
		Outbounds []map[string]any `json:"outbounds"`
	}
	if err := json.Unmarshal(body, &root); err != nil {
		result.Errors = append(result.Errors, newEntryError(0, "invalid_json", "", body, err))
		return result
	}
	for index, outbound := range root.Outbounds {
		n := singBoxOutbound(outbound)
		if n == nil {
			continue
		}
		if !n.Valid() {
			result.Errors = append(result.Errors, newEntryError(index+1, "invalid_node", n.Protocol, []byte(n.Protocol), fmt.Errorf("sing-box outbound is missing required fields")))
			continue
		}
		result.Nodes = append(result.Nodes, n)
	}
	return result
}

func singBoxOutbound(value map[string]any) *node.Node {
	typeName := strings.ToLower(jsonString(value["type"]))
	n := &node.Node{Name: jsonString(value["tag"]), Server: jsonString(value["server"]), Port: jsonInt(value["server_port"]), Extra: map[string]string{}}
	switch typeName {
	case "vless":
		n.Protocol, n.UUID, n.Flow = node.ProtoVLESS, jsonString(value["uuid"]), jsonString(value["flow"])
	case "vmess":
		n.Protocol, n.UUID, n.AlterID = node.ProtoVMess, jsonString(value["uuid"]), jsonInt(value["alter_id"])
	case "trojan":
		n.Protocol, n.Password = node.ProtoTrojan, jsonString(value["password"])
	case "shadowsocks":
		n.Protocol, n.Cipher, n.Password = node.ProtoSS, jsonString(value["method"]), jsonString(value["password"])
	case "hysteria":
		n.Protocol, n.Password = node.ProtoHysteria, firstJSON(value, "auth_str", "auth")
	case "hysteria2":
		n.Protocol, n.Password = node.ProtoHysteria2, jsonString(value["password"])
	case "tuic":
		n.Protocol, n.UUID, n.Password = node.ProtoTUIC, jsonString(value["uuid"]), jsonString(value["password"])
	case "wireguard":
		n.Protocol, n.Password, n.PublicKey = node.ProtoWireGuard, jsonString(value["private_key"]), jsonString(value["peer_public_key"])
		n.Extra["address"], n.Extra["reserved"] = jsonValueString(value["local_address"]), jsonValueString(value["reserved"])
	case "socks", "socks5":
		n.Protocol, n.Username, n.Password = node.ProtoSOCKS5, jsonString(value["username"]), jsonString(value["password"])
	case "http":
		n.Protocol, n.Username, n.Password = node.ProtoHTTP, jsonString(value["username"]), jsonString(value["password"])
	default:
		return nil
	}
	if tls := jsonMap(value["tls"]); tls != nil && jsonBool(tls["enabled"]) {
		n.Security, n.SNI, n.Insecure = "tls", jsonString(tls["server_name"]), jsonBool(tls["insecure"])
		n.ALPN = jsonValueString(tls["alpn"])
	}
	if transport := jsonMap(value["transport"]); transport != nil {
		n.Network, n.Path = jsonString(transport["type"]), jsonString(transport["path"])
		n.Host, n.ServiceName = jsonString(transport["host"]), jsonString(transport["service_name"])
	}
	return n
}

func jsonString(value any) string {
	text, _ := value.(string)
	return text
}

func jsonInt(value any) int {
	switch number := value.(type) {
	case float64:
		return int(number)
	case json.Number:
		integer, _ := number.Int64()
		return int(integer)
	}
	return 0
}

func jsonBool(value any) bool {
	boolean, _ := value.(bool)
	return boolean
}

func jsonMap(value any) map[string]any {
	object, _ := value.(map[string]any)
	return object
}

func jsonValueString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case []any:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			items = append(items, fmt.Sprint(item))
		}
		return strings.Join(items, ",")
	default:
		if value != nil {
			return fmt.Sprint(value)
		}
	}
	return ""
}

func firstJSON(value map[string]any, keys ...string) string {
	for _, key := range keys {
		if text := jsonString(value[key]); text != "" {
			return text
		}
	}
	return ""
}
