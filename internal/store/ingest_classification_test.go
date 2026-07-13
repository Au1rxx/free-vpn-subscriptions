package store

import (
	"strings"
	"testing"
)

func TestClassificationValueIsDatabaseBounded(t *testing.T) {
	if got := classificationValue("", "none"); got != "none" {
		t.Fatalf("empty=%q", got)
	}
	if got := classificationValue("TLS", "none"); got != "tls" {
		t.Fatalf("normalized=%q", got)
	}
	if got := classificationValue(strings.Repeat("x", 33), "none"); got != "other" {
		t.Fatalf("oversized=%q", got)
	}
	if got := classificationValue("tls\ninvalid", "none"); got != "other" {
		t.Fatalf("invalid=%q", got)
	}
}
