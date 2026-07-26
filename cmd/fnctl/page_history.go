package main

import (
	"fmt"
	"io"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/aggregate"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/pages"
)

func preparePageHistory(path string, summary aggregate.Summary, minPerCountry int, warnings io.Writer) []pages.HistoryPoint {
	current := pages.HistoryPoint{
		GeneratedAt:     time.Unix(summary.GeneratedAtUnix, 0).UTC(),
		Selected:        summary.TotalSelected,
		Verified:        summary.TotalVerified,
		MedianLatencyMS: summary.MedianLatencyMS,
		Countries:       qualifyingCountryCount(summary.ByCountry, minPerCountry),
	}
	points, err := pages.UpdateHistory(path, current)
	if err == nil {
		return points
	}
	if warnings != nil {
		fmt.Fprintf(warnings, "warning: page history: %v\n", err)
	}
	return []pages.HistoryPoint{current}
}

func qualifyingCountryCount(byCountry map[string]int, minPerCountry int) int {
	count := 0
	for country, nodes := range byCountry {
		if country == "" || country == "XX" || nodes <= 0 || nodes < minPerCountry {
			continue
		}
		count++
	}
	return count
}
