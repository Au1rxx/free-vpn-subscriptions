// Package main provides the fnctl CLI — entry point for the aggregator.
//
// The only relevant command is `aggregate`: it fetches every enabled source,
// probes every node for TCP reachability, deduplicates and ranks the alive
// set, resolves each node's country via GeoIP, and writes output files
// (clash.yaml, singbox.json, v2ray-base64.txt, per-country variants,
// status.json) plus a generated README.md.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/aggregate"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/geoip"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/pages"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/probe"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/readme"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/sources"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/subscribe"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/verify"
)

// runDeadline caps a single aggregate run. The hourly cron has a hard 6h
// GitHub Actions ceiling; we aim to finish well inside that so SIGINT from
// the runner cleanly aborts in-flight probes instead of leaving goroutines
// to be killed.
const runDeadline = 30 * time.Minute

var cfgPath string

func main() {
	root := &cobra.Command{
		Use:   "fnctl",
		Short: "free-vpn-subscriptions aggregator CLI",
	}
	root.PersistentFlags().StringVarP(&cfgPath, "config", "c", "config.yaml", "path to configuration file")
	root.AddCommand(newAggregateCmd())
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func newAggregateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "aggregate",
		Short: "Fetch, probe, rank, and emit subscription outputs",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				die(err)
			}

			ctx, cancel := signal.NotifyContext(
				context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer cancel()
			ctx, cancelDeadline := context.WithTimeout(ctx, runDeadline)
			defer cancelDeadline()

			fetched := fetchAll(ctx, cfg)
			fmt.Printf("fetched %d nodes from %d sources\n", len(fetched), countEnabled(cfg.Sources))

			alive := probe.TCP(ctx, fetched,
				time.Duration(cfg.Probe.TimeoutMS)*time.Millisecond,
				cfg.Probe.Concurrency)
			fmt.Printf("alive %d / %d after TCP probe\n", len(alive), len(fetched))

			if cfg.Probe.TLSVerify {
				before := len(alive)
				alive = probe.TLS(ctx, alive,
					time.Duration(cfg.Probe.TimeoutMS)*time.Millisecond,
					cfg.Probe.Concurrency)
				fmt.Printf("alive %d / %d after TLS handshake\n", len(alive), before)
			}

			// HTTP-over-proxy verification — this is where we move from
			// "TCP/TLS handshake succeeded" to "traffic actually flows
			// through the proxy". Sort by RTT first so the pool picks
			// the fastest TLS-alive nodes.
			verified := alive
			if cfg.Verify.Enabled {
				sortByLatency(alive)
				before := len(alive)
				verified = verify.Run(ctx, alive, verify.Config{
					Enabled:        cfg.Verify.Enabled,
					CandidatePool:  cfg.Verify.CandidatePool,
					BatchSize:      cfg.Verify.BatchSize,
					BasePort:       cfg.Verify.BasePort,
					Concurrency:    cfg.Verify.Concurrency,
					TimeoutMS:      cfg.Verify.TimeoutMS,
					Rounds:         cfg.Verify.Rounds,
					RoundGapMS:     cfg.Verify.RoundGapMS,
					Targets:        cfg.Verify.Targets,
					SingBoxBin:     cfg.Verify.SingBoxBin,
					StartupTimeout: time.Duration(cfg.Verify.StartupTimeoutMS) * time.Millisecond,
				})
				fmt.Printf("verified %d / %d after HTTP-over-proxy probe\n", len(verified), before)
			}

			enrichGeoIP(cfg, verified)

			selected, summary := aggregate.Run(verified, cfg.Aggregate)
			summary.TotalAlive = len(alive)
			summary.TotalVerified = len(verified)
			summary.TotalFetched = len(fetched)
			summary.GeneratedAtUnix = time.Now().Unix()
			summary.ByCountry = countByCountry(selected)
			fmt.Printf("selected %d nodes\n", len(selected))

			if err := writeOutputs(cfg, selected, summary); err != nil {
				die(err)
			}
			fmt.Println("outputs written to", cfg.Output.Dir)
		},
	}
}

// fetchAll fans out one goroutine per enabled source. Each fetch honors the
// passed context so an upstream deadline (or SIGINT) aborts in-flight HTTP.
func fetchAll(ctx context.Context, cfg *config.Config) []*node.Node {
	timeout := time.Duration(cfg.Probe.TimeoutMS*4) * time.Millisecond
	if timeout < 10*time.Second {
		timeout = 10 * time.Second
	}

	var (
		wg  sync.WaitGroup
		mu  sync.Mutex
		all []*node.Node
	)
	for _, s := range cfg.Sources {
		if !s.Enabled {
			continue
		}
		wg.Add(1)
		go func(src config.Source) {
			defer wg.Done()
			nodes, err := sources.Fetch(ctx, src, timeout)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  [skip] %s: %v\n", src.Name, err)
				return
			}
			if cfg.Probe.MaxNodesPerSource > 0 && len(nodes) > cfg.Probe.MaxNodesPerSource {
				nodes = nodes[:cfg.Probe.MaxNodesPerSource]
			}
			fmt.Fprintf(os.Stderr, "  [ok]   %s: %d nodes\n", src.Name, len(nodes))
			mu.Lock()
			all = append(all, nodes...)
			mu.Unlock()
		}(s)
	}
	wg.Wait()
	return all
}

