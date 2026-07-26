package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/starhistory"
)

func newStarHistoryCmd() *cobra.Command {
	return newStarHistoryCmdWithDeps(http.DefaultClient, "https://api.github.com", os.Getenv)
}

func newStarHistoryCmdWithDeps(client *http.Client, apiBase string, getenv func(string) string) *cobra.Command {
	var repo string
	var output string
	command := &cobra.Command{
		Use:   "star-history",
		Short: "Generate a self-hosted GitHub star history SVG",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if strings.TrimSpace(repo) == "" {
				return fmt.Errorf("--repo is required")
			}
			if strings.TrimSpace(output) == "" {
				return fmt.Errorf("--output is required")
			}
			token := strings.TrimSpace(getenv("GITHUB_TOKEN"))
			if token == "" {
				return fmt.Errorf("GITHUB_TOKEN is required")
			}
			stars, err := starhistory.Fetch(cmd.Context(), client, apiBase, repo, token)
			if err != nil {
				return err
			}
			body, err := starhistory.RenderSVG(repo, stars)
			if err != nil {
				return err
			}
			if err := writeAtomic(output, body); err != nil {
				return fmt.Errorf("write star history: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "stars=%d output=%s\n", len(stars), output)
			return nil
		},
	}
	command.Flags().StringVar(&repo, "repo", "", "GitHub repository in owner/repo form")
	command.Flags().StringVar(&output, "output", "assets/star-history.svg", "SVG output path")
	return command
}

func writeAtomic(path string, body []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	file, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	temp := file.Name()
	defer os.Remove(temp)
	if err := file.Chmod(0o644); err != nil {
		file.Close()
		return err
	}
	if _, err := file.Write(body); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(temp, path)
}
