package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

func RetryWithJitter(ctx context.Context, attempts int, baseDelay time.Duration, fn func() (*http.Response, error)) (*http.Response, error) {
	var lastErr error

	for range attempts {
		resp, err := fn()
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		lastErr = err

		// stop retry if context canceled
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// jitter: random delay
		jitter := time.Duration(rand.Int63n(int64(baseDelay)))
		time.Sleep(baseDelay + jitter)
	}

	return nil, lastErr
}

func DoJSONRequest(
	ctx context.Context,
	client *http.Client,
	method, url string,
	headers map[string]string,
	body any,
) (*http.Response, error) {

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return client.Do(req)
}

func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}
