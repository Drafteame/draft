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

// Start launches a new SSM tunnel for the given type+name. Postgres tunnels
// also get an RDS IAM auth token written to .pgpass; if the tunnel is already
// running, only the token is refreshed.
func Start(input StartInput) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	resolved, err := cfg.ResolveConnection(input.DBType, input.Name)
	if err != nil {
		return err
	}

	withToken := resolved.DBType == iamDBType

	// Resolve the IAM user before launching anything so a missing value does
	// not leave a tunnel without its token behind.
	iamUser := ""

	if withToken {
		iamUser, err = resolveIAMUser()
		if err != nil {
			return err
		}
	}

	state, err := loadRuntimeState()
	if err != nil {
		return err
	}

	key := runtimeKey(input.DBType, input.Name)

	if existing, ok := state.Connections[key]; ok {
		if isProcessAlive(existing.PID) {
			if withToken {
				return refreshToken(state, key, existing, resolved, iamUser)
			}

			return fmt.Errorf("connection %s/%s is already running (PID: %d, port: localhost:%d)",
				input.DBType, input.Name, existing.PID, existing.LocalPort)
		}

		// Dead entry — clean it up and continue. Its token may sit on another
		// local port (--port override), so the new entry would not replace it.
		cleanupIAMToken(existing.LocalPort, existing.IAMUser)
		delete(state.Connections, key)
	}

	localPort := resolved.LocalPort
	if input.LocalPort > 0 {
		localPort = input.LocalPort
	}

	if err := checkPortFree(localPort); err != nil {
		return err
	}

	logPath, err := tunnelLogPath(key)
	if err != nil {
		return err
	}

	pid, err := launchTunnel(resolved.Bastion, resolved.Host, resolved.RemotePort, localPort, logPath, tunnelProbe(resolved.DBType))
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
	log.Infof("  Log     : %s", logPath)

	// ElastiCache Serverless requires TLS; --insecure because the certificate
	// is issued for *.serverless.<region>.cache.amazonaws.com, not localhost.
	if resolved.Serverless {
		log.Infof("  Connect : redis-cli -p %d --tls --insecure", localPort)
	}

	log.Infof("  Stop    : draft db:connect stop %s %s", input.DBType, input.Name)

	if !withToken {
		return nil
	}

	result, err := issueIAMToken(resolved, localPort, iamUser)
	if err != nil {
		return fmt.Errorf("tunnel is up but the IAM token could not be issued: %w\n"+
			"retry with: draft db:connect start %s %s", err, input.DBType, input.Name)
	}

	entry := state.Connections[key]
	entry.IAMUser = iamUser
	state.Connections[key] = entry

	if err := saveRuntimeState(state); err != nil {
		log.Warnf("could not persist runtime state: %s", err)
	}

	printIAMTokenInfo(result, input.DBType, input.Name)

	return nil
}
