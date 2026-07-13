package store

import (
	"os"
	"path/filepath"
	"testing"

	appconfig "github.com/Au1rxx/free-vpn-subscriptions/internal/config"
)

func TestReadPasswordTrimsOnlyLineEndings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mysql-password")
	if err := os.WriteFile(path, []byte("p@ss word\r\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := ReadPassword(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != "p@ss word" {
		t.Fatalf("got %q", got)
	}
}

func TestReadPasswordRejectsEmptyCredential(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mysql-password")
	if err := os.WriteFile(path, []byte("\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadPassword(path); err == nil {
		t.Fatal("expected empty credential error")
	}
}

func TestNewMySQLConfigPreservesPasswordAndRequiresTLS(t *testing.T) {
	cfg := appconfig.DatabaseConfig{
		Address: "127.0.0.1:13306",
		User:    "db-user",
		TLSMode: "required",
	}
	password := "symbols:?@#%="
	got := NewMySQLConfig(cfg, password, "vpn_nodes")
	if got.Passwd != password {
		t.Fatal("password was changed")
	}
	if got.DBName != "vpn_nodes" || got.Net != "tcp" || got.Addr != cfg.Address {
		t.Fatalf("unexpected mysql config: net=%q addr=%q db=%q", got.Net, got.Addr, got.DBName)
	}
	if got.TLSConfig != "skip-verify" || !got.ParseTime || got.Loc.String() != "UTC" {
		t.Fatalf("tls/time config=%q/%v/%v", got.TLSConfig, got.ParseTime, got.Loc)
	}
}
