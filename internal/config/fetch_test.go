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
	if cfg.Fetch.ArchiveThresholdBytes != 50<<20 {
		t.Fatalf("archive_threshold_bytes=%d", cfg.Fetch.ArchiveThresholdBytes)
	}
	if cfg.Fetch.ArchiveWriteEnabled {
		t.Fatal("archive_write_enabled defaulted to true")
	}
	if cfg.Fetch.ArchiveDirectory != filepath.Join("output", ".raw-archive") {
		t.Fatalf("archive_directory=%q", cfg.Fetch.ArchiveDirectory)
	}
}

func TestFetchCeilingOverride(t *testing.T) {
	cfg, err := Load(writeFetchConfig(t, "fetch:\n  max_body_bytes: 1048576\n  max_decoded_bytes: 2097152\n  archive_directory: /srv/raw-archive\n  archive_threshold_bytes: 524288\n  archive_write_enabled: true\n"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Fetch.MaxBodyBytes != 1<<20 || cfg.Fetch.MaxDecodedBytes != 2<<20 {
		t.Fatalf("ceilings=%d/%d", cfg.Fetch.MaxBodyBytes, cfg.Fetch.MaxDecodedBytes)
	}
	if cfg.Fetch.ArchiveDirectory != "/srv/raw-archive" || cfg.Fetch.ArchiveThresholdBytes != 512<<10 || !cfg.Fetch.ArchiveWriteEnabled {
		t.Fatalf("archive config=%q/%d/%t", cfg.Fetch.ArchiveDirectory, cfg.Fetch.ArchiveThresholdBytes, cfg.Fetch.ArchiveWriteEnabled)
	}
}

func TestFetchCeilingRejectsInvalid(t *testing.T) {
	for name, body := range map[string]string{
		"negative_body":                     "fetch:\n  max_body_bytes: -1\n",
		"decoded_below":                     "fetch:\n  max_body_bytes: 2097152\n  max_decoded_bytes: 1048576\n",
		"negative_archive_threshold":        "fetch:\n  archive_threshold_bytes: -1\n",
		"enabled_archive_without_directory": "fetch:\n  archive_write_enabled: true\n  archive_directory: '   '\n",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := Load(writeFetchConfig(t, body))
			if err == nil || !strings.Contains(err.Error(), "fetch ") {
				t.Fatalf("err=%v", err)
			}
		})
	}
}
