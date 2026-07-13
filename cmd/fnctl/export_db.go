package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/exportdb"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

func newExportDBCmd() *cobra.Command {
	var output string
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
			return nil
		},
	}
	command.Flags().StringVar(&output, "output", "", "output directory (defaults to config output.dir)")
	command.Flags().IntVar(&shardSize, "shard-size", exportdb.DefaultShardSize, "maximum nodes per output file (1-2000)")
	return command
}
