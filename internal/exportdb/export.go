package exportdb

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/aggregate"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/emit"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

const DefaultShardSize = 2000

type Service struct {
	DB        *sql.DB
	Output    string
	ShardSize int
}

type Report struct {
	GeneratedAt time.Time      `json:"generated_at"`
	Candidates  int            `json:"candidate_count"`
	Stable      int            `json:"stable_count"`
	Collections map[string]int `json:"collections"`
	Files       int            `json:"file_count"`
	Bytes       int64          `json:"output_bytes"`
}

type item struct {
	node *node.Node
	meta store.ExportMeta
}

func (s Service) Run(ctx context.Context) (Report, error) {
	if s.DB == nil {
		return Report{}, errors.New("export database is required")
	}
	var nodes []*node.Node
	var metadata []store.ExportMeta
	for offset := 0; ; offset += 10000 {
		page, meta, err := store.ListExportable(ctx, s.DB, store.ExportQuery{Limit: 10000, Offset: offset})
		if err != nil {
			return Report{}, err
		}
		nodes = append(nodes, page...)
		metadata = append(metadata, meta...)
		if len(page) < 10000 {
			break
		}
	}
	startedAt := time.Now().UTC()
	runID, err := store.StartExportRun(ctx, s.DB, newRunUUID(), startedAt)
	if err != nil {
		return Report{}, fmt.Errorf("start export run: %w", err)
	}
	report, err := Generate(s.Output, nodes, metadata, s.ShardSize)
	if err != nil {
		_ = store.FailExportRun(ctx, s.DB, runID, time.Now().UTC(), err)
		return Report{}, err
	}
	members := exportMembers(nodes, metadata)
	if err := store.CompleteExportRun(ctx, s.DB, runID, time.Now().UTC(), report.Candidates, report.Candidates,
		report.Files, report.Bytes, report, members); err != nil {
		_ = store.FailExportRun(ctx, s.DB, runID, time.Now().UTC(), err)
		return Report{}, fmt.Errorf("complete export run: %w", err)
	}
	return report, nil
}

// Generate renders every collection into an invisible staging directory and
// only replaces managed output entries after all emitters have succeeded.
func Generate(root string, nodes []*node.Node, metadata []store.ExportMeta, shardSize int) (Report, error) {
	if len(nodes) != len(metadata) {
		return Report{}, fmt.Errorf("nodes and metadata differ: %d != %d", len(nodes), len(metadata))
	}
	if root == "" {
		return Report{}, errors.New("export output directory is required")
	}
	if shardSize <= 0 {
		shardSize = DefaultShardSize
	}
	if shardSize > DefaultShardSize {
		return Report{}, fmt.Errorf("shard size %d exceeds %d", shardSize, DefaultShardSize)
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return Report{}, err
	}
	lock := filepath.Join(root, ".export.lock")
	if err := os.Mkdir(lock, 0o700); err != nil {
		return Report{}, fmt.Errorf("acquire export lock: %w", err)
	}
	defer os.Remove(lock)

	staging := filepath.Join(root, ".next")
	if err := os.RemoveAll(staging); err != nil {
		return Report{}, err
	}
	if err := os.Mkdir(staging, 0o755); err != nil {
		return Report{}, err
	}
	defer os.RemoveAll(staging)

	all := make([]item, len(nodes))
	for index := range nodes {
		all[index] = item{node: nodes[index], meta: metadata[index]}
	}
	stable := selectItems(all, func(value item) bool {
		return value.meta.Grade == "S" || value.meta.Grade == "A" || value.meta.Grade == "B"
	})
	collections := map[string][]item{"all-verified": all, "stable": stable}
	addGroups(collections, "protocol", all, func(value item) string { return value.node.Protocol })
	addGroups(collections, "country", all, func(value item) string { return value.meta.Country })
	addGroups(collections, "network", all, func(value item) string { return value.meta.NetworkClass })
	for _, directory := range []string{"all-verified", "stable", "protocol", "country", "network"} {
		if err := os.MkdirAll(filepath.Join(staging, directory), 0o755); err != nil {
			return Report{}, err
		}
	}

	report := Report{GeneratedAt: time.Now().UTC(), Candidates: len(all), Stable: len(stable), Collections: make(map[string]int)}
	keys := make([]string, 0, len(collections))
	for key := range collections {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		report.Collections[key] = len(collections[key])
		if err := writeCollection(staging, key, collections[key], shardSize, &report); err != nil {
			return Report{}, fmt.Errorf("render collection %s: %w", key, err)
		}
	}
	legacy := stable
	if len(legacy) > DefaultShardSize {
		legacy = legacy[:DefaultShardSize]
	}
	if err := writeLegacy(staging, legacy, &report); err != nil {
		return Report{}, err
	}
	status, err := json.MarshalIndent(buildLegacyStatus(all, legacy, report.GeneratedAt), "", "  ")
	if err != nil {
		return Report{}, err
	}
	if err := writeOutput(staging, "status.json", append(status, '\n'), &report); err != nil {
		return Report{}, err
	}
	manifestReport := report
	manifestReport.Files++
	var manifest []byte
	for {
		var err error
		manifest, err = json.MarshalIndent(manifestReport, "", "  ")
		if err != nil {
			return Report{}, err
		}
		manifest = append(manifest, '\n')
		total := report.Bytes + int64(len(manifest))
		if manifestReport.Bytes == total {
			break
		}
		manifestReport.Bytes = total
	}
	if err := writeOutput(staging, "manifest.json", manifest, &report); err != nil {
		return Report{}, err
	}
	if err := promote(staging, root); err != nil {
		return Report{}, err
	}
	return report, nil
}

