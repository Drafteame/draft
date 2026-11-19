package connect

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadLogTail(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name     string
		content  *string
		maxBytes int
		want     string
	}{
		{
			name:     "missing file",
			maxBytes: 100,
			want:     "",
		},
		{
			name:     "short content is returned trimmed",
			content:  ptr("  line1\nline2\n\n"),
			maxBytes: 100,
			want:     "line1\nline2",
		},
		{
			name:     "long content keeps only whole trailing lines",
			content:  ptr("aaaaaaaaaa\nbbbbbbbbbb\ncccccccccc\n"),
			maxBytes: 15,
			want:     "cccccccccc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, strings.ReplaceAll(tt.name, " ", "_")+".log")

			if tt.content != nil {
				require.NoError(t, os.WriteFile(path, []byte(*tt.content), 0o600))
			}

			assert.Equal(t, tt.want, readLogTail(path, tt.maxBytes))
		})
	}
}

func TestWaitForTunnelDetectsEarlyExit(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "tunnel.log")

	logFile, err := os.Create(logPath)
	require.NoError(t, err)

	defer func() { _ = logFile.Close() }()

	cmd := exec.Command("sh", "-c", "echo 'dial tcp: lookup api-cache-dev: no such host' >&2; exit 1")
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	require.NoError(t, cmd.Start())

	exited := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(exited)
	}()

	start := time.Now()
	_, err = waitForTunnel(cmd, exited, logPath, freePort(t), nil, 10*time.Second)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "tunnel process exited during startup")
	assert.Contains(t, err.Error(), "no such host")
	assert.Less(t, time.Since(start), 5*time.Second, "early exit should not wait for the timeout")
}

// fakeTunnel stands in for the Session Manager plugin: a long-lived process
// whose session log already holds logContent, plus a listener on a free local
// port whose accepted connections are handed to handle.
func fakeTunnel(t *testing.T, logContent string, handle func(net.Conn)) (*exec.Cmd, chan struct{}, string, int) {
	t.Helper()

	logPath := filepath.Join(t.TempDir(), "tunnel.log")
	require.NoError(t, os.WriteFile(logPath, []byte(logContent), 0o600))

	cmd := exec.Command("sleep", "30")
	require.NoError(t, cmd.Start())

	exited := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(exited)
	}()

	t.Cleanup(func() { _ = cmd.Process.Kill() })

	port := freePort(t)

	ln, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)))
	require.NoError(t, err)

	t.Cleanup(func() { _ = ln.Close() })

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}

			go handle(conn)
		}
	}()

	return cmd, exited, logPath, port
}

func TestWaitForTunnelFailsWhenRemoteClosesTheConnection(t *testing.T) {
	cmd, exited, logPath, port := fakeTunnel(t,
		"Waiting for connections...\nConnection to destination port failed, check SSM Agent logs.\n",
		func(conn net.Conn) { _ = conn.Close() })

	start := time.Now()
	_, err := waitForTunnel(cmd, exited, logPath, port, postgresSSLRequest, 10*time.Second)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not connect to the remote host")
	assert.Contains(t, err.Error(), "Connection to destination port failed")
	assert.Less(t, time.Since(start), 5*time.Second)
}

func TestWaitForTunnelSucceedsWhenRemoteAnswersProbe(t *testing.T) {
	got := make(chan []byte, 1)

	cmd, exited, logPath, port := fakeTunnel(t, "Waiting for connections...\n", func(conn net.Conn) {
		defer func() { _ = conn.Close() }()

		buf := make([]byte, len(postgresSSLRequest))
		if _, err := conn.Read(buf); err == nil {
			got <- buf
			_, _ = conn.Write([]byte("N"))
		}
	})

	start := time.Now()
	pid, err := waitForTunnel(cmd, exited, logPath, port, postgresSSLRequest, 10*time.Second)

	require.NoError(t, err)
	assert.Equal(t, cmd.Process.Pid, pid)
	assert.Equal(t, postgresSSLRequest, <-got, "the probe must be the SSLRequest message")
	assert.Less(t, time.Since(start), tunnelSettleTimeout, "a reply must not wait for the settle timeout")
}

// silentTunnel is a fake tunnel whose remote side never answers, with the
// settle timeout shortened so the test does not wait the real three seconds.
func silentTunnel(t *testing.T) (*exec.Cmd, chan struct{}, string, int) {
	t.Helper()

	orig := tunnelSettleTimeout
	tunnelSettleTimeout = 300 * time.Millisecond

	t.Cleanup(func() { tunnelSettleTimeout = orig })

	hold := make(chan struct{})
	t.Cleanup(func() { close(hold) })

	return fakeTunnel(t, "Waiting for connections...\n", func(conn net.Conn) {
		<-hold
		_ = conn.Close()
	})
}

func TestWaitForTunnelSucceedsWhenRemoteStaysSilentWithoutProbe(t *testing.T) {
	cmd, exited, logPath, port := silentTunnel(t)

	start := time.Now()
	pid, err := waitForTunnel(cmd, exited, logPath, port, nil, 10*time.Second)

	require.NoError(t, err)
	assert.Equal(t, cmd.Process.Pid, pid)
	assert.GreaterOrEqual(t, time.Since(start), tunnelSettleTimeout, "silence is only trusted after the settle timeout")
}

func TestWaitForTunnelFailsWhenRemoteIgnoresProbe(t *testing.T) {
	cmd, exited, logPath, port := silentTunnel(t)

	_, err := waitForTunnel(cmd, exited, logPath, port, postgresSSLRequest, 10*time.Second)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "did not answer within 300ms")
}

func freePort(t *testing.T) int {
	t.Helper()

	for port := 59000; port < 59100; port++ {
		if checkPortFree(port) == nil {
			return port
		}
	}

	t.Fatal("no free port found")

	return 0
}

func ptr(s string) *string { return &s }
