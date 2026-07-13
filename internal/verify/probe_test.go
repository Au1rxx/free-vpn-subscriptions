package verify

import "testing"

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
