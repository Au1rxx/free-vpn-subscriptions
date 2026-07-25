package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/exportdb"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/pages"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/readme"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

func newExportDBCmd() *cobra.Command {
	var output string
	var siteRoot string
	var shardSize int
	command := &cobra.Command{
		Use:   "export-db",
		Short: "Export verified classified subscriptions from MySQL",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := loadDatabaseConfig()
			if err != nil {
				return err
			}
			db, err := store.Open(cmd.Context(), cfg.Database, cfg.Database.Name)
			if err != nil {
				return err
			}
			defer db.Close()
			if _, err := store.CheckServer(cmd.Context(), db); err != nil {
				return err
			}
			if output == "" {
				output = cfg.Output.Dir
			}
			report, err := (exportdb.Service{DB: db, Output: output, ShardSize: shardSize}).Run(cmd.Context())
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "candidates=%d stable=%d collections=%d files=%d bytes=%d output=%s\n",
				report.Candidates, report.Stable, len(report.Collections), report.Files, report.Bytes, output)
			if siteRoot == "" {
				return nil
			}
			locales, err := renderSite(cfg, siteRoot, report)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "site_root=%s readme_locales=%d pages=%t nodes=%d\n",
				siteRoot, locales, cfg.Output.Pages.Enabled, len(report.Selected))
			return nil
		},
	}
	command.Flags().StringVar(&output, "output", "", "output directory (defaults to config output.dir)")
	command.Flags().StringVar(&siteRoot, "site-root", "", "public repository root to regenerate README locales and the Pages site into")
	command.Flags().IntVar(&shardSize, "shard-size", exportdb.DefaultShardSize, "maximum nodes per output file (1-2000)")
	return command
}

// renderSite regenerates the public README locales and the Pages site from a
// database export. The legacy aggregate command does this inline; the database
// publisher needs it too, otherwise the site keeps serving the last legacy run.
func renderSite(cfg *config.Config, root string, report exportdb.Report) (int, error) {
	info, err := os.Stat(root)
	if err != nil {
		return 0, fmt.Errorf("site root: %w", err)
	}
	if !info.IsDir() {
		return 0, fmt.Errorf("site root %q is not a directory", root)
	}
	input := readme.Input{
		Title:          cfg.Readme.Title,
		RepoURL:        cfg.Readme.RepoURL,
		Nodes:          report.Selected,
		Summary:        report.Summary,
		MinPerCountry:  cfg.GeoIP.MinPerCountry,
		CountryEnabled: cfg.GeoIP.Enabled,
	}
	locales := readme.Locales()
	for _, loc := range locales {
		if err := os.WriteFile(filepath.Join(root, loc.FileName), []byte(readme.Generate(input, loc)), 0o644); err != nil {
			return 0, fmt.Errorf("readme %s: %w", loc.FileName, err)
		}
	}
	if !cfg.Output.Pages.Enabled {
		return len(locales), nil
	}
	directory := cfg.Output.Pages.Dir
	if !filepath.IsAbs(directory) {
		directory = filepath.Join(root, directory)
	}
	if err := pages.Generate(pages.Input{
		Title:         cfg.Readme.Title,
		RepoURL:       cfg.Readme.RepoURL,
		SiteURL:       cfg.Output.Pages.SiteURL,
		Summary:       report.Summary,
		Selected:      report.Selected,
		MinPerCountry: cfg.GeoIP.MinPerCountry,
	}, directory); err != nil {
		return 0, fmt.Errorf("pages: %w", err)
	}
	return len(locales), nil
}
