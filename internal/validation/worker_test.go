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

type fakeEngine struct{}

func (fakeEngine) Verify(context.Context, *node.Node, verify.Request) verify.Result {
	return verify.Result{ConfigAccepted: true, ProxyStarted: true, Passed: true, Attempts: 2, Successes: 2}
}

func TestWorkerRunOnceProcessesClaimedJobs(t *testing.T) {
	configured := &node.Node{Protocol: node.ProtoVLESS, Server: "192.0.2.1", Port: 443, UUID: "id"}
	body, _ := json.Marshal(configured)
	queue := &fakeQueue{jobs: []store.ValidationJob{{ID: 1, NodeConfigID: 2, Stage: "connectivity", NormalizedConfig: body}}}
	prober := &fakeProber{}
	worker := Worker{Queue: queue, Engine: fakeEngine{}, Checker: QuickChecker{Prober: prober, Resolver: fakeResolver{}}, ValidatorID: "ai-a1", Concurrency: 2}
	report, err := worker.RunOnce(context.Background(), 10)
	if err != nil || report.Claimed != 1 || report.Passed != 1 || queue.persisted != 1 {
		t.Fatalf("report=%+v persisted=%d err=%v", report, queue.persisted, err)
	}
}