func buildLegacyStatus(all, selected []item, generatedAt time.Time) aggregate.Summary {
	// Database export has no per-run fetch stage. Preserve the public JSON
	// contract by reporting the current exportable snapshot for fetch/alive/verified.
	status := aggregate.Summary{
		TotalFetched:    len(all),
		TotalAlive:      len(all),
		TotalVerified:   len(all),
		TotalSelected:   len(selected),
		BySource:        make(map[string]int),
		ByProtocol:      make(map[string]int),
		ByCountry:       make(map[string]int),
		GeneratedAtUnix: generatedAt.Unix(),
	}
	latencies := make([]int, 0, len(selected))
	for _, value := range selected {
		if value.node.SourceName != "" {
			status.BySource[value.node.SourceName]++
		}
		status.ByProtocol[value.node.Protocol]++
		if value.meta.Country != "" {
			status.ByCountry[value.meta.Country]++
		}
		if value.node.LatencyMS > 0 {
			latencies = append(latencies, value.node.LatencyMS)
		}
	}
	if len(latencies) > 0 {
		sort.Ints(latencies)
		status.MinLatencyMS = latencies[0]
		middle := len(latencies) / 2
		status.MedianLatencyMS = latencies[middle]
		if len(latencies)%2 == 0 {
			status.MedianLatencyMS = (latencies[middle-1] + latencies[middle]) / 2
		}
	}
	return status
}

func addGroups(collections map[string][]item, prefix string, values []item, key func(item) string) {
	for _, value := range values {
		segment := safeSegment(key(value))
		if segment == "" {
			continue
		}
		name := prefix + "/" + segment
		collections[name] = append(collections[name], value)
	}
}

func selectItems(values []item, keep func(item) bool) []item {
	result := make([]item, 0, len(values))
	for _, value := range values {
		if keep(value) {
			result = append(result, value)
		}
	}
	return result
}

func shard(values []item, size int) [][]item {
	if size <= 0 {
		return nil
	}
	if len(values) == 0 {
		return [][]item{nil}
	}
	result := make([][]item, 0, (len(values)+size-1)/size)
	for start := 0; start < len(values); start += size {
		end := start + size
		if end > len(values) {
			end = len(values)
		}
		result = append(result, values[start:end])
	}
	return result
}

func writeCollection(root, collection string, values []item, shardSize int, report *Report) error {
	for index, part := range shard(values, shardSize) {
		base := filepath.Join(collection, fmt.Sprintf("%%s-%04d", index+1))
		if err := writeFormats(root, base, part, report); err != nil {
			return err
		}
	}
	return nil
}

func writeLegacy(root string, values []item, report *Report) error {
	return writeFormats(root, "%s", values, report)
}

func writeFormats(root, pattern string, values []item, report *Report) error {
	clashNodes := supported(values, "clash")
	clash, err := emit.Clash(clashNodes)
	if err != nil {
		return err
	}
	if err := writeOutput(root, fmt.Sprintf(pattern, "clash")+".yaml", []byte(clash), report); err != nil {
		return err
	}
	singNodes := supported(values, "singbox")
	singbox, err := emit.Singbox(singNodes)
	if err != nil {
		return err
	}
	if err := writeOutput(root, fmt.Sprintf(pattern, "singbox")+".json", []byte(singbox), report); err != nil {
		return err
	}
	v2rayNodes := supported(values, "v2ray")
	return writeOutput(root, fmt.Sprintf(pattern, "v2ray-base64")+".txt", []byte(emit.V2RayBase64(v2rayNodes)), report)
}

