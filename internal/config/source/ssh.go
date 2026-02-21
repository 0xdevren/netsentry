package source

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHOptions holds the authentication and connection parameters for SSH sources.
type SSHOptions struct {
	// Host is the target device hostname or IP address.
	Host string
	// Port is the SSH port (defaults to 22 if zero).
	Port int
	// User is the SSH username.
	User string
	// PrivateKeyPath is the path to the PEM-encoded private key file.
	PrivateKeyPath string
	// Password is used for password authentication when PrivateKeyPath is empty.
	Password string
	// Command is the CLI command to execute to retrieve the configuration.
	// Defaults to "show running-config".
	Command string
	// Timeout is the connection and command execution timeout.
	Timeout time.Duration
}

// SSHSource retrieves device configuration via SSH command execution.
type SSHSource struct{}

// NewSSHSource constructs an SSHSource.
func NewSSHSource() *SSHSource {
	return &SSHSource{}
}

// Load establishes an SSH session to the target device, executes the
// configuration retrieval command, and returns the output.
func (s *SSHSource) Load(ctx context.Context, req LoadRequest) ([]byte, error) {
	opts := req.SSHOptions
	if opts == nil {
		return nil, fmt.Errorf("ssh source: SSHOptions are required")
	}

	port := opts.Port
	if port == 0 {
		port = 22
	}
	command := opts.Command
	if command == "" {
		command = "show running-config"
	}
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	authMethods, err := s.buildAuthMethods(opts)
	if err != nil {
		return nil, fmt.Errorf("ssh source: auth: %w", err)
	}

	cfg := &ssh.ClientConfig{
		User:            opts.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec // device fingerprints managed separately
		Timeout:         timeout,
	}

	addr := net.JoinHostPort(opts.Host, fmt.Sprintf("%d", port))
	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return nil, fmt.Errorf("ssh source: dial %s: %w", addr, err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("ssh source: new session: %w", err)
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf

	// Respect context cancellation during command execution.
	done := make(chan error, 1)
	go func() { done <- session.Run(command) }()

	select {
	case <-ctx.Done():
		_ = session.Signal(ssh.SIGTERM)
		return nil, fmt.Errorf("ssh source: context cancelled: %w", ctx.Err())
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("ssh source: run %q: %w", command, err)
		}
	}

	return buf.Bytes(), nil
}

// buildAuthMethods returns the appropriate SSH authentication methods.
func (s *SSHSource) buildAuthMethods(opts *SSHOptions) ([]ssh.AuthMethod, error) {
	if opts.PrivateKeyPath != "" {
		key, err := os.ReadFile(opts.PrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("read private key %q: %w", opts.PrivateKeyPath, err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("parse private key: %w", err)
		}
		return []ssh.AuthMethod{ssh.PublicKeys(signer)}, nil
	}
	if opts.Password != "" {
		return []ssh.AuthMethod{ssh.Password(opts.Password)}, nil
	}
	return nil, fmt.Errorf("either PrivateKeyPath or Password must be provided")
}
