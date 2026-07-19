package main

import "testing"

func TestClassificationBatchSizeRefreshesWhenAllNodesAreClassified(t *testing.T) {
	if got := classificationBatchSize(10000, 0, true); got != 10000 {
		t.Fatalf("batch=%d, want bounded refresh batch", got)
	}
	if got := classificationBatchSize(10000, 5000, true); got != 5000 {
		t.Fatalf("final unclassified batch=%d, want 5000", got)
	}
	if got := classificationBatchSize(10000, 0, false); got != 10000 {
		t.Fatalf("ordinary refresh batch=%d, want 10000", got)
	}
}

func TestRemainingClassificationCountDecrementsWithoutFullRescan(t *testing.T) {
	if got := remainingAfterClassification(2000000, 10000); got != 1990000 {
		t.Fatalf("remaining=%d", got)
	}
	if got := remainingAfterClassification(5000, 5000); got != 0 {
		t.Fatalf("final remaining=%d", got)
	}
	if got := remainingAfterClassification(0, 10000); got != 0 {
		t.Fatalf("refresh changed remaining=%d", got)
	}
}
