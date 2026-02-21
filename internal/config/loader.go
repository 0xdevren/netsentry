// Package config provides device configuration loading and source abstraction.
package config

import (
	"context"
	"fmt"

	"github.com/0xdevren/netsentry/internal/config/source"
)

// LoadOptions configures a configuration load operation.
type LoadOptions struct {
	// Source identifies the config source type: "filesystem", "ssh", "api", "git".
	Source string
	// Path is the filesystem path or remote URL/identifier for the configuration.
	Path string
	// SSHOptions is only used when Source is "ssh".
	SSHOptions *source.SSHOptions
	// APIOptions is only used when Source is "api".
	APIOptions *source.APIOptions
	// GitOptions is only used when Source is "git".
	GitOptions *source.GitOptions
}

// Loader abstracts configuration retrieval from multiple source backends.
type Loader struct {
	sources map[string]source.ConfigSource
}

// NewLoader constructs a Loader pre-registered with all built-in sources.
func NewLoader() *Loader {
	l := &Loader{sources: make(map[string]source.ConfigSource)}
	l.sources["filesystem"] = source.NewFilesystemSource()
	l.sources["ssh"] = source.NewSSHSource()
	l.sources["api"] = source.NewAPISource()
	l.sources["git"] = source.NewGitSource()
	return l
}

// Load retrieves the raw configuration text from the configured source.
func (l *Loader) Load(ctx context.Context, opts LoadOptions) ([]byte, error) {
	src, ok := l.sources[opts.Source]
	if !ok {
		return nil, fmt.Errorf("config loader: unknown source %q", opts.Source)
	}
	req := source.LoadRequest{
		Path:       opts.Path,
		SSHOptions: opts.SSHOptions,
		APIOptions: opts.APIOptions,
		GitOptions: opts.GitOptions,
	}
	data, err := src.Load(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("config loader: source %q: %w", opts.Source, err)
	}
	return data, nil
}
