package connect

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// deadPID returns the PID of a process that has already exited and been reaped.
func deadPID(t *testing.T) int {
	t.Helper()

	cmd := exec.Command("true")
	require.NoError(t, cmd.Run())

	return cmd.Process.Pid
}

func TestStatus_RemovesPgpassEntryOfDeadTunnel(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	pgpass := filepath.Join(home, "pgpass")
	t.Setenv("PGPASSFILE", pgpass)

	require.NoError(t, os.WriteFile(pgpass, []byte(
		"localhost:15432:*:dev@draftea.com:dead-token\n"+
			"localhost:15433:*:dev@draftea.com:live-token\n"+
			"other-host:5432:*:someone:secret\n"), 0600))

	deadKey := runtimeKey("postgres", "turbo-dev")
	liveKey := runtimeKey("postgres", "users-dev")

	require.NoError(t, saveRuntimeState(RuntimeState{Connections: map[string]RuntimeEntry{
		deadKey: {DBType: "postgres", Name: "turbo-dev", PID: deadPID(t), LocalPort: 15432, IAMUser: "dev@draftea.com"},
		liveKey: {DBType: "postgres", Name: "users-dev", PID: os.Getpid(), LocalPort: 15433, IAMUser: "dev@draftea.com"},
	}}))

	require.NoError(t, Status())

	data, err := os.ReadFile(pgpass)
	require.NoError(t, err)
	assert.Equal(t,
		"localhost:15433:*:dev@draftea.com:live-token\nother-host:5432:*:someone:secret\n",
		string(data),
		"only the dead tunnel's entry must be removed")

	saved, err := loadRuntimeState()
	require.NoError(t, err)
	assert.NotContains(t, saved.Connections, deadKey)
	assert.Contains(t, saved.Connections, liveKey)
}
