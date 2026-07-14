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
