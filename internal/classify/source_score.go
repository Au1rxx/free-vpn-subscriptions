package classify

// SourceInput contains normalized source-health measurements. Negative
// FreshnessHours means the source has never completed a successful fetch.
type SourceInput struct {
	FetchReliability float64
	ParseYield       float64
	UsableNodeRate   float64
	FreshnessHours   float64
}

// SourceBreakdown is the explainable 0-100 source-quality score.
type SourceBreakdown struct {
	FetchReliability int
	ParseYield       int
	UsableNodes      int
	Freshness        int
	Total            int
}

// ScoreSource applies the source-quality weights confirmed by the platform
// design: fetch 35, parse 25, usable nodes 25, and freshness 15.
func ScoreSource(input SourceInput) SourceBreakdown {
	result := SourceBreakdown{
		FetchReliability: weighted(input.FetchReliability, 35),
		ParseYield:       weighted(input.ParseYield, 25),
		UsableNodes:      weighted(input.UsableNodeRate, 25),
	}
	if input.FreshnessHours >= 0 {
		result.Freshness = weighted(1-clamp(input.FreshnessHours/(7*24)), 15)
	}
	result.Total = result.FetchReliability + result.ParseYield + result.UsableNodes + result.Freshness
	result.Total = min(max(result.Total, 0), 100)
	return result
}
