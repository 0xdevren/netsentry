package util

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// SHA256File computes the SHA-256 hash of the file at the given path and
// returns the hex-encoded digest string.
func SHA256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("hashing: open %q: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hashing: compute %q: %w", path, err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// SHA256Bytes computes the SHA-256 hash of the provided byte slice and returns
// the hex-encoded digest string.
func SHA256Bytes(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h)
}
