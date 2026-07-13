package validation

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPerformanceSamplerDefaultsCapsAndStopsEarly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.CopyN(w, zeroReader{}, 2<<20)
	}))
	defer server.Close()
	sampler := PerformanceSampler{}
	result := sampler.Sample(context.Background(), server.Client(), SampleRequest{URL: server.URL, Bytes: 2 << 20, Timeout: time.Second})
	if result.ErrorCode != "" || result.Bytes != 1<<20 || result.BytesPerSecond <= 0 {
		t.Fatalf("unexpected sample: %+v", result)
	}
	result = sampler.Sample(context.Background(), server.Client(), SampleRequest{URL: server.URL, Timeout: time.Second})
	if result.Bytes != 256<<10 {
		t.Fatalf("default bytes=%d", result.Bytes)
	}
}

type zeroReader struct{}

func (zeroReader) Read(buffer []byte) (int, error) {
	for index := range buffer {
		buffer[index] = 0
	}
	return len(buffer), nil
}
