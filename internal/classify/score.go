package classify

import "math"

type Input struct {
	Verified, CurrentAvailable                 bool
	Success7D, Stability30D, LatencyPercentile float64
	Consistency, ExitStability, Compatibility  float64
	FreshnessHours                             float64
	SourceCount                                int
}

type Breakdown struct {
	Current, Success7D, Stability30D, Latency, Consistency   int
	Freshness, SourceDiversity, ExitStability, Compatibility int
	Total                                                    int
	Grade                                                    string
}

func Score(input Input) Breakdown {
	result := Breakdown{}
	if input.CurrentAvailable {
		result.Current = 25
	}
	result.Success7D = weighted(input.Success7D, 20)
	result.Stability30D = weighted(input.Stability30D, 10)
	result.Latency = weighted(1-clamp(input.LatencyPercentile), 15)
	result.Consistency = weighted(input.Consistency, 10)
	if input.FreshnessHours >= 0 {
		result.Freshness = weighted(1-clamp(input.FreshnessHours/(7*24)), 8)
	}
	if input.SourceCount > 0 {
		result.SourceDiversity = min(input.SourceCount, 5)
	}
	result.ExitStability = weighted(input.ExitStability, 4)
	result.Compatibility = weighted(input.Compatibility, 3)
	result.Total = result.Current + result.Success7D + result.Stability30D + result.Latency + result.Consistency + result.Freshness + result.SourceDiversity + result.ExitStability + result.Compatibility
	result.Total = min(max(result.Total, 0), 100)
	result.Grade = gradeFor(result.Total, input.Verified)
	return result
}

func weighted(value float64, maximum int) int {
	return int(math.Round(clamp(value) * float64(maximum)))
}
func clamp(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}
func gradeFor(score int, verified bool) string {
	if !verified {
		return "U"
	}
	switch {
	case score >= 90:
		return "S"
	case score >= 80:
		return "A"
	case score >= 65:
		return "B"
	case score >= 50:
		return "C"
	default:
		return "D"
	}
}
