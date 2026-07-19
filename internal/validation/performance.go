package validation

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	defaultSampleBytes = int64(256 << 10)
	maximumSampleBytes = int64(1 << 20)
)

type SampleRequest struct {
	URL     string
	Bytes   int64
	Timeout time.Duration
}

type SampleResult struct {
	Bytes          int64
	Duration       time.Duration
	BytesPerSecond int64
	ErrorCode      string
}

type ProxyDialer interface {
	Do(*http.Request) (*http.Response, error)
}

type PerformanceSampler struct{}

func (PerformanceSampler) Sample(ctx context.Context, dialer ProxyDialer, request SampleRequest) SampleResult {
	if request.Bytes <= 0 {
		request.Bytes = defaultSampleBytes
	}
	if request.Bytes > maximumSampleBytes {
		request.Bytes = maximumSampleBytes
	}
	if request.Timeout <= 0 {
		request.Timeout = 15 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, request.Timeout)
	defer cancel()
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, request.URL, nil)
	if err != nil {
		return SampleResult{ErrorCode: "invalid_sample_url"}
	}
	started := time.Now()
	response, err := dialer.Do(httpRequest)
	if err != nil {
		code := "sample_request_failed"
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			code = "timeout"
		}
		return SampleResult{Duration: time.Since(started), ErrorCode: code}
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return SampleResult{Duration: time.Since(started), ErrorCode: "http_status"}
	}
	read, err := io.Copy(io.Discard, io.LimitReader(response.Body, request.Bytes))
	duration := time.Since(started)
	result := SampleResult{Bytes: read, Duration: duration}
	if err != nil {
		result.ErrorCode = "sample_read_failed"
		return result
	}
	if read == 0 {
		result.ErrorCode = "empty_sample"
		return result
	}
	if duration > 0 {
		result.BytesPerSecond = int64(float64(read) / duration.Seconds())
	}
	return result
}
