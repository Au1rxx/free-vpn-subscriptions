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
