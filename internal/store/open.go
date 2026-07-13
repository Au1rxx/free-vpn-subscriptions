package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	appconfig "github.com/Au1rxx/free-vpn-subscriptions/internal/config"
)

// ServerInfo captures the server properties required by the harvester.
type ServerInfo struct {
	Version   string
	ReadOnly  bool
	Cipher    string
	TimeZone  string
	Charset   string
	Collation string
}

// Open creates and verifies a bounded MySQL connection pool. It deliberately
// does not expose the generated DSN in returned errors.
func Open(ctx context.Context, cfg appconfig.DatabaseConfig, database string) (*sql.DB, error) {
	password, err := ReadPassword(cfg.PasswordFile)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("mysql", NewMySQLConfig(cfg, password, database).FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	applyPoolSettings(db, cfg)
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

func applyPoolSettings(db *sql.DB, cfg appconfig.DatabaseConfig) {
	maxOpen := cfg.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 20
	}
	maxIdle := cfg.MaxIdleConns
	if maxIdle < 0 || maxIdle > maxOpen {
		maxIdle = maxOpen / 2
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetConnMaxIdleTime(time.Minute)
}

// CheckServer reads and validates the properties that make a connection safe
// for persistent writes. The TLS cipher is read from this session.
func CheckServer(ctx context.Context, db *sql.DB) (ServerInfo, error) {
	var info ServerInfo
	var readOnly int
	err := db.QueryRowContext(ctx, `
		SELECT VERSION(), @@read_only, @@time_zone,
		       @@character_set_server, @@collation_server`).Scan(
		&info.Version, &readOnly, &info.TimeZone, &info.Charset, &info.Collation)
	if err != nil {
		return ServerInfo{}, fmt.Errorf("query database server properties: %w", err)
	}
	info.ReadOnly = readOnly != 0
	var statusName string
	if err := db.QueryRowContext(ctx, "SHOW SESSION STATUS LIKE 'Ssl_cipher'").Scan(&statusName, &info.Cipher); err != nil {
		return ServerInfo{}, fmt.Errorf("query database TLS status: %w", err)
	}
	if err := validateServerInfo(info); err != nil {
		return ServerInfo{}, err
	}
	return info, nil
}

func validateServerInfo(info ServerInfo) error {
	if info.Cipher == "" {
		return fmt.Errorf("database connection has no TLS cipher")
	}
	if info.ReadOnly {
		return fmt.Errorf("database server is read-only")
	}
	if info.TimeZone != "UTC" && info.TimeZone != "+00:00" {
		return fmt.Errorf("database time zone is %q, expected UTC", info.TimeZone)
	}
	return nil
}
