package ingest

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/discovery"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/sources"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/parse"
)

const parserVersion = "fnctl-2"

type Service struct {
	DB    *sql.DB
	Spool *Spool
}

type ImportSummary struct{ Sources, InsertedOrUpdated int }
type FetchSummary struct{ Sources, Success, NotModified, Failed, Spooled, Bytes int }
type ParseSummary struct{ Fetches, Nodes, Errors, NewEndpoints, NewConfigs, QueueJobs, Discovered int }

func (s *Service) ImportSeeds(ctx context.Context, configured []config.Source) (ImportSummary, error) {
	summary := ImportSummary{Sources: len(configured)}
	for _, seed := range configured {
		interval := time.Duration(seed.FetchIntervalSeconds) * time.Second
		if _, err := store.UpsertSource(ctx, s.DB, store.SourceRecord{
			Name: seed.Name, URL: seed.URL, FormatHint: seed.Format, Enabled: seed.Enabled,
			Kind: seed.Kind, DiscoveryMethod: seed.DiscoveryMethod, State: "active", Depth: seed.Depth,
			Priority: seed.Priority, FetchInterval: interval,
		}); err != nil {
			return summary, err
		}
		summary.InsertedOrUpdated++
	}
	return summary, nil
}

func (s *Service) Fetch(ctx context.Context, limit int) (FetchSummary, error) {
	claimed, err := store.ClaimDueSources(ctx, s.DB, limit)
	if err != nil {
		return FetchSummary{}, err
	}
	summary := FetchSummary{Sources: len(claimed)}
	for _, source := range claimed {
		started := time.Now().UTC()
		response, fetchErr := sources.FetchRaw(ctx, sources.Request{
			URL: source.URL, ETag: source.ETag, LastModified: source.LastModified,
			Timeout: 30 * time.Second, MaxBodyBytes: 20 << 20, MaxDecodedBytes: 64 << 20, MaxRedirects: 5,
		})
		write := store.FetchWrite{SourceID: source.ID, StartedAt: started, FinishedAt: time.Now().UTC()}
		if fetchErr != nil {
			write.ErrorCode, write.ErrorSummary = sources.ErrorCode(fetchErr), fetchErr.Error()
			summary.Failed++
		} else {
			write.StatusCode, write.FinalURL, write.ETag = response.StatusCode, response.FinalURL, response.ETag
			write.LastModified, write.ContentType, write.ContentEncoding = response.LastModified, response.ContentType, response.ContentEncoding
			write.Body, write.Duration = response.Body, response.Duration
			summary.Bytes += len(response.Body)
			if response.StatusCode == 304 {
				summary.NotModified++
			} else {
				summary.Success++
			}
		}
		if _, err := store.FinishFetch(ctx, s.DB, write); err != nil {
			if s.Spool == nil {
				return summary, err
			}
			if spoolErr := s.Spool.Enqueue(envelopeFromWrite(write)); spoolErr != nil {
				return summary, fmt.Errorf("persist fetch: %v; spool: %w", err, spoolErr)
			}
			summary.Spooled++
		}
	}
	return summary, nil
}

func (s *Service) Parse(ctx context.Context, limit int) (ParseSummary, error) {
	inputs, err := store.ClaimUnparsedFetches(ctx, s.DB, limit)
	if err != nil {
		return ParseSummary{}, err
	}
	summary := ParseSummary{Fetches: len(inputs)}
	for _, input := range inputs {
		result := parse.Parse(input.Body, parse.Format(input.FormatHint))
		persisted, err := store.PersistParseResult(ctx, s.DB, input.SourceID, input.FetchID, result, parserVersion)
		if err != nil {
			return summary, err
		}
		summary.Nodes += len(result.Nodes)
		summary.Errors += len(result.Errors)
		summary.NewEndpoints += persisted.NewEndpoints
		summary.NewConfigs += persisted.NewConfigs
		summary.QueueJobs += persisted.QueueJobs
		for _, discoveredURL := range result.DiscoveredURLs {
			if !discovery.LikelySubscriptionURL(discoveredURL) {
				continue
			}
			if _, err := store.UpsertSource(ctx, s.DB, store.SourceRecord{
				Name: "nested-source", URL: discoveredURL, FormatHint: "auto", Enabled: true,
				Kind: "nested-subscription", DiscoveryMethod: "content-link", State: "active", Depth: 1,
			}); err == nil {
				summary.Discovered++
			}
		}
	}
	return summary, nil
}

// PersistFetch implements spool.Persister.
func (s *Service) PersistFetch(ctx context.Context, envelope FetchEnvelope) error {
	_, err := store.FinishFetch(ctx, s.DB, writeFromEnvelope(envelope))
	return err
}

func envelopeFromWrite(write store.FetchWrite) FetchEnvelope {
	return FetchEnvelope{SourceID: write.SourceID, FetchedAt: write.FinishedAt, StatusCode: write.StatusCode,
		FinalURL: write.FinalURL, ETag: write.ETag, LastModified: write.LastModified,
		ContentType: write.ContentType, ContentEncoding: write.ContentEncoding, Body: write.Body,
		DurationMS: uint64(write.Duration / time.Millisecond), ErrorCode: write.ErrorCode, ErrorSummary: write.ErrorSummary}
}

func writeFromEnvelope(envelope FetchEnvelope) store.FetchWrite {
	return store.FetchWrite{SourceID: envelope.SourceID, StartedAt: envelope.FetchedAt, FinishedAt: envelope.FetchedAt,
		StatusCode: envelope.StatusCode, FinalURL: envelope.FinalURL, ETag: envelope.ETag,
		LastModified: envelope.LastModified, ContentType: envelope.ContentType, ContentEncoding: envelope.ContentEncoding,
		Body: envelope.Body, Duration: time.Duration(envelope.DurationMS) * time.Millisecond,
		ErrorCode: envelope.ErrorCode, ErrorSummary: envelope.ErrorSummary}
}
