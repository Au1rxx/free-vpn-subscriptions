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

const parserVersion = "fnctl-4"

type Service struct {
	DB      *sql.DB
	Spool   *Spool
	Archive *PayloadArchive
	// Response ceilings for a single upstream fetch. Zero falls back to the
	// sources package defaults.
	MaxBodyBytes, MaxDecodedBytes int64
	ArchiveThresholdBytes         int64
	ArchiveWriteEnabled           bool
}

type ImportSummary struct{ Sources, InsertedOrUpdated int }
type FetchSummary struct {
	Sources, Success, NotModified, Failed, Spooled, Replayed, Quarantined, Bytes int
}
type ParseSummary struct{ Fetches, Nodes, Errors, NewEndpoints, NewConfigs, QueueJobs, Discovered int }

func (s *Service) ImportSeeds(ctx context.Context, configured []config.Source) (ImportSummary, error) {
	summary := ImportSummary{Sources: len(configured)}
	for _, seed := range configured {
		interval := time.Duration(seed.FetchIntervalSeconds) * time.Second
		if _, err := store.UpsertSource(ctx, s.DB, store.SourceRecord{
			Name: seed.Name, URL: seed.URL, FormatHint: seed.Format, ProtocolHint: seed.ProtocolHint, Enabled: seed.Enabled,
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
	summary := FetchSummary{}
	replay, err := s.ReplaySpool(ctx)
	summary.Replayed, summary.Quarantined = replay.Persisted, replay.Quarantined
	if err != nil {
		return summary, fmt.Errorf("replay fetch spool: %w", err)
	}
	claimed, err := store.ClaimDueSources(ctx, s.DB, limit)
	if err != nil {
		return summary, err
	}
	summary.Sources = len(claimed)
	for _, source := range claimed {
		started := time.Now().UTC()
		response, fetchErr := sources.FetchRaw(ctx, sources.Request{
			URL: source.URL, ETag: source.ETag, LastModified: source.LastModified,
			Timeout: 30 * time.Second, MaxRedirects: 5,
			MaxBodyBytes: s.MaxBodyBytes, MaxDecodedBytes: s.MaxDecodedBytes,
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
		persistWrite, persistErr := persistPreparedFetch(ctx, write, s.ArchiveWriteEnabled, s.ArchiveThresholdBytes,
			func(ctx context.Context, candidate store.FetchWrite) (store.FetchWrite, error) {
				return archiveFetchWrite(ctx, s.Archive, candidate)
			},
			func(ctx context.Context, candidate store.FetchWrite) error {
				_, err := store.FinishFetch(ctx, s.DB, candidate)
				return err
			})
		if persistWrite.ErrorCode == "payload_storage_required" && write.StatusCode >= 200 && write.StatusCode < 300 && write.StatusCode != 304 {
			summary.Success--
			summary.Failed++
		}
		if persistErr != nil {
			if s.Spool == nil {
				return summary, persistErr
			}
			if spoolErr := s.Spool.Enqueue(envelopeFromWrite(persistWrite)); spoolErr != nil {
				return summary, fmt.Errorf("persist fetch: %v; spool: %w", persistErr, spoolErr)
			}
			summary.Spooled++
		}
	}
	return summary, nil
}

// ReplaySpool persists fetches retained during a prior transient database
// failure before claiming new network work.
func (s *Service) ReplaySpool(ctx context.Context) (ReplayReport, error) {
	if s.Spool == nil {
		return ReplayReport{}, nil
	}
	return s.Spool.Replay(ctx, s)
}

func (s *Service) Parse(ctx context.Context, limit int) (ParseSummary, error) {
	return s.ParseStorage(ctx, store.StorageDatabase, limit)
}

func (s *Service) ParseStorage(ctx context.Context, storageKind string, limit int) (ParseSummary, error) {
	inputs, err := store.ClaimUnparsedFetchesByStorage(ctx, s.DB, storageKind, limit)
	if err != nil {
		return ParseSummary{}, err
	}
	summary := ParseSummary{Fetches: len(inputs)}
	for _, input := range inputs {
		if input.StorageKind == store.StorageFilesystem {
			if s.Archive == nil {
				return summary, fmt.Errorf("read filesystem fetch %d: payload archive is unavailable", input.FetchID)
			}
			input.Body, err = s.Archive.Read(ctx, ArchiveReference{
				Key: input.ArchiveKey, SHA256: input.PayloadHash,
				OriginalBytes: input.OriginalBytes, CompressedBytes: input.CompressedBytes,
			})
			if err != nil {
				return summary, fmt.Errorf("read filesystem fetch %d: %w", input.FetchID, err)
			}
		}
		result := parse.ParseWithProtocolHint(input.Body, parse.Format(input.FormatHint), input.ProtocolHint)
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

func (s *Service) RecoveredFetchExists(ctx context.Context, sourceID uint64, fetchedAt time.Time, digest [32]byte) (bool, error) {
	return store.RecoveredFetchExists(ctx, s.DB, sourceID, fetchedAt, digest)
}

func (s *Service) PersistRecoveredFetch(ctx context.Context, envelope FetchEnvelope) error {
	return s.PersistFetch(ctx, envelope)
}

func envelopeFromWrite(write store.FetchWrite) FetchEnvelope {
	return FetchEnvelope{SourceID: write.SourceID, FetchedAt: write.FinishedAt, StatusCode: write.StatusCode,
		FinalURL: write.FinalURL, ETag: write.ETag, LastModified: write.LastModified,
		ContentType: write.ContentType, ContentEncoding: write.ContentEncoding, Body: write.Body,
		ExternalPayload: write.ExternalPayload, ResponseBytes: write.ResponseBytes,
		DurationMS: uint64(write.Duration / time.Millisecond), ErrorCode: write.ErrorCode, ErrorSummary: write.ErrorSummary}
}

func writeFromEnvelope(envelope FetchEnvelope) store.FetchWrite {
	return store.FetchWrite{SourceID: envelope.SourceID, StartedAt: envelope.FetchedAt, FinishedAt: envelope.FetchedAt,
		StatusCode: envelope.StatusCode, FinalURL: envelope.FinalURL, ETag: envelope.ETag,
		LastModified: envelope.LastModified, ContentType: envelope.ContentType, ContentEncoding: envelope.ContentEncoding,
		Body: envelope.Body, ExternalPayload: envelope.ExternalPayload, ResponseBytes: envelope.ResponseBytes,
		Duration:  time.Duration(envelope.DurationMS) * time.Millisecond,
		ErrorCode: envelope.ErrorCode, ErrorSummary: envelope.ErrorSummary}
}

func storageFailureWrite(write store.FetchWrite, cause error) store.FetchWrite {
	responseBytes := int64(len(write.Body))
	write.Body = nil
	write.ExternalPayload = nil
	write.ResponseBytes = responseBytes
	write.ErrorCode = "payload_storage_required"
	write.ErrorSummary = cause.Error()
	return write
}

type archiveFetchFunc func(context.Context, store.FetchWrite) (store.FetchWrite, error)
type persistFetchFunc func(context.Context, store.FetchWrite) error

func persistPreparedFetch(ctx context.Context, write store.FetchWrite, archiveEnabled bool, archiveThresholdBytes int64,
	archive archiveFetchFunc, persist persistFetchFunc) (store.FetchWrite, error) {
	candidate := write
	if archiveEnabled && shouldArchiveBody(int64(len(write.Body)), archiveThresholdBytes) {
		archived, err := archive(ctx, write)
		if err != nil {
			failure := storageFailureWrite(write, err)
			return failure, persist(ctx, failure)
		}
		candidate = archived
	}
	if err := persist(ctx, candidate); err != nil {
		if !store.IsDatabasePayloadTooLarge(err) {
			return candidate, err
		}
		storageErr := err
		if archiveEnabled {
			archived, archiveErr := archive(ctx, write)
			if archiveErr == nil {
				candidate = archived
				return candidate, persist(ctx, candidate)
			}
			storageErr = archiveErr
		}
		failure := storageFailureWrite(write, storageErr)
		return failure, persist(ctx, failure)
	}
	return candidate, nil
}
