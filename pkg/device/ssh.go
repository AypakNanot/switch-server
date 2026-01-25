package device

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHAdapter implements ProtocolAdapter for SSH protocol
type SSHAdapter struct {
	client    *ssh.Client
	session   *ssh.Session
	connected bool
}

// NewSSHAdapter creates a new SSH adapter
func NewSSHAdapter() *SSHAdapter {
	return &SSHAdapter{}
}

// Connect establishes an SSH connection
func (a *SSHAdapter) Connect(ctx context.Context, config *ConnectionConfig) error {
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(config.Timeout) * time.Second,
	}

	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return NewConnectionError(err)
	}

	a.client = client
	a.connected = true
	return nil
}

// Disconnect closes the SSH connection
func (a *SSHAdapter) Disconnect(ctx context.Context) error {
	if a.session != nil {
		a.session.Close()
		a.session = nil
	}
	if a.client != nil {
		err := a.client.Close()
		a.client = nil
		a.connected = false
		return err
	}
	return nil
}

// ExecuteCommand executes a single command
func (a *SSHAdapter) ExecuteCommand(ctx context.Context, cmd string) (*CommandResult, error) {
	if !a.connected {
		return nil, NewConnectionError(fmt.Errorf("not connected"))
	}

	session, err := a.client.NewSession()
	if err != nil {
		return nil, NewCommandFailedError(err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, NewCommandFailedError(err)
	}

	startTime := time.Now()
	output, err := session.Output(cmd)
	duration := time.Since(startTime)

	result := &CommandResult{
		Command:   cmd,
		Output:    string(output),
		Duration:  duration.Milliseconds(),
		Success:   err == nil,
		Timestamp: time.Now().Unix(),
	}

	if err != nil {
		result.Error = err.Error()
		if ctx.Err() == context.DeadlineExceeded {
			return nil, NewCommandTimeoutError()
		}
		return result, NewCommandFailedError(err)
	}

	return result, nil
}

// IsConnected returns whether the SSH connection is active
func (a *SSHAdapter) IsConnected() bool {
	return a.connected && a.client != nil
}

// ProtocolType returns the protocol type
func (a *SSHAdapter) ProtocolType() ProtocolType {
	return ProtocolSSH
}

// NewSSHAdapterFunc creates a new SSH adapter (factory function)
func NewSSHAdapterFunc() ProtocolAdapter {
	return NewSSHAdapter()
}
