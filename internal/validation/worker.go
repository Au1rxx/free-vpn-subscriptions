package validation

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/verify"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

type Queue interface {
	Claim(context.Context, string, int, time.Duration) ([]store.ValidationJob, error)
	Extend(context.Context, uint64, string, time.Duration) error
	Persist(context.Context, store.ValidationWrite) error
}

type VerificationEngine interface {
	Verify(context.Context, *node.Node, verify.Request) verify.Result
}

type SQLQueue struct {
	DB              *sql.DB
	InitialProtocol string
}

func (q SQLQueue) Claim(ctx context.Context, owner string, limit int, lease time.Duration) ([]store.ValidationJob, error) {
	if q.InitialProtocol != "" {
		return store.ClaimInitialValidationJobsByProtocol(ctx, q.DB, owner, q.InitialProtocol, limit, lease)
	}
	return store.ClaimValidationJobs(ctx, q.DB, owner, limit, lease)
}
func (q SQLQueue) Extend(ctx context.Context, jobID uint64, owner string, lease time.Duration) error {
	return store.ExtendLease(ctx, q.DB, jobID, owner, lease)
}
func (q SQLQueue) Persist(ctx context.Context, write store.ValidationWrite) error {
	return store.PersistValidationResult(ctx, q.DB, write)
}

type Worker struct {
	Queue       Queue
	Checker     QuickChecker
	Engine      VerificationEngine
	ValidatorID string
	Concurrency int
	Lease       time.Duration
	Request     verify.Request
	Performance SampleRequest
}

type WorkerReport struct{ Claimed, Passed, Partial, Failed, PersistErrors int }

func (w Worker) RunOnce(ctx context.Context, limit int) (WorkerReport, error) {
	if w.Queue == nil || w.ValidatorID == "" {
		return WorkerReport{}, fmt.Errorf("validation queue and validator ID are required")
	}
	if limit < 1 || limit > 1000 {
		return WorkerReport{}, fmt.Errorf("validation limit must be between 1 and 1000")
	}
	if w.Concurrency <= 0 {
		w.Concurrency = 20
	}
	if w.Lease <= 0 {
		w.Lease = 2 * time.Minute
	}
	if w.Engine == nil {
		w.Engine = verify.Engine{}
	}
	jobs, err := w.Queue.Claim(ctx, w.ValidatorID, limit, w.Lease)
	if err != nil {
		return WorkerReport{}, err
	}
	jobs = fairProtocolOrder(jobs)
	report := WorkerReport{Claimed: len(jobs)}
	semaphore := make(chan struct{}, w.Concurrency)
	var wait sync.WaitGroup
	var mutex sync.Mutex
	var firstPersistError error
	for _, job := range jobs {
		if ctx.Err() != nil {
			break
		}
		job := job
		wait.Add(1)
		semaphore <- struct{}{}
		go func() {
			defer wait.Done()
			defer func() { <-semaphore }()
			passed, partial, persistErr := w.processJob(ctx, job)
			mutex.Lock()
			defer mutex.Unlock()
			switch {
			case persistErr != nil:
				report.PersistErrors++
				if firstPersistError == nil {
					firstPersistError = persistErr
				}
			case passed:
				report.Passed++
			case partial:
				report.Partial++
			default:
				report.Failed++
			}
		}()
	}
	wait.Wait()
	if ctx.Err() != nil {
		return report, ctx.Err()
	}
	if report.PersistErrors > 0 {
		return report, fmt.Errorf("%d validation results failed to persist: %w", report.PersistErrors, firstPersistError)
	}
	return report, nil
}

