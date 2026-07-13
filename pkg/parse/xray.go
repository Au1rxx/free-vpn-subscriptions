package parse

import (
	"encoding/json"
	"fmt"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// Xray converts supported client outbounds from an Xray/V2Ray JSON config.
func Xray(body []byte) Result {
	result := Result{Format: FormatXray}
	var root struct {
		Outbounds []map[string]any `json:"outbounds"`
	}
	if err := json.Unmarshal(body, &root); err != nil {
		result.Errors = append(result.Errors, newEntryError(0, "invalid_json", "", body, err))
		return result
	}
	for index, outbound := range root.Outbounds {
		for _, n := range xrayOutbound(outbound) {
			if !n.Valid() {
				result.Errors = append(result.Errors, newEntryError(index+1, "invalid_node", n.Protocol, []byte(n.Protocol), fmt.Errorf("Xray outbound is missing required fields")))
				continue
			}
			result.Nodes = append(result.Nodes, n)
		}
	}
	return result
}

func xrayOutbound(value map[string]any) []*node.Node {
	protocol := jsonString(value["protocol"])
	settings := jsonMap(value["settings"])
	stream := jsonMap(value["streamSettings"])
	var nodes []*node.Node
	if protocol == "vless" || protocol == "vmess" {
		for _, rawServer := range jsonSlice(settings["vnext"]) {
			server := jsonMap(rawServer)
			for _, rawUser := range jsonSlice(server["users"]) {
				user := jsonMap(rawUser)
				n := &node.Node{Protocol: protocol, Server: jsonString(server["address"]), Port: jsonInt(server["port"]), UUID: jsonString(user["id"]), AlterID: jsonInt(user["alterId"]), Flow: jsonString(user["flow"])}
				applyXrayStream(n, stream)
				nodes = append(nodes, n)
			}
		}
		return nodes
	}
	for _, rawServer := range jsonSlice(settings["servers"]) {
		server := jsonMap(rawServer)
		n := &node.Node{Server: firstJSON(server, "address", "server"), Port: jsonInt(server["port"])}
		switch protocol {
		case "trojan":
			n.Protocol, n.Password = node.ProtoTrojan, jsonString(server["password"])
		case "shadowsocks":
			n.Protocol, n.Cipher, n.Password = node.ProtoSS, jsonString(server["method"]), jsonString(server["password"])
		case "socks":
			n.Protocol = node.ProtoSOCKS5
			users := jsonSlice(server["users"])
			if len(users) > 0 {
				user := jsonMap(users[0])
				n.Username, n.Password = jsonString(user["user"]), jsonString(user["pass"])
			}
		case "http":
			n.Protocol = node.ProtoHTTP
			users := jsonSlice(server["users"])
			if len(users) > 0 {
				user := jsonMap(users[0])
				n.Username, n.Password = jsonString(user["user"]), jsonString(user["pass"])
			}
		default:
			continue
		}
		applyXrayStream(n, stream)
		nodes = append(nodes, n)
	}
	return nodes
}

func applyXrayStream(n *node.Node, stream map[string]any) {
	if stream == nil {
		return
	}
	n.Network, n.Security = jsonString(stream["network"]), jsonString(stream["security"])
	if tls := jsonMap(stream["tlsSettings"]); tls != nil {
		n.SNI, n.ALPN, n.Insecure = jsonString(tls["serverName"]), jsonValueString(tls["alpn"]), jsonBool(tls["allowInsecure"])
	}
	if reality := jsonMap(stream["realitySettings"]); reality != nil {
		n.PublicKey, n.ShortID, n.SpiderX = jsonString(reality["publicKey"]), jsonString(reality["shortId"]), jsonString(reality["spiderX"])
	}
	if ws := jsonMap(stream["wsSettings"]); ws != nil {
		n.Path = jsonString(ws["path"])
		if headers := jsonMap(ws["headers"]); headers != nil {
			n.Host = jsonString(headers["Host"])
		}
	}
	if grpc := jsonMap(stream["grpcSettings"]); grpc != nil {
		n.ServiceName = jsonString(grpc["serviceName"])
	}
}

func jsonSlice(value any) []any {
	items, _ := value.([]any)
	return items
}
