package verify

import (
	"context"
	"net/http"
	"testing"
)

func TestFinalizeEngineResultPartialAndPassed(t *testing.T) {
	result := Result{Attempts: 2, Targets: []TargetResult{{OK: true, TotalMS: 10}, {OK: false, ErrorCode: "timeout"}}}
	finalizeEngineResult(&result)
	if result.Passed || !result.PartialSuccess || result.Successes != 1 || result.HTTPMedianMS != 10 {
		t.Fatalf("partial result: %+v", result)
	}
	result = Result{Attempts: 2, Targets: []TargetResult{{OK: true, TotalMS: 10}, {OK: true, TotalMS: 20}}}
	finalizeEngineResult(&result)
	if !result.Passed || result.PartialSuccess || result.HTTPMedianMS != 15 {
		t.Fatalf("passed result: %+v", result)
	}
}

func TestTargetURLHashDoesNotExposeURL(t *testing.T) {
	hash := targetURLHash("https://example.com/path?token=secret")
	if len(hash) != 16 || hash == "https://example.com/path?token=secret" {
		t.Fatalf("unsafe target hash %q", hash)
	}
}

func TestPerformanceProbeRunsOnlyAfterConnectivitySuccess(t *testing.T) {
	called := 0
	request := Request{PerformanceProbe: func(context.Context, ProxyDialer) PerformanceResult {
		called++
		return PerformanceResult{Attempted: true, Bytes: 256 << 10, BytesPerSecond: 1024}
	}}
	client := &http.Client{}
	failed := Result{}
	runPerformanceProbe(context.Background(), client, request, &failed)
	if called != 0 || failed.Performance.Attempted {
		t.Fatalf("failed connectivity sampled: called=%d result=%+v", called, failed.Performance)
	}
	passed := Result{Successes: 1}
	runPerformanceProbe(context.Background(), client, request, &passed)
	if called != 1 || passed.Performance.Bytes != 256<<10 || passed.Performance.BytesPerSecond != 1024 {
		t.Fatalf("performance callback result: called=%d result=%+v", called, passed.Performance)
	}
}
