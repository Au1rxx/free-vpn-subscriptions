package store

import (
	"bytes"
	"crypto/sha256"
	"strings"
	"testing"
	"testing/fstest"

	dbmigrations "github.com/Au1rxx/free-vpn-subscriptions/db/migrations"
)

func TestEmbeddedMigrationFilesParse(t *testing.T) {
	got, err := loadMigrations(dbmigrations.Files)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 11 || got[0].Version != "0001" || got[10].Version != "0011" {
		t.Fatalf("embedded migrations=%v", got)
	}
}

func TestSplitStatements(t *testing.T) {
	source := "-- fnctl:statement\nCREATE TABLE a(id INT);\n-- fnctl:statement\nCREATE TABLE b(id INT);\n"
	got, err := splitStatements(source)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"CREATE TABLE a(id INT);", "CREATE TABLE b(id INT);"}
	if len(got) != len(want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("statement %d=%q want %q", i, got[i], want[i])
		}
	}
}

func TestSplitStatementsRequiresDelimiter(t *testing.T) {
	if _, err := splitStatements("CREATE TABLE a(id INT);"); err == nil {
		t.Fatal("expected missing delimiter error")
	}
}

func TestLoadMigrationsSortsAndChecksums(t *testing.T) {
	files := fstest.MapFS{
		"0002_second.sql": &fstest.MapFile{Data: []byte("-- fnctl:statement\nSELECT 2;")},
		"0001_first.sql":  &fstest.MapFile{Data: []byte("-- fnctl:statement\nSELECT 1;")},
		"README.md":       &fstest.MapFile{Data: []byte("ignored")},
	}
	got, err := loadMigrations(files)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Version != "0001" || got[1].Version != "0002" {
		t.Fatalf("unexpected order: %#v", got)
	}
	want := sha256.Sum256(files["0001_first.sql"].Data)
	if got[0].Checksum != want {
		t.Fatal("checksum mismatch")
	}
}

func TestVerifyMigrationChecksumMismatchStops(t *testing.T) {
	migration := Migration{Version: "0002", Name: "0002_sources.sql", Checksum: sha256.Sum256([]byte("changed"))}
	applied := AppliedMigration{Version: "0002", Checksum: sha256.Sum256([]byte("original"))}
	err := verifyMigration(migration, applied)
	if err == nil || !strings.Contains(err.Error(), "migration checksum mismatch") {
		t.Fatalf("err=%v", err)
	}
	if bytes.Equal(migration.Checksum[:], applied.Checksum[:]) {
		t.Fatal("fixture checksums unexpectedly equal")
	}
}

func TestBootstrapMigrationSeparatesAdminAndDatabaseStatements(t *testing.T) {
	migration := Migration{Version: "0001", Statements: []string{
		"CREATE DATABASE vpn_nodes;",
		"CREATE TABLE schema_migrations(version VARCHAR(32));",
	}}
	admin, database, err := bootstrapStatements(migration)
	if err != nil {
		t.Fatal(err)
	}
	if len(admin) != 1 || !strings.HasPrefix(admin[0], "CREATE DATABASE") {
		t.Fatalf("admin statements=%#v", admin)
	}
	if len(database) != 1 || !strings.HasPrefix(database[0], "CREATE TABLE") {
		t.Fatalf("database statements=%#v", database)
	}
}
