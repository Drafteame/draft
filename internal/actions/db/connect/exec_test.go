package connect

import (
	"os"
	"os/exec"
	"path/filepath"
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
	_, err = waitForTunnel(cmd, exited, logPath, freePort(t), 10*time.Second)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "tunnel process exited during startup")
	assert.Contains(t, err.Error(), "no such host")
	assert.Less(t, time.Since(start), 5*time.Second, "early exit should not wait for the timeout")
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
