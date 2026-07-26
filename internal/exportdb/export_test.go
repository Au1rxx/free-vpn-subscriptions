package exportdb

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	var scope struct {
		SchemaVersion   int    `json:"schema_version"`
		DataSource      string `json:"data_source"`
		StatisticsScope string `json:"statistics_scope"`
	}
	if err := json.Unmarshal(statusBody, &scope); err != nil {
		t.Fatal(err)
	}
	if scope.SchemaVersion != 2 || scope.DataSource != "database" || scope.StatisticsScope != "exportable_snapshot" {
		t.Fatalf("status scope=%+v", scope)
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

func TestGenerateExposesSiteRenderInputs(t *testing.T) {
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
	report, err := Generate(t.TempDir(), nodes, metadata, 2)
	if err != nil {
		t.Fatal(err)
	}
	// Only the stable grades (A, B) are published, so the site renders those.
	if len(report.Selected) != 2 {
		t.Fatalf("selected=%d", len(report.Selected))
	}
	if report.Summary.TotalSelected != 2 || report.Summary.TotalVerified != 3 {
		t.Fatalf("summary=%+v", report.Summary)
	}
	if report.Summary.ByCountry["US"] != 1 || report.Summary.ByCountry["JP"] != 1 {
		t.Fatalf("by_country=%v", report.Summary.ByCountry)
	}
	if report.Summary.GeneratedAtUnix == 0 {
		t.Fatal("generated_at_unix is zero")
	}
}

func TestManifestExcludesSiteRenderInputs(t *testing.T) {
	root := t.TempDir()
	nodes := []*node.Node{{Protocol: node.ProtoVLESS, Server: "one.example", Port: 443, UUID: "a", Country: "US", LatencyMS: 10}}
	metadata := []store.ExportMeta{{ConfigID: 1, Grade: "A", Score: 88, Country: "US", NetworkClass: "cloud"}}
	if _, err := Generate(root, nodes, metadata, 2); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"Summary", "Selected", "by_country", "total_selected"} {
		if strings.Contains(string(body), field) {
			t.Fatalf("manifest leaked %q: %s", field, body)
		}
	}
}

// Shard contents must not depend on the ranking of the input: quality scores
// are re-measured every run, and a rank-sensitive layout rewrites whole shards
// on every publish even when the node set is unchanged.
func TestCollectionShardsAreIndependentOfInputOrder(t *testing.T) {
	build := func() ([]*node.Node, []store.ExportMeta) {
		var nodes []*node.Node
		var metadata []store.ExportMeta
		for i := 1; i <= 6; i++ {
			nodes = append(nodes, &node.Node{
				Protocol: node.ProtoVLESS,
				Server:   fmt.Sprintf("host%02d.example", i),
				Port:     443,
				UUID:     fmt.Sprintf("uuid-%02d", i),
				Name:     fmt.Sprintf("vless-%02d", i),
				Country:  "US",
			})
			metadata = append(metadata, store.ExportMeta{
				ConfigID: uint64(i), Grade: "A", Score: 90, Country: "US", NetworkClass: "cloud",
			})
		}
		return nodes, metadata
	}

	ranked, rankedMeta := build()
	rootA := t.TempDir()
	if _, err := Generate(rootA, ranked, rankedMeta, 2); err != nil {
		t.Fatal(err)
	}

	// Same nodes, reversed ranking — as if every score had been re-measured.
	reversed, reversedMeta := build()
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
		reversedMeta[i], reversedMeta[j] = reversedMeta[j], reversedMeta[i]
	}
	rootB := t.TempDir()
	if _, err := Generate(rootB, reversed, reversedMeta, 2); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{
		"all-verified/clash-0001.yaml", "all-verified/clash-0002.yaml",
		"all-verified/singbox-0001.json", "country/US/clash-0001.yaml",
	} {
		a, err := os.ReadFile(filepath.Join(rootA, name))
		if err != nil {
			t.Fatal(err)
		}
		b, err := os.ReadFile(filepath.Join(rootB, name))
		if err != nil {
			t.Fatal(err)
		}
		if string(a) != string(b) {
			t.Errorf("%s differs when only the input ranking changed", name)
		}
	}
}
