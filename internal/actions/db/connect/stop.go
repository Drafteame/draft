package connect

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/Drafteame/draft/internal/pkg/log"
)

// StopInput holds parameters for stopping one or more tunnels.
type StopInput struct {
	DBType string // empty = all types
	Name   string // empty = all within type
}

// Stop terminates active tunnels based on the input:
//   - DBType="" → stop all active connections
//   - DBType set, Name="" → stop all active connections of that type
//   - DBType set, Name set → stop that specific connection
func Stop(input StopInput) error {
	state, err := loadRuntimeState()
	if err != nil {
		return err
	}

	if input.DBType == "" {
		return stopAll(state)
	}

	if input.Name == "" {
		return stopByType(state, input.DBType)
	}

	return stopOne(state, input.DBType, input.Name)
}

func stopAll(state RuntimeState) error {
	if len(state.Connections) == 0 {
		log.Info("No active connections to stop.")
		return nil
	}

	var errs []string

	for key, entry := range state.Connections {
		if err := killEntry(state, key, entry); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if err := saveRuntimeState(state); err != nil {
		log.Warnf("could not update runtime state: %s", err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("some stops failed:\n%s", strings.Join(errs, "\n"))
	}

	return nil
}

func stopByType(state RuntimeState, dbType string) error {
	var keys []string

	for key, entry := range state.Connections {
		if entry.DBType == dbType {
			keys = append(keys, key)
		}
	}

	if len(keys) == 0 {
		return fmt.Errorf("no active %s connections found", dbType)
	}

	var errs []string

	for _, key := range keys {
		if err := killEntry(state, key, state.Connections[key]); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if err := saveRuntimeState(state); err != nil {
		log.Warnf("could not update runtime state: %s", err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("some stops failed:\n%s", strings.Join(errs, "\n"))
	}

	return nil
}

func stopOne(state RuntimeState, dbType, name string) error {
	key := runtimeKey(dbType, name)

	entry, ok := state.Connections[key]
	if !ok {
		return fmt.Errorf("no active connection found for %s/%s", dbType, name)
	}

	if err := killEntry(state, key, entry); err != nil {
		return err
	}

	if err := saveRuntimeState(state); err != nil {
		log.Warnf("could not update runtime state: %s", err)
	}

	return nil
}

// killEntry terminates the process for a single connection entry and removes it from state.
func killEntry(state RuntimeState, key string, entry RuntimeEntry) error {
	if !isProcessAlive(entry.PID) {
		delete(state.Connections, key)
		log.Warnf("connection %s/%s was not running (stale state cleaned up)", entry.DBType, entry.Name)

		return nil
	}

	// Kill the entire process group (PGID == PID because of Setpgid=true at launch).
	if err := syscall.Kill(-entry.PID, syscall.SIGKILL); err != nil {
		return fmt.Errorf("failed to stop %s/%s (PID %d): %w", entry.DBType, entry.Name, entry.PID, err)
	}

	delete(state.Connections, key)
	log.Successf("Tunnel stopped: %s/%s (PID: %d)", entry.DBType, entry.Name, entry.PID)

	return nil
}
