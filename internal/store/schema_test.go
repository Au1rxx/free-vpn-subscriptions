package store

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"

	dbmigrations "github.com/Au1rxx/free-vpn-subscriptions/db/migrations"
)

var requiredTables = []string{
	"schema_migrations", "sources", "source_links", "source_fetches", "raw_payloads",
	"parse_runs", "parse_errors", "endpoints", "node_configs", "node_source_stats",
	"node_source_daily", "validation_queue", "validation_batches", "validation_attempts",
	"node_current_status", "node_daily_stats", "node_classifications", "export_runs",
	"export_members", "storage_metrics", "storage_policies", "node_tombstones",
}

func TestMigrationFilesDeclareAllTablesAndChineseComments(t *testing.T) {
	ddl := readMigrationDDL(t)
	for _, table := range requiredTables {
		marker := "CREATE TABLE IF NOT EXISTS `" + table + "`"
		if !strings.Contains(ddl, marker) {
			t.Errorf("missing table %s", table)
		}
	}
	if got := strings.Count(ddl, "CREATE TABLE IF NOT EXISTS `"); got != len(requiredTables) {
		t.Errorf("table declarations=%d want=%d", got, len(requiredTables))
	}
	for _, block := range tableBlocks(ddl) {
		if !regexp.MustCompile(`(?s)COMMENT='[^']*[\x{4e00}-\x{9fff}][^']*'`).MatchString(block) {
			t.Errorf("table block lacks Chinese comment: %.80s", block)
		}
		for _, line := range strings.Split(block, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "`") && !strings.Contains(trimmed, " COMMENT '") {
				t.Errorf("column lacks comment: %s", trimmed)
			}
		}
	}
}

func TestMigrationFilesContainCriticalUniqueKeys(t *testing.T) {
	ddl := readMigrationDDL(t)
	keys := []string{
		"uk_sources_canonical_url", "uk_raw_payloads_sha256", "uk_endpoints_host_port",
		"uk_node_configs_fingerprint", "uk_node_source_stats_pair", "uk_validation_queue_node_stage",
		"uk_export_members_pair", "uk_node_tombstones_fingerprint",
	}
	for _, key := range keys {
		if !strings.Contains(ddl, "`"+key+"`") {
			t.Errorf("missing unique key %s", key)
		}
	}
}

func readMigrationDDL(t *testing.T) string {
	t.Helper()
	entries, err := fs.ReadDir(dbmigrations.Files, ".")
	if err != nil {
		t.Fatal(err)
	}
	var all strings.Builder
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		body, err := fs.ReadFile(dbmigrations.Files, entry.Name())
		if err != nil {
			t.Fatal(err)
		}
		all.Write(body)
		all.WriteByte('\n')
	}
	return all.String()
}

func tableBlocks(ddl string) []string {
	parts := strings.Split(ddl, "CREATE TABLE IF NOT EXISTS `")
	if len(parts) < 2 {
		return nil
	}
	blocks := make([]string, 0, len(parts)-1)
	for _, part := range parts[1:] {
		if end := strings.Index(part, ";"); end >= 0 {
			part = part[:end+1]
		}
		blocks = append(blocks, part)
	}
	return blocks
}
