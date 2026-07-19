package classify

import "testing"

func TestScoreWeightsAndGrades(t *testing.T) {
	perfect := Score(Input{Verified: true, CurrentAvailable: true, Success7D: 1, Stability30D: 1, LatencyPercentile: 0,
		Consistency: 1, FreshnessHours: 0, SourceCount: 5, ExitStability: 1, Compatibility: 1})
	if perfect.Total != 100 || perfect.Grade != "S" {
		t.Fatalf("perfect=%+v", perfect)
	}
	if perfect.Current != 25 || perfect.Success7D != 20 || perfect.Stability30D != 10 || perfect.Latency != 15 || perfect.Consistency != 10 || perfect.Freshness != 8 || perfect.SourceDiversity != 5 || perfect.ExitStability != 4 || perfect.Compatibility != 3 {
		t.Fatalf("weights=%+v", perfect)
	}
	if gradeFor(90, true) != "S" || gradeFor(80, true) != "A" || gradeFor(65, true) != "B" || gradeFor(50, true) != "C" || gradeFor(0, true) != "D" || gradeFor(100, false) != "U" {
		t.Fatal("grade boundaries changed")
	}
}

func TestScoreClampsInputs(t *testing.T) {
	result := Score(Input{Verified: true, Success7D: 2, Stability30D: -1, LatencyPercentile: 2, Consistency: -1, SourceCount: -5})
	if result.Total < 0 || result.Total > 100 || result.Success7D != 20 || result.Stability30D != 0 || result.Latency != 0 {
		t.Fatalf("clamping failed: %+v", result)
	}
}
