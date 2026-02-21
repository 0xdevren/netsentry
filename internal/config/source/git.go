package source

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// GitOptions configures a Git repository config source.
type GitOptions struct {
	// RepoURL is the remote Git repository URL.
	RepoURL string
	// Branch is the branch to check out (defaults to "main").
	Branch string
	// FilePath is the path within the repository to the config file.
	FilePath string
	// CloneDir is the local directory used for the clone (defaults to a temp dir).
	CloneDir string
}

// gitSource fetches configurations from a Git repository.
type gitSource struct{}

// NewGitSource constructs a gitSource.
func NewGitSource() ConfigSource {
	return &gitSource{}
}

// Load clones or fetches the repository and reads the specified file.
func (g *gitSource) Load(ctx context.Context, req LoadRequest) ([]byte, error) {
	opts := req.GitOptions
	if opts == nil {
		return nil, fmt.Errorf("git source: GitOptions are required")
	}
	if opts.RepoURL == "" {
		return nil, fmt.Errorf("git source: RepoURL is required")
	}
	if opts.FilePath == "" {
		return nil, fmt.Errorf("git source: FilePath is required")
	}

	branch := opts.Branch
	if branch == "" {
		branch = "main"
	}

	cloneDir := opts.CloneDir
	if cloneDir == "" {
		tmp, err := os.MkdirTemp("", "netsentry-git-*")
		if err != nil {
			return nil, fmt.Errorf("git source: create temp dir: %w", err)
		}
		defer os.RemoveAll(tmp)
		cloneDir = tmp
	}

	args := []string{"clone", "--depth=1", "--branch", branch, opts.RepoURL, cloneDir}
	cmd := exec.CommandContext(ctx, "git", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git source: clone: %w: %s", err, string(out))
	}

	targetPath := filepath.Join(cloneDir, opts.FilePath)
	data, err := os.ReadFile(targetPath)
	if err != nil {
		return nil, fmt.Errorf("git source: read %q: %w", targetPath, err)
	}
	return data, nil
}
