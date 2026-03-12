package connect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"syscall"
	"time"
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
// Returns the PID of the launched process (which equals its PGID).
func launchTunnel(bastion BastionConfig, host string, remotePort, localPort int) (int, error) {
	params, err := buildSSMParams(host, remotePort, localPort)
	if err != nil {
		return 0, fmt.Errorf("failed to build SSM parameters: %w", err)
	}

	var stderrBuf bytes.Buffer

	cmd := exec.Command("aws", "ssm", "start-session",
		"--target", bastion.Target,
		"--document-name", "AWS-StartPortForwardingSessionToRemoteHost",
		"--parameters", params,
		"--region", bastion.Region,
		"--profile", bastion.Profile,
	)
	cmd.Stdout = nil
	cmd.Stderr = &stderrBuf
	// New process group: detaches from the shell's group so the terminal is
	// not held waiting, and allows group-kill on stop.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("failed to start tunnel: %w", err)
	}

	// Poll until the local port is accepting connections (tunnel established)
	// or the process dies (tunnel failed). Timeout: 10 seconds.
	pid, err := waitForTunnel(cmd, &stderrBuf, localPort, 10*time.Second)
	if err != nil {
		// Best-effort: kill any partial process that may still be running.
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		return 0, err
	}

	return pid, nil
}

// waitForTunnel polls until:
//   - the local port accepts a TCP connection → success
//   - the process dies                        → error with stderr
//   - the timeout is reached                  → kill process and error
func waitForTunnel(cmd *exec.Cmd, stderrBuf *bytes.Buffer, localPort int, timeout time.Duration) (int, error) {
	deadline := time.Now().Add(timeout)
	addr := fmt.Sprintf("127.0.0.1:%d", localPort)

	for time.Now().Before(deadline) {
		time.Sleep(300 * time.Millisecond)

		// Check if the process died.
		if cmd.Process.Signal(syscall.Signal(0)) != nil {
			return 0, tunnelError(stderrBuf, "tunnel process died during startup")
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

	return 0, tunnelError(stderrBuf,
		fmt.Sprintf("tunnel timed out: local port %d never became available (check bastion connectivity and AWS credentials)", localPort))
}

func tunnelError(stderrBuf *bytes.Buffer, fallback string) error {
	msg := strings.TrimSpace(stderrBuf.String())
	if msg == "" {
		msg = fallback
	}

	return fmt.Errorf("%s", msg)
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
