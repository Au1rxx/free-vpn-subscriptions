package store

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestValidationResultRejectsMissingIdentity(t *testing.T) {
	if err := PersistValidationResult(context.Background(), nil, ValidationWrite{}); err == nil {
		t.Fatal("missing validation identity was accepted")
	}
}

func TestValidationResultLocksNodeBeforeQueue(t *testing.T) {
	source, err := os.ReadFile("validation_results.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	nodeLock := "SELECT consecutive_failures FROM node_configs WHERE node_config_id=? FOR UPDATE"
	queueLock := "FROM validation_queue WHERE validation_job_id=? FOR UPDATE"
	nodeIndex := strings.Index(text, nodeLock)
	queueIndex := strings.Index(text, queueLock)
	if nodeIndex < 0 || queueIndex < 0 {
		t.Fatalf("required row locks are missing: node=%d queue=%d", nodeIndex, queueIndex)
	}
	if nodeIndex > queueIndex {
		t.Fatalf("validation lock order must be node_configs before validation_queue: node=%d queue=%d", nodeIndex, queueIndex)
	}
}
