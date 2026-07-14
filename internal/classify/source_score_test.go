package classify

import "testing"

func TestSourceScoreWeights(t *testing.T) {
	perfect := ScoreSource(SourceInput{
		FetchReliability: 1,
		ParseYield:       1,
		UsableNodeRate:   1,
		FreshnessHours:   0,
	})
	if perfect.Total != 100 || perfect.FetchReliability != 35 || perfect.ParseYield != 25 ||
		perfect.UsableNodes != 25 || perfect.Freshness != 15 {
		t.Fatalf("perfect=%+v", perfect)
	}

	half := ScoreSource(SourceInput{
		FetchReliability: .5,
		ParseYield:       .5,
		UsableNodeRate:   .5,
		FreshnessHours:   84,
	})
	if half.Total != 52 || half.FetchReliability != 18 || half.ParseYield != 13 ||
		half.UsableNodes != 13 || half.Freshness != 8 {
		t.Fatalf("half=%+v", half)
	}
}

func TestSourceScoreMissingAndOutOfRangeInputs(t *testing.T) {
	missing := ScoreSource(SourceInput{FreshnessHours: -1})
	if missing.Total != 0 {
		t.Fatalf("missing=%+v", missing)
	}

	clamped := ScoreSource(SourceInput{
		FetchReliability: 2,
		ParseYield:       -1,
		UsableNodeRate:   4,
		FreshnessHours:   400,
	})
	if clamped.Total != 60 || clamped.FetchReliability != 35 || clamped.ParseYield != 0 ||
		clamped.UsableNodes != 25 || clamped.Freshness != 0 {
		t.Fatalf("clamped=%+v", clamped)
	}
}
