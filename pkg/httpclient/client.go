package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	timeout    time.Duration
	maxRetries int
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		baseURL:    baseURL,
		timeout:    timeout,
		maxRetries: 3,
	}
}

func (c *Client) Get(ctx context.Context, path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Country-Search-API/1.0")
	
	var resp *http.Response
	var lastErr error
	
	// Retry logic
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		resp, lastErr = c.httpClient.Do(req)
		if lastErr == nil && resp.StatusCode < 500 {
			break
		}
		
		if resp != nil {
			resp.Body.Close()
		}
		
		if attempt < c.maxRetries-1 {
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}
	
	if lastErr != nil {
		return nil, fmt.Errorf("request failed after retries: %w", lastErr)
	}
	
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	
	return body, nil
}
