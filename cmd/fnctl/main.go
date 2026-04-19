// Package main provides the fnctl CLI — entry point for the aggregator.
//
// The only relevant command is `aggregate`: it fetches every enabled source,
// probes every node for TCP reachability, deduplicates and ranks the alive
// set, and writes output files (clash.yaml, singbox.json, v2ray-base64.txt,
// status.json) plus a generated README.md.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/aggregate"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/node"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/probe"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/readme"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/sources"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/subscribe"
)

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

			fetched := fetchAll(cfg)
			fmt.Printf("fetched %d nodes from %d sources\n", len(fetched), countEnabled(cfg.Sources))

			alive := probe.TCP(fetched,
				time.Duration(cfg.Probe.TimeoutMS)*time.Millisecond,
				cfg.Probe.Concurrency)
			fmt.Printf("alive %d / %d after TCP probe\n", len(alive), len(fetched))

			selected, summary := aggregate.Run(alive, cfg.Aggregate)
			summary.TotalFetched = len(fetched)
			summary.GeneratedAtUnix = time.Now().Unix()
			fmt.Printf("selected %d nodes\n", len(selected))

			if err := writeOutputs(cfg, selected, summary); err != nil {
				die(err)
			}
			fmt.Println("outputs written to", cfg.Output.Dir)
		},
	}
}

// fetchAll fans out one goroutine per enabled source.
func fetchAll(cfg *config.Config) []*node.Node {
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
			nodes, err := sources.Fetch(src, timeout)
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

func countEnabled(srcs []config.Source) int {
	n := 0
	for _, s := range srcs {
		if s.Enabled {
			n++
		}
	}
	return n
}

func writeOutputs(cfg *config.Config, selected []*node.Node, summary aggregate.Summary) error {
	outDir := cfg.Output.Dir
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	formats := map[string]bool{}
	for _, f := range cfg.Output.Formats {
		formats[f] = true
	}

	if formats["clash"] {
		content, err := subscribe.Clash(selected)
		if err != nil {
			return fmt.Errorf("clash: %w", err)
		}
		if err := write(filepath.Join(outDir, "clash.yaml"), content); err != nil {
			return err
		}
	}
	if formats["singbox"] {
		content, err := subscribe.Singbox(selected)
		if err != nil {
			return fmt.Errorf("singbox: %w", err)
		}
		if err := write(filepath.Join(outDir, "singbox.json"), content); err != nil {
			return err
		}
	}
	if formats["v2ray-base64"] {
		if err := write(filepath.Join(outDir, "v2ray-base64.txt"),
			subscribe.V2RayBase64(selected)); err != nil {
			return err
		}
	}

	statusJSON, _ := json.MarshalIndent(summary, "", "  ")
	if err := write(filepath.Join(outDir, "status.json"), string(statusJSON)); err != nil {
		return err
	}

	md := readme.Generate(readme.Input{
		Title:   cfg.Readme.Title,
		RepoURL: cfg.Readme.RepoURL,
		Nodes:   selected,
		Summary: summary,
	})
	return write("README.md", md)
}

func write(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
