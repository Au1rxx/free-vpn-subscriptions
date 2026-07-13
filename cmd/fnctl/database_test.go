package main

import (
	"strings"
	"testing"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

func TestRootContainsDatabaseCommands(t *testing.T) {
	root := newRootCmd()
	want := map[string]bool{
		"aggregate": false, "migrate": false, "db-status": false,
		"import-seeds": false, "fetch": false, "parse": false, "discover": false, "ingest-status": false,
	}
	for _, command := range root.Commands() {
		if _, ok := want[command.Name()]; ok {
			want[command.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("missing command %s", name)
		}
	}
}

func TestFormatDatabaseStatusIsBoundedAndContainsNoDSN(t *testing.T) {
	status := store.DatabaseStatus{
		Server: store.ServerInfo{
			Version: "9.7.1-cloud", Cipher: "TLS_AES_128_GCM_SHA256",
			TimeZone: "UTC", Charset: "utf8mb4", Collation: "utf8mb4_0900_ai_ci",
		},
		AppliedMigrations:   6,
		BusinessTables:      22,
		EmptyTableComments:  0,
		EmptyColumnComments: 0,
		EnabledPolicies:     6,
		DataBytes:           1024,
		IndexBytes:          2048,
	}
	out := formatDatabaseStatus(status)
	for _, want := range []string{
		"version=9.7.1-cloud", "tls=TLS_AES_128_GCM_SHA256", "migrations=6", "tables=22",
		"empty_table_comments=0", "empty_column_comments=0", "enabled_policies=6",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("status missing %q: %s", want, out)
		}
	}
	if strings.Contains(out, "@tcp(") || strings.Contains(out, "password") {
		t.Fatalf("status leaked connection details: %s", out)
	}
	if lines := strings.Count(out, "\n"); lines > 10 {
		t.Fatalf("status output too long: %d lines", lines)
	}
}
