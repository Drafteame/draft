package connect

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/Drafteame/draft/internal/pkg/files"
)

const stateFileName = "dbconnect.state.json"

func stateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	return filepath.Join(home, ".draft", stateFileName), nil
}

func loadRuntimeState() (RuntimeState, error) {
	path, err := stateFilePath()
	if err != nil {
		return RuntimeState{}, err
	}

	if !files.Exists(path) {
		return RuntimeState{Connections: map[string]RuntimeEntry{}}, nil
	}

	data, err := files.Read(path)
	if err != nil {
		return RuntimeState{}, fmt.Errorf("failed to read state file: %w", err)
	}

	var state RuntimeState
	if err := json.Unmarshal(data, &state); err != nil {
		return RuntimeState{}, fmt.Errorf("failed to parse state file: %w", err)
	}

	if state.Connections == nil {
		state.Connections = map[string]RuntimeEntry{}
	}

	return state, nil
}

func saveRuntimeState(state RuntimeState) error {
	path, err := stateFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	return files.Create(path, data)
}

// isProcessAlive returns true if the process with the given PID is still running.
func isProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	return proc.Signal(syscall.Signal(0)) == nil
}
