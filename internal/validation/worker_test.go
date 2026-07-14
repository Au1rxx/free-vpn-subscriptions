package validation

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/store"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/verify"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

type fakeQueue struct {
	jobs      []store.ValidationJob
	mu        sync.Mutex
	persisted int
}

func (q *fakeQueue) Claim(context.Context, string, int, time.Duration) ([]store.ValidationJob, error) {
	return q.jobs, nil
}
func (q *fakeQueue) Extend(context.Context, uint64, string, time.Duration) error { return nil }
func (q *fakeQueue) Persist(_ context.Context, _ store.ValidationWrite) error {
	q.mu.Lock()
	q.persisted++
	q.mu.Unlock()
	return nil
}

type fakeEngine struct {
	mu                     sync.Mutex
	performanceProbeStates []bool
}

func (e *fakeEngine) Verify(_ context.Context, _ *node.Node, request verify.Request) verify.Result {
	e.mu.Lock()
	e.performanceProbeStates = append(e.performanceProbeStates, request.PerformanceProbe != nil)
	e.mu.Unlock()
	return verify.Result{ConfigAccepted: true, ProxyStarted: true, Passed: true, Attempts: 2, Successes: 2}
}

func TestWorkerRunOnceProcessesClaimedJobs(t *testing.T) {
	configured := &node.Node{Protocol: node.ProtoVLESS, Server: "192.0.2.1", Port: 443, UUID: "id"}
	body, _ := json.Marshal(configured)
	queue := &fakeQueue{jobs: []store.ValidationJob{{ID: 1, NodeConfigID: 2, Stage: "connectivity", NormalizedConfig: body}}}
	prober := &fakeProber{}
	engine := &fakeEngine{}
	worker := Worker{Queue: queue, Engine: engine, Checker: QuickChecker{Prober: prober, Resolver: fakeResolver{}}, ValidatorID: "ai-a1", Concurrency: 2}
	report, err := worker.RunOnce(context.Background(), 10)
	if err != nil || report.Claimed != 1 || report.Passed != 1 || queue.persisted != 1 {
		t.Fatalf("report=%+v persisted=%d err=%v", report, queue.persisted, err)
	}
}

func TestWorkerSchedulesPerformanceProbeOnlyWhenDue(t *testing.T) {
	configured := &node.Node{Protocol: node.ProtoVLESS, Server: "192.0.2.1", Port: 443, UUID: "id"}
	body, _ := json.Marshal(configured)
	queue := &fakeQueue{jobs: []store.ValidationJob{
		{ID: 1, NodeConfigID: 1, Stage: "connectivity", NormalizedConfig: body, PerformanceDue: true},
		{ID: 2, NodeConfigID: 2, Stage: "connectivity", NormalizedConfig: body, PerformanceDue: false},
	}}
	engine := &fakeEngine{}
	worker := Worker{Queue: queue, Engine: engine, Checker: QuickChecker{Prober: &fakeProber{}, Resolver: fakeResolver{}},
		ValidatorID: "ai-a1", Concurrency: 1,
		Performance: SampleRequest{URL: "https://speed.example.test/sample", Bytes: 256 << 10, Timeout: time.Second}}
	if _, err := worker.RunOnce(context.Background(), 10); err != nil {
		t.Fatal(err)
	}
	if len(engine.performanceProbeStates) != 2 || !engine.performanceProbeStates[0] || engine.performanceProbeStates[1] {
		t.Fatalf("performance probe states=%v", engine.performanceProbeStates)
	}
}
