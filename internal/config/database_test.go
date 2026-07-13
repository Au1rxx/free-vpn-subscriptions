package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDatabaseDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	data := []byte(`sources:
  - name: fixture
    url: https://example.test/sub
    format: uri-list
    enabled: true
`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Database.Address != "127.0.0.1:13306" {
		t.Fatalf("address=%q", cfg.Database.Address)
	}
	if cfg.Database.Name != "vpn_nodes" {
		t.Fatalf("name=%q", cfg.Database.Name)
	}
	if cfg.Database.TLSMode != "required" {
		t.Fatalf("tls_mode=%q", cfg.Database.TLSMode)
	}
	if cfg.Database.MaxOpenConns != 20 || cfg.Database.MaxIdleConns != 10 {
		t.Fatalf("pool=%d/%d", cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns)
	}
}

func TestDatabaseRejectsUnknownTLSMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	data := []byte(`sources:
  - name: fixture
    url: https://example.test/sub
    format: uri-list
    enabled: true
database:
  tls_mode: plaintext
`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected invalid tls_mode error")
	}
}

func TestLoadExpandsSystemdCredentialDirectory(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CREDENTIALS_DIRECTORY", filepath.Join(dir, "credentials"))
	path := filepath.Join(dir, "config.yaml")
	data := []byte(`sources:
  - name: fixture
    url: https://example.test/sub
    format: auto
database:
  enabled: true
  password_file: "%d/mysql-password"
`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "credentials", "mysql-password")
	if cfg.Database.PasswordFile != want {
		t.Fatalf("password_file=%q want %q", cfg.Database.PasswordFile, want)
	}
}

func TestLoadDetailedSourceMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	data := []byte(`sources:
  - name: catalog
    url: https://example.test/sub
    format: auto
    kind: github-wiki
    discovery_method: researched-seed
    depth: 1
    priority: 90
    fetch_interval_seconds: 1800
    enabled: true
`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	source := cfg.Sources[0]
	if source.Kind != "github-wiki" || source.DiscoveryMethod != "researched-seed" || source.Depth != 1 || source.Priority != 90 || source.FetchIntervalSeconds != 1800 {
		t.Fatalf("source metadata=%+v", source)
	}
}
