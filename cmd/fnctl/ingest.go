package main

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/discovery"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/ingest"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

func newImportSeedsCmd() *cobra.Command {
	return &cobra.Command{Use: "import-seeds", Short: "Import configured sources into MySQL", RunE: func(cmd *cobra.Command, _ []string) error {
		cfg, db, service, err := openIngestService(cmd.Context())
		if err != nil {
			return err
		}
		defer db.Close()
		summary, err := service.ImportSeeds(cmd.Context(), cfg.Sources)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "sources=%d upserted=%d\n", summary.Sources, summary.InsertedOrUpdated)
		return nil
	}}
}

func newFetchCmd() *cobra.Command {
	limit := 100
	command := &cobra.Command{Use: "fetch", Short: "Fetch one bounded batch into MySQL", RunE: func(cmd *cobra.Command, _ []string) error {
		_, db, service, err := openIngestService(cmd.Context())
		if err != nil {
			return err
		}
		defer db.Close()
		summary, err := service.Fetch(cmd.Context(), limit)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "sources=%d success=%d not_modified=%d failed=%d spooled=%d bytes=%d\n",
			summary.Sources, summary.Success, summary.NotModified, summary.Failed, summary.Spooled, summary.Bytes)
		return nil
	}}
	command.Flags().IntVar(&limit, "limit", 100, "maximum due sources to fetch (1-1000)")
	return command
}

func newParseCmd() *cobra.Command {
	limit := 100
	command := &cobra.Command{Use: "parse", Short: "Parse one bounded pending fetch batch", RunE: func(cmd *cobra.Command, _ []string) error {
		_, db, service, err := openIngestService(cmd.Context())
		if err != nil {
			return err
		}
		defer db.Close()
		summary, err := service.Parse(cmd.Context(), limit)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "fetches=%d nodes=%d errors=%d new_endpoints=%d new_configs=%d queue_jobs=%d discovered=%d\n",
			summary.Fetches, summary.Nodes, summary.Errors, summary.NewEndpoints, summary.NewConfigs, summary.QueueJobs, summary.Discovered)
		return nil
	}}
	command.Flags().IntVar(&limit, "limit", 100, "maximum pending fetches to parse (1-1000)")
	return command
}

func newIngestStatusCmd() *cobra.Command {
	return &cobra.Command{Use: "ingest-status", Short: "Show bounded ingestion counters", RunE: func(cmd *cobra.Command, _ []string) error {
		_, db, _, err := openIngestService(cmd.Context())
		if err != nil {
			return err
		}
		defer db.Close()
		status, err := store.ReadIngestStatus(cmd.Context(), db)
		if err != nil {
			return err
		}
		fmt.Fprint(cmd.OutOrStdout(), formatIngestStatus(status))
		return nil
	}}
}

func formatIngestStatus(status store.IngestStatus) string {
	var output strings.Builder
	fmt.Fprintf(&output, "sources=%d fetches=%d pending_fetches=%d parse_runs=%d\n", status.Sources, status.Fetches, status.PendingFetches, status.ParseRuns)
	fmt.Fprintf(&output, "endpoints=%d configs=%d parse_errors=%d queue_pending=%d\n", status.Endpoints, status.Configs, status.ParseErrors, status.QueuePending)
	protocols := make([]string, 0, len(status.ByProtocol))
	for protocol := range status.ByProtocol {
		protocols = append(protocols, protocol)
	}
	sort.Strings(protocols)
	for _, protocol := range protocols {
		fmt.Fprintf(&output, "protocol=%s configs=%d\n", protocol, status.ByProtocol[protocol])
	}
	return output.String()
}

