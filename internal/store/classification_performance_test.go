package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
)

func TestListClassificationCandidatesIsBoundedIntegration(t *testing.T) {
	configPath := os.Getenv("VPN_NODE_TEST_CONFIG")
	if configPath == "" {
		t.Skip("VPN_NODE_TEST_CONFIG is not set")
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := Open(ctx, cfg.Database, cfg.Database.Name)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := ListClassificationCandidates(ctx, db, 100); err != nil {
		t.Fatalf("bounded classification candidate query: %v", err)
	}
}
