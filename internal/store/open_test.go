package store

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"

	appconfig "github.com/Au1rxx/free-vpn-subscriptions/internal/config"
)

type nopConnector struct{}

func (nopConnector) Connect(context.Context) (driver.Conn, error) {
	return nil, errors.New("not used")
}

func (nopConnector) Driver() driver.Driver { return nopDriver{} }

type nopDriver struct{}

func (nopDriver) Open(string) (driver.Conn, error) { return nil, errors.New("not used") }

func TestApplyPoolSettings(t *testing.T) {
	db := sql.OpenDB(nopConnector{})
	defer db.Close()
	applyPoolSettings(db, appconfig.DatabaseConfig{MaxOpenConns: 12, MaxIdleConns: 8})
	if got := db.Stats().MaxOpenConnections; got != 12 {
		t.Fatalf("max open=%d", got)
	}
}

func TestValidateServerInfoRequiresTLSAndWritablePrimary(t *testing.T) {
	tests := []struct {
		name string
		info ServerInfo
		want string
	}{
		{name: "no tls", info: ServerInfo{Version: "9.7.1", Cipher: "", TimeZone: "UTC"}, want: "TLS cipher"},
		{name: "read only", info: ServerInfo{Version: "9.7.1", Cipher: "TLS_AES_128_GCM_SHA256", ReadOnly: true, TimeZone: "UTC"}, want: "read-only"},
		{name: "not utc", info: ServerInfo{Version: "9.7.1", Cipher: "TLS_AES_128_GCM_SHA256", TimeZone: "+08:00"}, want: "UTC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateServerInfo(tt.info)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err=%v want substring %q", err, tt.want)
			}
		})
	}
	if err := validateServerInfo(ServerInfo{Version: "9.7.1", Cipher: "TLS_AES_128_GCM_SHA256", TimeZone: "UTC"}); err != nil {
		t.Fatalf("valid server rejected: %v", err)
	}
}

func TestOpenRejectsMissingCredentialWithoutLeakingSecret(t *testing.T) {
	_, err := Open(context.Background(), appconfig.DatabaseConfig{
		Address:      "127.0.0.1:13306",
		User:         "db-user",
		PasswordFile: "/does/not/exist/unique-secret",
		TLSMode:      "required",
	}, "vpn_nodes")
	if err == nil {
		t.Fatal("expected credential read failure")
	}
	if strings.Contains(err.Error(), "db-user") {
		t.Fatalf("database user leaked in error: %v", err)
	}
}
