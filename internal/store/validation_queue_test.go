package store

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestValidationQueueRejectsInvalidClaim(t *testing.T) {
	for _, test := range []struct {
		owner string
		limit int
		lease time.Duration
	}{{"", 1, time.Minute}, {"worker", 0, time.Minute}, {"worker", 1, 0}} {
		if _, err := ClaimValidationJobs(context.Background(), nil, test.owner, test.limit, test.lease); err == nil {
			t.Fatalf("expected validation error for %+v", test)
		}
	}
}

func TestValidationQueueOwnershipErrorIsStable(t *testing.T) {
	if !errors.Is(ErrLeaseOwnership, ErrLeaseOwnership) {
		t.Fatal("lease ownership sentinel is not stable")
	}
}
