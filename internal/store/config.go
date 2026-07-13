// Package store provides MySQL persistence for harvested proxy data.
package store

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	appconfig "github.com/Au1rxx/free-vpn-subscriptions/internal/config"
)

// ReadPassword reads a systemd credential or other mode-0600 password file.
// Only trailing line endings are removed so spaces and punctuation remain
// valid password characters.
func ReadPassword(path string) (string, error) {
	if path == "" {
		return "", errors.New("database credential path is empty")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read database credential: %w", err)
	}
	password := strings.TrimRight(string(b), "\r\n")
	if password == "" {
		return "", errors.New("database credential is empty")
	}
	return password, nil
}

// NewMySQLConfig constructs a driver configuration without interpolating
// credentials into logs or hand-built DSN strings.
func NewMySQLConfig(cfg appconfig.DatabaseConfig, password, database string) *mysql.Config {
	tlsMode := cfg.TLSMode
	if tlsMode == "required" {
		// The SSH tunnel and TLS encryption are mandatory. Identity verification
		// is enabled later when the Oracle CA bundle is installed.
		tlsMode = "skip-verify"
	}
	return &mysql.Config{
		User:              cfg.User,
		Passwd:            password,
		Net:               "tcp",
		Addr:              cfg.Address,
		DBName:            database,
		Collation:         "utf8mb4_0900_ai_ci",
		Loc:               time.UTC,
		ParseTime:         true,
		TLSConfig:         tlsMode,
		Timeout:           10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		RejectReadOnly:    true,
		CheckConnLiveness: true,
	}
}
