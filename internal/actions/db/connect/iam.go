package connect

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Drafteame/draft/internal/pkg/aws"
	"github.com/Drafteame/draft/internal/pkg/log"
)

const iamUserEnv = "DRAFTEA_EMAIL"

// iamTokenTTL is how long an RDS IAM auth token is accepted by the database.
const iamTokenTTL = 15 * time.Minute

// iamDBType is the only engine that supports IAM database authentication here.
const iamDBType = "postgres"

// generateRDSAuthToken signs the token; a variable so tests can stub it.
var generateRDSAuthToken = aws.GenerateRDSAuthToken

// generateAuthToken builds an RDS IAM auth token for the real RDS host and
// remote port, using the bastion profile/region of the environment.
func generateAuthToken(bastion BastionConfig, host string, remotePort int, user string) (string, error) {
	return generateRDSAuthToken(context.Background(), aws.RDSAuthTokenInput{
		Profile: bastion.Profile,
		Region:  bastion.Region,
		Host:    host,
		Port:    remotePort,
		User:    user,
	})
}

// iamTokenResult describes a token that was written to the password file.
type iamTokenResult struct {
	User       string
	PgpassPath string
	ExpiresAt  time.Time
}

// issueIAMToken generates a token for the resolved postgres connection and
// stores it in the password file as localhost:<localPort>:*:<user>:<token>.
func issueIAMToken(resolved ResolvedConnection, localPort int, user string) (iamTokenResult, error) {
	if resolved.DBType != iamDBType {
		return iamTokenResult{}, fmt.Errorf("IAM authentication is only supported for %s connections", iamDBType)
	}

	token, err := generateAuthToken(resolved.Bastion, resolved.Host, resolved.RemotePort, user)
	if err != nil {
		return iamTokenResult{}, err
	}

	path, err := pgpassPath()
	if err != nil {
		return iamTokenResult{}, err
	}

	if err := upsertPgpassEntry(path, localPort, user, token); err != nil {
		return iamTokenResult{}, err
	}

	return iamTokenResult{
		User:       user,
		PgpassPath: path,
		ExpiresAt:  time.Now().Add(iamTokenTTL),
	}, nil
}

// cleanupIAMToken removes the password-file entry written for a tunnel.
func cleanupIAMToken(localPort int, user string) {
	if user == "" {
		return
	}

	path, err := pgpassPath()
	if err != nil {
		log.Warnf("could not remove IAM token from password file: %s", err)
		return
	}

	if err := removePgpassEntry(path, localPort, user); err != nil {
		log.Warnf("could not remove IAM token from password file: %s", err)
	}
}

// resolveIAMUser returns the postgres user for IAM auth: the Draftea email
// exported by home-manager in DRAFTEA_EMAIL. There is no override on purpose:
// the token must be for the identity that owns the AWS credentials.
func resolveIAMUser() (string, error) {
	if user := strings.TrimSpace(os.Getenv(iamUserEnv)); user != "" {
		return user, nil
	}

	return "", fmt.Errorf("%s is not set: run 'hms-set-email name.lastname@draftea.com'", iamUserEnv)
}

// printIAMTokenInfo prints the user the token was issued for and how to refresh it.
func printIAMTokenInfo(result iamTokenResult, dbType, name string) {
	log.Infof("  IAM user: %s", result.User)
	log.Infof("  Refresh : draft db:connect start %s %s", dbType, name)
}

// refreshToken issues a new token for a tunnel that is already running and
// records the user in the runtime state, without touching the tunnel.
func refreshToken(state RuntimeState, key string, entry RuntimeEntry, resolved ResolvedConnection, user string) error {
	result, err := issueIAMToken(resolved, entry.LocalPort, user)
	if err != nil {
		return err
	}

	// The email may have changed since the tunnel was started (hms-set-email);
	// a different user than last time leaves a stale entry behind, drop it.
	if entry.IAMUser != "" && entry.IAMUser != user {
		cleanupIAMToken(entry.LocalPort, entry.IAMUser)
	}

	entry.IAMUser = user
	state.Connections[key] = entry

	if err := saveRuntimeState(state); err != nil {
		log.Warnf("could not persist runtime state: %s", err)
	}

	log.Successf("Tunnel already running, token refreshed: %s/%s (PID: %d)", entry.DBType, entry.Name, entry.PID)
	printIAMTokenInfo(result, entry.DBType, entry.Name)

	return nil
}
