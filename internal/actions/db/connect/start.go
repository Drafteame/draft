package connect

import (
	"fmt"
	"time"

	"github.com/Drafteame/draft/internal/pkg/log"
)

// StartInput holds parameters for starting a tunnel.
type StartInput struct {
	DBType    string
	Name      string // "{service}-{env}", e.g. "turbo-dev"
	LocalPort int    // 0 means use the default computed from config
}

// Start launches a new SSM tunnel for the given type+name.
func Start(input StartInput) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	resolved, err := cfg.ResolveConnection(input.DBType, input.Name)
	if err != nil {
		return err
	}

	state, err := loadRuntimeState()
	if err != nil {
		return err
	}

	key := runtimeKey(input.DBType, input.Name)

	if existing, ok := state.Connections[key]; ok {
		if isProcessAlive(existing.PID) {
			return fmt.Errorf("connection %s/%s is already running (PID: %d, port: localhost:%d)",
				input.DBType, input.Name, existing.PID, existing.LocalPort)
		}

		// Dead entry — clean it up and continue
		delete(state.Connections, key)
	}

	localPort := resolved.LocalPort
	if input.LocalPort > 0 {
		localPort = input.LocalPort
	}

	if err := checkPortFree(localPort); err != nil {
		return err
	}

	pid, err := launchTunnel(resolved.Bastion, resolved.Host, resolved.RemotePort, localPort)
	if err != nil {
		return err
	}

	state.Connections[key] = RuntimeEntry{
		DBType:     input.DBType,
		Name:       resolved.Name,
		Env:        resolved.Env,
		PID:        pid,
		LocalPort:  localPort,
		RemotePort: resolved.RemotePort,
		Host:       resolved.Host,
		StartedAt:  time.Now(),
	}

	if err := saveRuntimeState(state); err != nil {
		log.Warnf("could not persist runtime state: %s", err)
	}

	log.Successf("Tunnel started: %s/%s", input.DBType, input.Name)
	log.Infof("  Env     : %s", resolved.Env)
	log.Infof("  PID     : %d", pid)
	log.Infof("  Local   : localhost:%d", localPort)
	log.Infof("  Remote  : %s:%d", resolved.Host, resolved.RemotePort)
	log.Infof("  Stop    : draft dbconnect stop %s %s", input.DBType, input.Name)

	return nil
}
