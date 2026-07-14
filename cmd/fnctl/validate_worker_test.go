package main

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateWorkerInitialProtocolFlagDefaultsEmpty(t *testing.T) {
	command := newValidateWorkerCmd()
	flag := command.Flags().Lookup("initial-protocol")
	if flag == nil {
		t.Fatal("missing --initial-protocol flag")
	}
	if flag.DefValue != "" {
		t.Fatalf("initial protocol default=%q, want empty", flag.DefValue)
	}
}

func TestValidateWorkerRejectsUnsupportedInitialProtocolBeforeConfigLoad(t *testing.T) {
	previousConfigPath := cfgPath
	cfgPath = filepath.Join(t.TempDir(), "missing.yaml")
	t.Cleanup(func() { cfgPath = previousConfigPath })

	command := newValidateWorkerCmd()
	command.SetArgs([]string{"--initial-protocol", "invalid"})
	err := command.ExecuteContext(context.Background())
	if err == nil || !strings.Contains(err.Error(), "unsupported initial protocol") {
		t.Fatalf("error=%v, want unsupported initial protocol before config load", err)
	}
	if strings.Contains(err.Error(), "missing.yaml") {
		t.Fatalf("validated config before protocol: %v", err)
	}
}