func (w Worker) processJob(ctx context.Context, job store.ValidationJob) (bool, bool, error) {
	started := time.Now().UTC()
	var configured node.Node
	if err := json.Unmarshal(job.NormalizedConfig, &configured); err != nil {
		result := verify.Result{Protocol: job.Protocol, ErrorCode: "stored_config_invalid", ErrorSummary: boundedValidationError(err)}
		persistErr := w.Queue.Persist(ctx, store.ValidationWrite{JobID: job.ID, NodeConfigID: job.NodeConfigID, Owner: w.ValidatorID,
			ValidatorID: w.ValidatorID, Stage: job.Stage, Engine: "sing-box", StartedAt: started, FinishedAt: time.Now().UTC(), Result: result})
		return false, false, persistErr
	}
	jobCtx, cancel := context.WithCancel(ctx)
	leaseDone := make(chan struct{})
	go func() {
		defer close(leaseDone)
		interval := w.Lease / 2
		if interval < time.Second {
			interval = time.Second
		}
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-jobCtx.Done():
				return
			case <-ticker.C:
				_ = w.Queue.Extend(jobCtx, job.ID, w.ValidatorID, w.Lease)
			}
		}
	}()
	quick := w.Checker.Check(jobCtx, &configured)
	var result verify.Result
	if quick.Passed {
		request := w.Request
		if job.PerformanceDue && w.Performance.URL != "" {
			request.PerformanceProbe = func(ctx context.Context, dialer verify.ProxyDialer) verify.PerformanceResult {
				sample := (PerformanceSampler{}).Sample(ctx, dialer, w.Performance)
				return verify.PerformanceResult{Attempted: true, Bytes: sample.Bytes,
					BytesPerSecond: sample.BytesPerSecond, ErrorCode: sample.ErrorCode}
			}
		}
		result = w.Engine.Verify(jobCtx, &configured, request)
	} else {
		result = verify.Result{Node: &configured, Protocol: configured.Protocol, ErrorCode: quick.ErrorCode,
			ErrorSummary: quick.ErrorSummary, Targets: nil}
	}
	persistErr := w.persistWithRetry(jobCtx, store.ValidationWrite{JobID: job.ID, NodeConfigID: job.NodeConfigID,
		Owner: w.ValidatorID, ValidatorID: w.ValidatorID, Stage: job.Stage, Engine: "sing-box",
		StartedAt: started, FinishedAt: time.Now().UTC(), Result: result})
	cancel()
	<-leaseDone
	return result.Passed, result.PartialSuccess, persistErr
}

func (w Worker) persistWithRetry(ctx context.Context, write store.ValidationWrite) error {
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		if err = w.Queue.Persist(ctx, write); err == nil {
			return nil
		}
		if errors.Is(err, store.ErrLeaseOwnership) || errors.Is(err, store.ErrStaleValidation) || ctx.Err() != nil {
			return err
		}
		delay := time.Duration(attempt+1) * 100 * time.Millisecond
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return err
}

func (w Worker) Run(ctx context.Context, limit int, idleDelay time.Duration) error {
	if idleDelay <= 0 {
		idleDelay = 5 * time.Second
	}
	for ctx.Err() == nil {
		report, err := w.RunOnce(ctx, limit)
		if err != nil {
			return err
		}
		if report.Claimed > 0 {
			continue
		}
		select {
		case <-ctx.Done():
		case <-time.After(idleDelay):
		}
	}
	return ctx.Err()
}

func fairProtocolOrder(jobs []store.ValidationJob) []store.ValidationJob {
	groups := make(map[string][]store.ValidationJob)
	var protocols []string
	for _, job := range jobs {
		if _, exists := groups[job.Protocol]; !exists {
			protocols = append(protocols, job.Protocol)
		}
		groups[job.Protocol] = append(groups[job.Protocol], job)
	}
	ordered := make([]store.ValidationJob, 0, len(jobs))
	for len(ordered) < len(jobs) {
		for _, protocol := range protocols {
			if len(groups[protocol]) == 0 {
				continue
			}
			ordered = append(ordered, groups[protocol][0])
			groups[protocol] = groups[protocol][1:]
		}
	}
	return ordered
}