func newDiscoverCmd() *cobra.Command {
	var kind, sourceURL, query, tokenFile string
	limit := 500
	command := &cobra.Command{Use: "discover", Short: "Discover and register public source candidates", RunE: func(cmd *cobra.Command, _ []string) error {
		adapter, err := selectDiscoverer(kind)
		if err != nil {
			return err
		}
		candidates, err := adapter.Discover(cmd.Context(), discovery.Seed{URL: sourceURL, Query: query, TokenFile: tokenFile})
		if err != nil {
			return err
		}
		if len(candidates) > limit {
			candidates = candidates[:limit]
		}
		_, db, _, err := openIngestService(cmd.Context())
		if err != nil {
			return err
		}
		defer db.Close()
		upserted := 0
		for _, candidate := range candidates {
			if !discovery.LikelySubscriptionURL(candidate.URL) {
				continue
			}
			if _, err := store.UpsertSource(cmd.Context(), db, store.SourceRecord{
				Name: candidate.Name, URL: candidate.URL, Kind: candidate.Kind, FormatHint: "auto",
				DiscoveryMethod: kind, State: "active", Enabled: true, Depth: candidate.Depth,
			}); err == nil {
				upserted++
			}
		}
		fmt.Fprintf(cmd.OutOrStdout(), "kind=%s candidates=%d upserted=%d\n", kind, len(candidates), upserted)
		return nil
	}}
	command.Flags().StringVar(&kind, "kind", "web", "github|gitlab|gitea|gitee|telegram|web|api")
	command.Flags().StringVar(&sourceURL, "url", "", "public page or API seed URL")
	command.Flags().StringVar(&query, "query", "", "code search query")
	command.Flags().StringVar(&tokenFile, "token-file", "", "optional credential file")
	command.Flags().IntVar(&limit, "limit", 500, "maximum candidates to register")
	return command
}

func newPruneDiscoveryCmd() *cobra.Command {
	var kind, method string
	command := &cobra.Command{Use: "prune-discovery", Short: "Disable one exact class of automatically discovered sources", RunE: func(cmd *cobra.Command, _ []string) error {
		if kind == "" || method == "" {
			return fmt.Errorf("kind and method are required")
		}
		_, db, _, err := openIngestService(cmd.Context())
		if err != nil {
			return err
		}
		defer db.Close()
		affected, err := store.DisableDiscoveredSources(cmd.Context(), db, kind, method)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "kind=%s method=%s disabled=%d\n", kind, method, affected)
		return nil
	}}
	command.Flags().StringVar(&kind, "kind", "", "exact source kind")
	command.Flags().StringVar(&method, "method", "", "exact discovery method")
	return command
}

func openIngestService(ctx context.Context) (*config.Config, *sql.DB, *ingest.Service, error) {
	cfg, err := loadDatabaseConfig()
	if err != nil {
		return nil, nil, nil, err
	}
	db, err := store.Open(ctx, cfg.Database, cfg.Database.Name)
	if err != nil {
		return nil, nil, nil, err
	}
	if _, err := store.CheckServer(ctx, db); err != nil {
		db.Close()
		return nil, nil, nil, err
	}
	spool, err := ingest.NewSpool(filepath.Join(cfg.Output.Dir, ".spool"), 2<<30)
	if err != nil {
		db.Close()
		return nil, nil, nil, err
	}
	return cfg, db, &ingest.Service{DB: db, Spool: spool}, nil
}

func selectDiscoverer(kind string) (discovery.Discoverer, error) {
	switch kind {
	case "github":
		return discovery.GitHubDiscoverer{}, nil
	case "gitlab":
		return discovery.GitLabDiscoverer{}, nil
	case "gitea", "codeberg":
		return discovery.GiteaDiscoverer{}, nil
	case "gitee":
		return discovery.GiteeDiscoverer{}, nil
	case "telegram":
		return discovery.TelegramDiscoverer{}, nil
	case "web":
		return discovery.WebDiscoverer{}, nil
	case "api":
		return discovery.APIDiscoverer{}, nil
	default:
		return nil, fmt.Errorf("unsupported discovery kind %q", kind)
	}
}
