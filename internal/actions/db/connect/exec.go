package connect

import (
	"bytes"
	"encoding/json"
	"errors"
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

// tunnelSettleTimeout bounds how long start keeps a probe connection open
// through the tunnel waiting for proof that the bastion reached the remote
// host. A refused remote port closes the forwarded connection at once and a
// failed DNS lookup usually ends the session a couple of seconds after the
// local port opens, but an unknown name can also just hang. With a probe, a
// reply is the proof and silence for the whole window is a failure; without
// one, silence is all a TLS-only endpoint offers and it is taken as success.
//
// A variable so tests can shorten it.
var tunnelSettleTimeout = 3 * time.Second

// postgresSSLRequest is the PostgreSQL wire-protocol SSLRequest message. Every
// PostgreSQL endpoint (cluster or RDS Proxy) answers it right away with one
// byte ('S' or 'N'), so it proves the tunnel ends at a live server.
var postgresSSLRequest = []byte{0, 0, 0, 8, 4, 210, 22, 47}

// tunnelProbe returns the bytes sent through a freshly opened tunnel to make
// the remote side answer. Engines without a harmless plaintext greeting (TLS
// only redis/docdb endpoints) get nil and are confirmed by silence instead.
func tunnelProbe(dbType string) []byte {
	if dbType == iamDBType {
		return postgresSSLRequest
	}

	return nil
}

// launchTunnel starts the SSM port-forwarding process in the background and
// waits until a connection through the tunnel is confirmed (see
// confirmTunnel) or the process dies (tunnel failed).
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
func launchTunnel(bastion BastionConfig, host string, remotePort, localPort int, logPath string, probe []byte) (int, error) {
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

	// Poll until the local port is accepting connections and the remote side
	// is confirmed, or the process dies (tunnel failed). Timeout: 10 seconds.
	pid, err := waitForTunnel(cmd, exited, logPath, localPort, probe, 10*time.Second)
	if err != nil {
		// Best-effort: kill any partial process that may still be running.
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		return 0, err
	}

	return pid, nil
}

// waitForTunnel polls until:
//   - the local port accepts a TCP connection → confirmTunnel decides
//   - the process exits                       → error with session output
//   - the timeout is reached                  → kill process and error with session output
func waitForTunnel(cmd *exec.Cmd, exited <-chan struct{}, logPath string, localPort int, probe []byte, timeout time.Duration) (int, error) {
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
			return confirmTunnel(cmd, exited, logPath, conn, probe)
		}
	}

	// Timeout reached — port never became available.
	// Kill the process (it's running but not working) and report failure.
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)

	return 0, tunnelError(logPath,
		fmt.Sprintf("tunnel timed out: local port %d never became available (check bastion connectivity and AWS credentials)", localPort))
}

// confirmTunnel checks that a connection accepted on the local port really
// reaches the remote host. The Session Manager plugin opens the local port
// before the bastion has connected anywhere, so accepting alone proves
// nothing: a wrong host name or a closed remote port only shows up afterwards,
// as the session exiting or the forwarded connection being closed.
//
// The probe (if any) is written and then one byte is awaited for up to
// tunnelSettleTimeout. A reply means the remote answered; EOF or the process
// exiting means the bastion could not reach it; silence is a failure when a
// probe was sent and a success otherwise. The connection is closed before
// returning so the tunnel starts clean for the real client.
func confirmTunnel(cmd *exec.Cmd, exited <-chan struct{}, logPath string, conn net.Conn, probe []byte) (int, error) {
	defer func() { _ = conn.Close() }()

	readDone := make(chan error, 1)

	go func() {
		if len(probe) > 0 {
			if _, err := conn.Write(probe); err != nil {
				readDone <- err
				return
			}
		}

		_ = conn.SetReadDeadline(time.Now().Add(tunnelSettleTimeout))

		_, err := conn.Read(make([]byte, 1))
		readDone <- err
	}()

	select {
	case <-exited:
		return 0, tunnelError(logPath, "tunnel process exited right after opening the local port")
	case err := <-readDone:
		if err == nil {
			return cmd.Process.Pid, nil
		}

		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			if len(probe) == 0 {
				return cmd.Process.Pid, nil
			}

			return 0, tunnelError(logPath, fmt.Sprintf(
				"tunnel is open but the remote host did not answer within %s (does it exist and can the bastion reach it?)",
				tunnelSettleTimeout))
		}

		// The forwarded connection was closed: the bastion could not reach the
		// remote host. Give the session a moment to log the reason or exit so
		// the error carries it.
		select {
		case <-exited:
		case <-time.After(time.Second):
		}

		return 0, tunnelError(logPath, "tunnel is open but the bastion could not connect to the remote host")
	}
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
