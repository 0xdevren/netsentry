package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileExists reports whether the file or directory at path exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ReadFile reads the entire file at path and returns its content.
func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("filesystem: read %q: %w", path, err)
	}
	return data, nil
}

// WriteFile writes data to the file at path, creating any missing parent
// directories with mode 0o755.
func WriteFile(path string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("filesystem: mkdir %q: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("filesystem: write %q: %w", path, err)
	}
	return nil
}

// EnsureDir creates the directory at path (including all parents) if it does
// not already exist.
func EnsureDir(path string) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("filesystem: ensure dir %q: %w", path, err)
	}
	return nil
}

// AbsPath returns the absolute form of path, resolving symlinks.
func AbsPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("filesystem: abs path %q: %w", path, err)
	}
	return abs, nil
}
