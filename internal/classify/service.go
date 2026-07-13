package classify

import (
	"context"
	"database/sql"
	"net"
	"sort"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/geoip"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
)

type Service struct {
	DB      *sql.DB
	Network *geoip.NetworkClassifier
}
type Report struct{ Candidates, Classified int }

func (s Service) Run(ctx context.Context, limit int, now time.Time) (Report, error) {
	candidates, err := store.ListClassificationCandidates(ctx, s.DB, limit)
	if err != nil {
		return Report{}, err
	}
	percentiles := latencyPercentiles(candidates)
	updates := make([]store.ClassificationUpdate, 0, len(candidates))
	for _, candidate := range candidates {
		verified := candidate.LastValidationAt.Valid
		breakdown := Score(Input{Verified: verified, CurrentAvailable: candidate.Availability == "available",
			Success7D: candidate.Success7D, Stability30D: candidate.Stability30D,
			LatencyPercentile: percentiles[candidate.NodeConfigID], Consistency: candidate.Consistency,
			FreshnessHours: now.Sub(candidate.LastSeenAt).Hours(), SourceCount: candidate.SourceCount,
			ExitStability: candidate.ExitStability, Compatibility: candidate.Compatibility})
		network := s.Network.Classify(net.ParseIP(candidate.EntryHost))
		updates = append(updates, store.ClassificationUpdate{NodeConfigID: candidate.NodeConfigID, Protocol: candidate.Protocol,
			Transport: candidate.Transport, Security: candidate.Security, ExitCountry: candidate.ExitCountry, ExitASN: candidate.ExitASN,
			IPVersion: candidate.IPVersion, EntryCountry: network.Country, EntryRegion: network.Region, EntryCity: network.City,
			EntryTimeZone: network.TimeZone, EntryASN: network.ASN, EntryOrganization: network.Organization, ProviderClass: network.ProviderClass,
			FreshnessClass: freshnessClass(now.Sub(candidate.LastSeenAt)), StabilityClass: stabilityClass(candidate.Success7D),
			Score: breakdown.Total, Grade: breakdown.Grade, Breakdown: breakdown})
	}
	if err := store.WriteClassifications(ctx, s.DB, updates, now); err != nil {
		return Report{}, err
	}
	return Report{Candidates: len(candidates), Classified: len(updates)}, nil
}

func latencyPercentiles(candidates []store.ClassificationCandidate) map[uint64]float64 {
	groups := make(map[string][]store.ClassificationCandidate)
	for _, candidate := range candidates {
		if candidate.LatencyMS > 0 {
			groups[candidate.Protocol] = append(groups[candidate.Protocol], candidate)
		}
	}
	result := make(map[uint64]float64)
	for _, group := range groups {
		sort.Slice(group, func(i, j int) bool { return group[i].LatencyMS < group[j].LatencyMS })
		for index, candidate := range group {
			percentile := float64(0)
			if len(group) > 1 {
				percentile = float64(index) / float64(len(group)-1)
			}
			result[candidate.NodeConfigID] = percentile
		}
	}
	return result
}

func freshnessClass(age time.Duration) string {
	if age <= 24*time.Hour {
		return "new"
	}
	if age <= 7*24*time.Hour {
		return "fresh"
	}
	if age <= 30*24*time.Hour {
		return "aging"
	}
	return "old"
}
func stabilityClass(rate float64) string {
	if rate >= .9 {
		return "stable"
	}
	if rate >= .5 {
		return "flaky"
	}
	return "unstable"
}
