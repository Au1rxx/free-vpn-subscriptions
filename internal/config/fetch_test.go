package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const fetchTestSources = `sources:
  - name: fixture
    url: https://example.test/sub
    format: uri-list
    enabled: true
`

func writeFetchConfig(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(fetchTestSources+body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestFetchCeilingDefaults(t *testing.T) {
	cfg, err := Load(writeFetchConfig(t, ""))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Fetch.MaxBodyBytes != 192<<20 {
		t.Fatalf("max_body_bytes=%d", cfg.Fetch.MaxBodyBytes)
	}
	if cfg.Fetch.MaxDecodedBytes != 384<<20 {
		t.Fatalf("max_decoded_bytes=%d", cfg.Fetch.MaxDecodedBytes)
	}
}

func TestFetchCeilingOverride(t *testing.T) {
	cfg, err := Load(writeFetchConfig(t, "fetch:\n  max_body_bytes: 1048576\n  max_decoded_bytes: 2097152\n"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Fetch.MaxBodyBytes != 1<<20 || cfg.Fetch.MaxDecodedBytes != 2<<20 {
		t.Fatalf("ceilings=%d/%d", cfg.Fetch.MaxBodyBytes, cfg.Fetch.MaxDecodedBytes)
	}
}

func TestFetchCeilingRejectsInvalid(t *testing.T) {
	for name, body := range map[string]string{
		"negative_body": "fetch:\n  max_body_bytes: -1\n",
		"decoded_below": "fetch:\n  max_body_bytes: 2097152\n  max_decoded_bytes: 1048576\n",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := Load(writeFetchConfig(t, body))
			if err == nil || !strings.Contains(err.Error(), "fetch ") {
				t.Fatalf("err=%v", err)
			}
		})
	}
}
