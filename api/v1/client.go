package v1

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ApiClient struct {
	token  string
	client HTTPClient
}

func NewApiClient(token string) *ApiClient {
	return &ApiClient{
		token: token,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// fetchWithRetry performs a GET request with retries and exponential backoff
func (a *ApiClient) fetchWithRetry(url string, maxTries int) ([]byte, error) {
	var lastErr error
	backoff := 500 * time.Millisecond

	for attempt := 1; attempt <= maxTries; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+a.token)

		resp, err := a.client.Do(req)
		if err != nil {
			// Network error, retry
			lastErr = fmt.Errorf("network error on attempt %d: %w", attempt, err)
		} else {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				lastErr = fmt.Errorf("failed to read body on attempt %d: %w", attempt, err)
			} else {
				switch resp.StatusCode {
				case 200:
					return body, nil
				case 429:
					// Too many requests: respect Retry-After header
					retryAfter := resp.Header.Get("Retry-After")
					if retryAfter, err := time.ParseDuration(retryAfter + "s"); err == nil {
						time.Sleep(retryAfter)
					} else {
						time.Sleep(backoff)
					}
				case 500, 502, 503, 504:
					// Server errors, retry
					time.Sleep(backoff)
				default:
					// 4xx errors (except 429) should not retry
					return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
				}
			}
		}

		// exponential backoff for next attempt
		backoff *= 2
	}

	return nil, fmt.Errorf("all attempts failed, last error: %w", lastErr)
}