// enrichGeoIP populates n.Country on every node. Soft-failures on DB/open so
// the pipeline still produces global outputs even without country tags.
func enrichGeoIP(cfg *config.Config, nodes []*node.Node) {
	if !cfg.GeoIP.Enabled {
		return
	}
	if err := geoip.EnsureDB(cfg.GeoIP.DBURL, cfg.GeoIP.DBPath); err != nil {
		fmt.Fprintf(os.Stderr, "  [warn] geoip db: %v\n", err)
		return
	}
	r, err := geoip.Open(cfg.GeoIP.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [warn] geoip open: %v\n", err)
		return
	}
	defer r.Close()
	r.Enrich(nodes, 50)
	fmt.Fprintf(os.Stderr, "  [ok]   geoip enriched %d nodes\n", len(nodes))
}

// sortByLatency orders nodes ascending by LatencyMS in-place. Used before
// verify so the candidate pool is the fastest TLS-alive subset — faster
// nodes are both more likely to pass HTTP verification and more valuable
// to publish.
func sortByLatency(nodes []*node.Node) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].LatencyMS < nodes[j].LatencyMS
	})
}

func countEnabled(srcs []config.Source) int {
	n := 0
	for _, s := range srcs {
		if s.Enabled {
			n++
		}
	}
	return n
}

func countByCountry(ns []*node.Node) map[string]int {
	m := map[string]int{}
	for _, n := range ns {
		cc := n.Country
		if cc == "" {
			cc = "XX"
		}
		m[cc]++
	}
	return m
}

func writeOutputs(cfg *config.Config, selected []*node.Node, summary aggregate.Summary) error {
	outDir := cfg.Output.Dir
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	if err := emitSet(cfg, outDir, "", selected); err != nil {
		return err
	}

	// Per-country outputs: one set per country with ≥ MinPerCountry nodes.
	if cfg.GeoIP.Enabled && cfg.GeoIP.MinPerCountry > 0 {
		byCC := groupByCountry(selected)
		countryDir := filepath.Join(outDir, "by-country")
		// Wipe stale per-country files so dropped countries don't linger.
		_ = os.RemoveAll(countryDir)
		if err := os.MkdirAll(countryDir, 0o755); err != nil {
			return err
		}
		for _, cc := range sortedCountries(byCC) {
			nodes := byCC[cc]
			if len(nodes) < cfg.GeoIP.MinPerCountry {
				continue
			}
			if err := emitSet(cfg, countryDir, cc, nodes); err != nil {
				return fmt.Errorf("emit %s: %w", cc, err)
			}
		}
	}

	statusJSON, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("status marshal: %w", err)
	}
	if err := write(filepath.Join(outDir, "status.json"), string(statusJSON)); err != nil {
		return err
	}

	input := readme.Input{
		Title:          cfg.Readme.Title,
		RepoURL:        cfg.Readme.RepoURL,
		Nodes:          selected,
		Summary:        summary,
		MinPerCountry:  cfg.GeoIP.MinPerCountry,
		CountryEnabled: cfg.GeoIP.Enabled,
	}
	for _, loc := range readme.Locales() {
		if err := write(loc.FileName, readme.Generate(input, loc)); err != nil {
			return fmt.Errorf("readme %s: %w", loc.FileName, err)
		}
	}

	if cfg.Output.Pages.Enabled {
		if err := pages.Generate(pages.Input{
			Title:         cfg.Readme.Title,
			RepoURL:       cfg.Readme.RepoURL,
			SiteURL:       cfg.Output.Pages.SiteURL,
			Summary:       summary,
			Selected:      selected,
			MinPerCountry: cfg.GeoIP.MinPerCountry,
		}, cfg.Output.Pages.Dir); err != nil {
			return fmt.Errorf("pages: %w", err)
		}
	}

	return nil
}

// emitSet writes clash / singbox / v2ray-base64 files for the given node list.
// When suffix is empty, files are named clash.yaml / singbox.json / v2ray-base64.txt.
// When suffix is non-empty (e.g. "HK"), files are named clash-HK.yaml etc.
func emitSet(cfg *config.Config, dir, suffix string, nodes []*node.Node) error {
	formats := map[string]bool{}
	for _, f := range cfg.Output.Formats {
		formats[f] = true
	}
	tag := ""
	if suffix != "" {
		tag = "-" + suffix
	}

	if formats["clash"] {
		content, err := subscribe.Clash(nodes)
		if err != nil {
			return fmt.Errorf("clash: %w", err)
		}
		if err := write(filepath.Join(dir, "clash"+tag+".yaml"), content); err != nil {
			return err
		}
	}
	if formats["singbox"] {
		content, err := subscribe.Singbox(nodes)
		if err != nil {
			return fmt.Errorf("singbox: %w", err)
		}
		if err := write(filepath.Join(dir, "singbox"+tag+".json"), content); err != nil {
			return err
		}
	}
	if formats["v2ray-base64"] {
		if err := write(filepath.Join(dir, "v2ray-base64"+tag+".txt"),
			subscribe.V2RayBase64(nodes)); err != nil {
			return err
		}
	}
	return nil
}

func groupByCountry(ns []*node.Node) map[string][]*node.Node {
	m := map[string][]*node.Node{}
	for _, n := range ns {
		cc := n.Country
		if cc == "" {
			continue // skip unknown — no point publishing an "XX" bucket
		}
		m[cc] = append(m[cc], n)
	}
	return m
}

func sortedCountries(m map[string][]*node.Node) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func write(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
