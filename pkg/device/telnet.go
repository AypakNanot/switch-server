package device

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

const telnetTimeout = 10 * time.Second

// TelnetAdapter implements ProtocolAdapter for Telnet protocol
type TelnetAdapter struct {
	conn      net.Conn
	reader    *bufio.Reader
	writer    io.Writer
	connected bool
	prompt    string
}

// NewTelnetAdapter creates a new Telnet adapter
func NewTelnetAdapter() *TelnetAdapter {
	return &TelnetAdapter{
		prompt: "#|>|\\$",
	}
}

// Connect establishes a Telnet connection
func (a *TelnetAdapter) Connect(ctx context.Context, config *ConnectionConfig) error {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)

	dialer := net.Dialer{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return NewConnectionError(err)
	}

	a.conn = conn
	a.reader = bufio.NewReader(conn)
	a.writer = conn
	a.connected = true

	if config.Username != "" {
		if err := a.login(ctx, config); err != nil {
			a.Disconnect(ctx)
			return err
		}
	}

	return nil
}

func (a *TelnetAdapter) login(ctx context.Context, config *ConnectionConfig) error {
	if err := a.waitFor(ctx, "ogin:"); err != nil {
		return NewAuthError(err)
	}

	if _, err := fmt.Fprintf(a.writer, "%s\r\n", config.Username); err != nil {
		return NewAuthError(err)
	}

	if err := a.waitFor(ctx, "assword:"); err != nil {
		return NewAuthError(err)
	}

	if _, err := fmt.Fprintf(a.writer, "%s\r\n", config.Password); err != nil {
		return NewAuthError(err)
	}

	if err := a.waitForPrompt(ctx); err != nil {
		return NewAuthError(err)
	}

	return nil
}

// Disconnect closes the Telnet connection
func (a *TelnetAdapter) Disconnect(ctx context.Context) error {
	if a.conn != nil {
		err := a.conn.Close()
		a.conn = nil
		a.reader = nil
		a.writer = nil
		a.connected = false
		return err
	}
	return nil
}

// ExecuteCommand executes a single command
func (a *TelnetAdapter) ExecuteCommand(ctx context.Context, cmd string) (*CommandResult, error) {
	if !a.connected {
		return nil, NewConnectionError(fmt.Errorf("not connected"))
	}

	if _, err := fmt.Fprintf(a.writer, "%s\r\n", cmd); err != nil {
		return nil, NewCommandFailedError(err)
	}

	if _, err := a.reader.ReadString('\n'); err != nil {
		return nil, NewCommandFailedError(err)
	}

	startTime := time.Now()
	output, err := a.readUntilPrompt(ctx)
	duration := time.Since(startTime)

	result := &CommandResult{
		Command:   cmd,
		Output:    output,
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

func (a *TelnetAdapter) waitFor(ctx context.Context, waitStr string) error {
	timeout := time.After(telnetTimeout)
	buf := make([]byte, 1)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timeout waiting for '%s'", waitStr)
		default:
			n, err := a.reader.Read(buf)
			if err != nil {
				return err
			}
			if n > 0 && strings.Contains(string(buf), waitStr) {
				return nil
			}
		}
	}
}

func (a *TelnetAdapter) waitForPrompt(ctx context.Context) error {
	timeout := time.After(telnetTimeout)
	buf := make([]byte, 4096)
	output := ""

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timeout waiting for prompt")
		default:
			a.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, err := a.reader.Read(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				return err
			}
			if n > 0 {
				output += string(buf[:n])
				if strings.Contains(output, "#") || strings.Contains(output, ">") || strings.Contains(output, "$") {
					return nil
				}
			}
		}
	}
}

func (a *TelnetAdapter) readUntilPrompt(ctx context.Context) (string, error) {
	timeout := time.After(telnetTimeout)
	buf := make([]byte, 4096)
	var output strings.Builder

	for {
		select {
		case <-ctx.Done():
			return output.String(), ctx.Err()
		case <-timeout:
			return output.String(), fmt.Errorf("timeout reading output")
		default:
			a.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, err := a.reader.Read(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					out := output.String()
					if strings.Contains(out, "#") || strings.Contains(out, ">") || strings.Contains(out, "$") {
						lines := strings.Split(out, "\n")
						if len(lines) > 0 {
							lastLine := strings.TrimSpace(lines[len(lines)-1])
							if len(lastLine) <= 3 && (strings.HasSuffix(lastLine, "#") || strings.HasSuffix(lastLine, ">") || strings.HasSuffix(lastLine, "$")) {
								out = strings.Join(lines[:len(lines)-1], "\n")
							}
						}
						return out, nil
					}
					continue
				}
				return output.String(), err
			}
			if n > 0 {
				output.Write(buf[:n])
				out := output.String()
				if strings.Contains(out, "#") || strings.Contains(out, ">") || strings.Contains(out, "$") {
					lines := strings.Split(out, "\n")
					if len(lines) > 0 {
						lastLine := strings.TrimSpace(lines[len(lines)-1])
						if len(lastLine) <= 3 && (strings.HasSuffix(lastLine, "#") || strings.HasSuffix(lastLine, ">") || strings.HasSuffix(lastLine, "$")) {
							out = strings.Join(lines[:len(lines)-1], "\n")
						}
					}
					return out, nil
				}
			}
		}
	}
}

// IsConnected returns whether the Telnet connection is active
func (a *TelnetAdapter) IsConnected() bool {
	return a.connected && a.conn != nil
}

// ProtocolType returns the protocol type
func (a *TelnetAdapter) ProtocolType() ProtocolType {
	return ProtocolTelnet
}

// NewTelnetAdapterFunc creates a new Telnet adapter (factory function)
func NewTelnetAdapterFunc() ProtocolAdapter {
	return NewTelnetAdapter()
}
