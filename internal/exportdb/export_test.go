package exportdb

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/aggregate"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

func TestShardBoundsAtTenThousandNodes(t *testing.T) {
	items := make([]item, 10000)
	shards := shard(items, 2000)
	if len(shards) != 5 {
		t.Fatalf("shards=%d, want 5", len(shards))
	}
	for index, part := range shards {
		if len(part) > 2000 {
			t.Fatalf("shard %d contains %d nodes", index, len(part))
		}
	}
}

func TestGenerateClassifiedAndLegacyOutputs(t *testing.T) {
	root := t.TempDir()
	nodes := []*node.Node{
		{Protocol: node.ProtoVLESS, Server: "one.example", Port: 443, UUID: "a", Country: "US", LatencyMS: 10},
		{Protocol: node.ProtoTUIC, Server: "two.example", Port: 443, UUID: "b", Password: "p", Country: "JP", LatencyMS: 30},
		{Protocol: node.ProtoSS, Server: "three.example", Port: 8388, Cipher: "aes-128-gcm", Password: "p", Country: "US", LatencyMS: 50},
	}
	metadata := []store.ExportMeta{
		{ConfigID: 1, Grade: "A", Score: 88, Country: "US", NetworkClass: "cloud"},
		{ConfigID: 2, Grade: "B", Score: 70, Country: "JP", NetworkClass: "hosting"},
		{ConfigID: 3, Grade: "D", Score: 40, Country: "US", NetworkClass: "isp"},
	}
	report, err := Generate(root, nodes, metadata, 2)
	if err != nil {
		t.Fatal(err)
	}
	if report.Candidates != 3 || report.Stable != 2 || report.Files == 0 {
		t.Fatalf("report=%+v", report)
	}
	for _, path := range []string{
		"clash.yaml", "singbox.json", "v2ray-base64.txt", "manifest.json", "status.json",
		"stable/clash-0001.yaml", "all-verified/clash-0001.yaml",
		"protocol/vless/singbox-0001.json", "country/US/v2ray-base64-0001.txt",
		"network/cloud/clash-0001.yaml",
	} {
		if _, err := os.Stat(filepath.Join(root, path)); err != nil {
			t.Errorf("missing %s: %v", path, err)
		}
	}
	body, err := os.ReadFile(filepath.Join(root, "v2ray-base64.txt"))
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(decoded), "three.example") {
		t.Fatal("legacy stable output contains grade D node")
	}
	statusBody, err := os.ReadFile(filepath.Join(root, "status.json"))
	if err != nil {
		t.Fatal(err)
	}
	var status aggregate.Summary
	if err := json.Unmarshal(statusBody, &status); err != nil {
		t.Fatal(err)
	}
	if status.TotalFetched != 3 || status.TotalAlive != 3 || status.TotalVerified != 3 || status.TotalSelected != 2 {
		t.Fatalf("status totals=%+v", status)
	}
	if status.ByProtocol[node.ProtoVLESS] != 1 || status.ByProtocol[node.ProtoTUIC] != 1 ||
		status.ByCountry["US"] != 1 || status.ByCountry["JP"] != 1 {
		t.Fatalf("status breakdowns=%+v", status)
	}
	if status.MinLatencyMS != 10 || status.MedianLatencyMS != 20 || status.GeneratedAtUnix == 0 {
		t.Fatalf("status quality=%+v", status)
	}
	if _, err := os.Stat(filepath.Join(root, ".next")); !os.IsNotExist(err) {
		t.Fatalf("staging directory remains: %v", err)
	}
}

func TestGenerateRejectsMisalignedMetadata(t *testing.T) {
	_, err := Generate(t.TempDir(), []*node.Node{{}}, nil, 2000)
	if err == nil {
		t.Fatal("expected metadata alignment error")
	}
}

func TestExportMembersIncludeDetailedCollections(t *testing.T) {
	nodes := []*node.Node{{Protocol: node.ProtoVLESS, Country: "US"}}
	metadata := []store.ExportMeta{{ConfigID: 9, Grade: "A", Score: 85, Country: "US", NetworkClass: "cloud", Reason: "verified_a"}}
	members := exportMembers(nodes, metadata)
	collections := map[string]bool{}
	for _, member := range members {
		collections[member.Collection] = true
		if member.ConfigID != 9 || member.Rank != 1 || member.Score != 85 || member.Reason != "verified_a" {
			t.Fatalf("member=%+v", member)
		}
	}
	for _, want := range []string{"all-verified", "stable", "protocol/vless", "country/US", "network/cloud"} {
		if !collections[want] {
			t.Errorf("missing collection %s", want)
		}
	}
	if uuid := newRunUUID(); len(uuid) != 36 || uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		t.Fatalf("invalid run UUID %q", uuid)
	}
}
