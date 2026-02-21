// Package source defines the ConfigSource interface and all built-in
// source implementations for loading raw device configuration data.
package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LoadRequest encapsulates all parameters needed to retrieve a configuration.
type LoadRequest struct {
	// Path is a filesystem path, URL, or device address depending on the source.
	Path string
	// SSHOptions is populated when the source is "ssh".
	SSHOptions *SSHOptions
	// APIOptions is populated when the source is "api".
	APIOptions *APIOptions
	// GitOptions is populated when the source is "git".
	GitOptions *GitOptions
}

// ConfigSource is the interface that all configuration source backends implement.
type ConfigSource interface {
	// Load retrieves the raw device configuration and returns it as bytes.
	Load(ctx context.Context, req LoadRequest) ([]byte, error)
}

// APIOptions configures an HTTP API config source.
type APIOptions struct {
	// URL is the HTTP endpoint to fetch the config from.
	URL string
	// Token is an optional Bearer token for authentication.
	Token string
	// Method is the HTTP method (defaults to "GET").
	Method string
	// Timeout for the HTTP request.
	Timeout time.Duration
}

// apiSource fetches configurations from an HTTP API.
type apiSource struct{}

// NewAPISource constructs an HTTP API source.
func NewAPISource() ConfigSource {
	return &apiSource{}
}

// Load performs an HTTP request to fetch the configuration data.
func (a *apiSource) Load(ctx context.Context, req LoadRequest) ([]byte, error) {
	opts := req.APIOptions
	if opts == nil {
		return nil, fmt.Errorf("api source: APIOptions are required")
	}
	url := opts.URL
	if url == "" {
		url = req.Path
	}
	if url == "" {
		return nil, fmt.Errorf("api source: URL or Path is required")
	}

	method := opts.Method
	if method == "" {
		method = http.MethodGet
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("api source: create request: %w", err)
	}

	if opts.Token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+opts.Token)
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("api source: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("api source: unexpected status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("api source: read response body: %w", err)
	}

	return data, nil
}