func supported(values []item, format string) []*node.Node {
	result := make([]*node.Node, 0, len(values))
	for _, value := range values {
		n := value.node
		supported := false
		switch format {
		case "clash", "v2ray":
			supported = n.Protocol == node.ProtoVLESS || n.Protocol == node.ProtoVMess || n.Protocol == node.ProtoTrojan || n.Protocol == node.ProtoSS || n.Protocol == node.ProtoHysteria2
		case "singbox":
			supported = emit.SingboxOutbound(n, "probe") != nil
		}
		if supported {
			result = append(result, n)
		}
	}
	return result
}

func writeOutput(root, name string, body []byte, report *Report) error {
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, body, 0o644); err != nil {
		return err
	}
	report.Files++
	report.Bytes += int64(len(body))
	return nil
}

func promote(staging, root string) error {
	entries, err := os.ReadDir(staging)
	if err != nil {
		return err
	}
	backupRoot := filepath.Join(root, ".previous")
	if err := os.RemoveAll(backupRoot); err != nil {
		return err
	}
	if err := os.Mkdir(backupRoot, 0o755); err != nil {
		return err
	}
	defer os.RemoveAll(backupRoot)
	var backedUp []string
	for _, entry := range entries {
		target := filepath.Join(root, entry.Name())
		if _, err := os.Lstat(target); err == nil {
			if err := os.Rename(target, filepath.Join(backupRoot, entry.Name())); err != nil {
				for _, previous := range backedUp {
					_ = os.Rename(filepath.Join(backupRoot, previous), filepath.Join(root, previous))
				}
				return err
			}
			backedUp = append(backedUp, entry.Name())
		} else if !os.IsNotExist(err) {
			for _, previous := range backedUp {
				_ = os.Rename(filepath.Join(backupRoot, previous), filepath.Join(root, previous))
			}
			return err
		}
	}
	var moved []string
	for _, entry := range entries {
		name := entry.Name()
		if err := os.Rename(filepath.Join(staging, name), filepath.Join(root, name)); err != nil {
			for _, published := range moved {
				_ = os.RemoveAll(filepath.Join(root, published))
			}
			for _, previous := range entries {
				_ = os.Rename(filepath.Join(backupRoot, previous.Name()), filepath.Join(root, previous.Name()))
			}
			return fmt.Errorf("publish %s: %w", name, err)
		}
		moved = append(moved, name)
	}
	return nil
}

func safeSegment(value string) string {
	value = strings.TrimSpace(value)
	var result strings.Builder
	for _, char := range value {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '-' || char == '_' {
			result.WriteRune(char)
		}
	}
	return result.String()
}

func exportMembers(nodes []*node.Node, metadata []store.ExportMeta) []store.ExportMember {
	all := make([]item, len(nodes))
	for index := range nodes {
		all[index] = item{node: nodes[index], meta: metadata[index]}
	}
	collections := map[string][]item{"all-verified": all}
	collections["stable"] = selectItems(all, func(value item) bool {
		return value.meta.Grade == "S" || value.meta.Grade == "A" || value.meta.Grade == "B"
	})
	addGroups(collections, "protocol", all, func(value item) string { return value.node.Protocol })
	addGroups(collections, "country", all, func(value item) string { return value.meta.Country })
	addGroups(collections, "network", all, func(value item) string { return value.meta.NetworkClass })
	keys := make([]string, 0, len(collections))
	for key := range collections {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var result []store.ExportMember
	for _, key := range keys {
		for index, value := range collections[key] {
			result = append(result, store.ExportMember{ConfigID: value.meta.ConfigID, Collection: key,
				Rank: index + 1, Score: value.meta.Score, Grade: value.meta.Grade, Reason: value.meta.Reason})
		}
	}
	return result
}

func newRunUUID() string {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		stamp := time.Now().UTC().UnixNano()
		return fmt.Sprintf("00000000-0000-4000-8000-%012x", uint64(stamp)&0xffffffffffff)
	}
	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", value[0:4], value[4:6], value[6:8], value[8:10], value[10:16])
}
