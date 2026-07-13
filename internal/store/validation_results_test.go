package store

import (
	"context"
	"testing"
)

func TestValidationResultRejectsMissingIdentity(t *testing.T) {
	if err := PersistValidationResult(context.Background(), nil, ValidationWrite{}); err == nil {
		t.Fatal("missing validation identity was accepted")
	}
}
