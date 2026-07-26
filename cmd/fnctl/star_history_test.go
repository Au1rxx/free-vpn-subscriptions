package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStarHistoryCommandRequiresToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	command := newStarHistoryCmd()
	command.SetArgs([]string{"--repo", "acme/tool", "--output", filepath.Join(t.TempDir(), "stars.svg")})
	err := command.Execute()
	if err == nil || !strings.Contains(err.Error(), "GITHUB_TOKEN") {
		t.Fatalf("err=%v", err)
	}
}

func TestStarHistoryCommandWritesSVGAtomically(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			http.Error(w, "bad token", http.StatusUnauthorized)
			return
		}
		fmt.Fprint(w, `[{"starred_at":"2026-07-01T10:00:00Z"}]`)
	}))
	defer server.Close()

	output := filepath.Join(t.TempDir(), "nested", "stars.svg")
	command := newStarHistoryCmdWithDeps(server.Client(), server.URL, func(string) string {
		return "test-token"
	})
	command.SetArgs([]string{"--repo", "acme/tool", "--output", output})
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "<svg") {
		t.Fatalf("not SVG: %q", body)
	}
	if matches, err := filepath.Glob(output + ".*.tmp"); err != nil || len(matches) != 0 {
		t.Fatalf("temporary files remain: %v err=%v", matches, err)
	}
}

func TestRootCommandRegistersStarHistory(t *testing.T) {
	command, _, err := newRootCmd().Find([]string{"star-history"})
	if err != nil {
		t.Fatal(err)
	}
	if command.Name() != "star-history" {
		t.Fatalf("name=%q", command.Name())
	}
}
