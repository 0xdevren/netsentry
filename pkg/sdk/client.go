// Package sdk provides a Go client library for the NetSentry HTTP API.
package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the NetSentry API client.
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// ClientOptions configures the SDK client.
type ClientOptions struct {
	// BaseURL is the NetSentry API base URL (e.g. "http://localhost:8080").
	BaseURL string
	// Token is an optional Bearer token for authenticated endpoints.
	Token string
	// Timeout is the HTTP request timeout (default: 30s).
	Timeout time.Duration
}

// NewClient constructs a Client with the given options.
func NewClient(opts ClientOptions) *Client {
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &Client{
		baseURL:    opts.BaseURL,
		token:      opts.Token,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// ValidateRequest is the input for Client.Validate.
type ValidateRequest struct {
	Config      string `json:"config"`
	PolicyYAML  string `json:"policy_yaml"`
	Strict      bool   `json:"strict"`
	Concurrency int    `json:"concurrency"`
}

// Validate submits a configuration and policy for validation via the API and
// returns the raw JSON report bytes.
func (c *Client) Validate(ctx context.Context, req ValidateRequest) ([]byte, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("sdk: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/validate", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("sdk: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sdk: request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("sdk: read body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("sdk: API error %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}

// Health calls GET /healthz and returns true if the server is healthy.
func (c *Client) Health(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/healthz", nil)
	if err != nil {
		return false, fmt.Errorf("sdk: build request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("sdk: request: %w", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}
