package connect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Drafteame/draft/internal/pkg/dirs"
)

// checkPortFree returns an error if the given local TCP port is already in use.
func checkPortFree(port int) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("local port %d is already in use", port)
	}

	_ = ln.Close()

	return nil
}

// launchTunnel starts the SSM port-forwarding process in the background and
// waits until the local port is accepting connections (tunnel ready) or the
// process dies (tunnel failed).
//
// The subprocess runs in its own process group (Setpgid=true) so that:
//   - the shell does not hang after `draft` exits
//   - stop can kill the entire group with Kill(-pid, SIGKILL), which also
//     terminates the Session Manager Plugin child process
//
// Stdout and stderr of `aws ssm` / `session-manager-plugin` are written to
// logPath (a file, not a pipe, so the tunnel keeps working after `draft`
// exits). If startup fails, the tail of that log is included in the error so
// problems such as `no such host` are visible instead of a generic timeout.
//
// Returns the PID of the launched process (which equals its PGID).
func launchTunnel(bastion BastionConfig, host string, remotePort, localPort int, logPath string) (int, error) {
	params, err := buildSSMParams(host, remotePort, localPort)
	if err != nil {
		return 0, fmt.Errorf("failed to build SSM parameters: %w", err)
	}

	if err := dirs.Create(filepath.Dir(logPath)); err != nil {
		return 0, fmt.Errorf("failed to create tunnel log directory: %w", err)
	}

	logFile, err := os.Create(logPath)
	if err != nil {
		return 0, fmt.Errorf("failed to create tunnel log file: %w", err)
	}
	// The child inherits its own copy of the descriptor; ours can be closed.
	defer func() { _ = logFile.Close() }()

	cmd := exec.Command("aws", "ssm", "start-session",
		"--target", bastion.Target,
		"--document-name", "AWS-StartPortForwardingSessionToRemoteHost",
		"--parameters", params,
		"--region", bastion.Region,
		"--profile", bastion.Profile,
	)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	// New process group: detaches from the shell's group so the terminal is
	// not held waiting, and allows group-kill on stop.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("failed to start tunnel: %w", err)
	}

	// Reap the process in the background so an early exit is detected
	// (an un-reaped zombie still answers Signal(0)).
	exited := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(exited)
	}()

	// Poll until the local port is accepting connections (tunnel established)
	// or the process dies (tunnel failed). Timeout: 10 seconds.
	pid, err := waitForTunnel(cmd, exited, logPath, localPort, 10*time.Second)
	if err != nil {
		// Best-effort: kill any partial process that may still be running.
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		return 0, err
	}

	return pid, nil
}

// waitForTunnel polls until:
//   - the local port accepts a TCP connection → success
//   - the process exits                       → error with session output
//   - the timeout is reached                  → kill process and error with session output
func waitForTunnel(cmd *exec.Cmd, exited <-chan struct{}, logPath string, localPort int, timeout time.Duration) (int, error) {
	deadline := time.Now().Add(timeout)
	addr := fmt.Sprintf("127.0.0.1:%d", localPort)

	for time.Now().Before(deadline) {
		select {
		case <-exited:
			return 0, tunnelError(logPath, "tunnel process exited during startup")
		case <-time.After(300 * time.Millisecond):
		}

		// Check if the local port is now accepting connections.
		conn, dialErr := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if dialErr == nil {
			_ = conn.Close()
			return cmd.Process.Pid, nil
		}
	}

	// Timeout reached — port never became available.
	// Kill the process (it's running but not working) and report failure.
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)

	return 0, tunnelError(logPath,
		fmt.Sprintf("tunnel timed out: local port %d never became available (check bastion connectivity and AWS credentials)", localPort))
}

// maxLogTailBytes caps how much session output is embedded in an error.
const maxLogTailBytes = 2048

// tunnelError builds a startup error from msg plus the tail of the session log, if any.
func tunnelError(logPath, msg string) error {
	out := readLogTail(logPath, maxLogTailBytes)
	if out == "" {
		return fmt.Errorf("%s", msg)
	}

	return fmt.Errorf("%s\nsession output (%s):\n%s", msg, logPath, out)
}

func readLogTail(path string, maxBytes int) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	if len(data) > maxBytes {
		data = data[len(data)-maxBytes:]
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			data = data[i+1:]
		}
	}

	return strings.TrimSpace(string(data))
}

func buildSSMParams(host string, remotePort, localPort int) (string, error) {
	params := map[string][]string{
		"host":            {host},
		"portNumber":      {fmt.Sprintf("%d", remotePort)},
		"localPortNumber": {fmt.Sprintf("%d", localPort)},
	}

	b, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
