package source

import (
	"context"
	"fmt"
	"os"
)

// FilesystemSource loads configuration from the local filesystem.
type FilesystemSource struct{}

// NewFilesystemSource constructs a FilesystemSource.
func NewFilesystemSource() *FilesystemSource {
	return &FilesystemSource{}
}

// Load reads the file at req.Path and returns its content.
func (f *FilesystemSource) Load(_ context.Context, req LoadRequest) ([]byte, error) {
	if req.Path == "" {
		return nil, fmt.Errorf("filesystem source: path is required")
	}
	data, err := os.ReadFile(req.Path)
	if err != nil {
		return nil, fmt.Errorf("filesystem source: read %q: %w", req.Path, err)
	}
	return data, nil
}
