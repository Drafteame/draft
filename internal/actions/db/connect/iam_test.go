package connect

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Drafteame/draft/internal/pkg/aws"
)

// stubAuthToken replaces the SDK token signing for the duration of the test
// and records the input it was invoked with.
func stubAuthToken(t *testing.T, token string, err error) *aws.RDSAuthTokenInput {
	t.Helper()

	var got aws.RDSAuthTokenInput

	orig := generateRDSAuthToken
	generateRDSAuthToken = func(_ context.Context, input aws.RDSAuthTokenInput) (string, error) {
		got = input
		return token, err
	}

	t.Cleanup(func() { generateRDSAuthToken = orig })

	return &got
}

func TestGenerateAuthToken_PassesResolvedValuesToSigner(t *testing.T) {
	got := stubAuthToken(t, "host:5432/?Action=connect&X-Amz-Signature=abc", nil)

	bastion := BastionConfig{Profile: "dev-profile", Region: "us-east-2"}

	token, err := generateAuthToken(bastion, "turbo-dev.cluster-aaa.us-east-2.rds.amazonaws.com", 5432, "dev@draftea.com")
	require.NoError(t, err)

	assert.Equal(t, "host:5432/?Action=connect&X-Amz-Signature=abc", token)
	assert.Equal(t, aws.RDSAuthTokenInput{
		Profile: "dev-profile",
		Region:  "us-east-2",
		Host:    "turbo-dev.cluster-aaa.us-east-2.rds.amazonaws.com",
		Port:    5432,
		User:    "dev@draftea.com",
	}, *got)
}

func TestGenerateAuthToken_PropagatesSignerError(t *testing.T) {
	stubAuthToken(t, "", errors.New("boom"))

	_, err := generateAuthToken(BastionConfig{}, "host", 5432, "u")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}

func TestIssueIAMToken_RejectsNonPostgres(t *testing.T) {
	got := stubAuthToken(t, "token", nil)

	_, err := issueIAMToken(ResolvedConnection{DBType: "redis"}, 56100, "u")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "only supported for postgres")
	assert.Empty(t, got.Host, "the signer must not be called for unsupported engines")
}

func TestIssueIAMToken_WritesPgpassForLocalPort(t *testing.T) {
	stubAuthToken(t, "tok:en", nil)

	pgpass := filepath.Join(t.TempDir(), ".pgpass")
	t.Setenv("PGPASSFILE", pgpass)

	cfg := loadTestConfig(t)
	resolved, err := cfg.ResolveConnection("postgres", "turbo-dev")
	require.NoError(t, err)

	result, err := issueIAMToken(resolved, 15432, "dev@draftea.com")
	require.NoError(t, err)

	assert.Equal(t, "dev@draftea.com", result.User)
	assert.Equal(t, pgpass, result.PgpassPath)
	assert.WithinDuration(t, result.ExpiresAt, result.ExpiresAt, iamTokenTTL)

	data, err := os.ReadFile(pgpass)
	require.NoError(t, err)
	assert.Equal(t, `localhost:15432:*:dev@draftea.com:tok\:en`+"\n", string(data),
		"entry must point at the local tunnel port, not the remote RDS port")
}

func TestResolveIAMUser(t *testing.T) {
	t.Setenv("DRAFTEA_EMAIL", " env@draftea.com ")

	user, err := resolveIAMUser()
	require.NoError(t, err)
	assert.Equal(t, "env@draftea.com", user)

	t.Setenv("DRAFTEA_EMAIL", "  ")

	_, err = resolveIAMUser()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DRAFTEA_EMAIL is not set")
	assert.Contains(t, err.Error(), "hms-set-email")
}

func TestRefreshToken_UsesTunnelPortAndRecordsUser(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	pgpass := filepath.Join(home, "pgpass")
	t.Setenv("PGPASSFILE", pgpass)

	// Previous token for another user must be dropped.
	require.NoError(t, os.WriteFile(pgpass, []byte("localhost:15432:*:old@draftea.com:old\n"), 0600))

	got := stubAuthToken(t, "token-value", nil)

	cfg := loadTestConfig(t)
	resolved, err := cfg.ResolveConnection("postgres", "turbo-prod")
	require.NoError(t, err)

	key := runtimeKey("postgres", "turbo-prod")
	state := RuntimeState{Connections: map[string]RuntimeEntry{
		key: {DBType: "postgres", Name: "turbo-prod", PID: os.Getpid(), LocalPort: 15432, IAMUser: "old@draftea.com"},
	}}

	require.NoError(t, refreshToken(state, key, state.Connections[key], resolved, "dev@draftea.com"))

	assert.Equal(t, "turbo-prod.cluster-ddd.us-east-2.rds.amazonaws.com", got.Host)
	assert.Equal(t, "prod-profile", got.Profile)

	data, err := os.ReadFile(pgpass)
	require.NoError(t, err)
	assert.Equal(t, "localhost:15432:*:dev@draftea.com:token-value\n", string(data),
		"entry must use the tunnel's local port and the old user's entry must be gone")

	assert.Equal(t, "dev@draftea.com", state.Connections[key].IAMUser)

	saved, err := loadRuntimeState()
	require.NoError(t, err)
	assert.Equal(t, "dev@draftea.com", saved.Connections[key].IAMUser, "state must be persisted")
}

func TestStart_FailsBeforeCallingAWSWhenUserIsUnknown(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("DRAFTEA_EMAIL", "")

	require.NoError(t, os.MkdirAll(filepath.Join(home, ".draft"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(home, ".draft", configFileName), []byte(testConfigYAML), 0644))

	got := stubAuthToken(t, "token", nil)

	err := Start(StartInput{DBType: "postgres", Name: "turbo-dev", Token: true})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DRAFTEA_EMAIL is not set")
	assert.Empty(t, got.Host)
}

func TestStart_RejectsTokenForNonPostgres(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("DRAFTEA_EMAIL", "dev@draftea.com")

	require.NoError(t, os.MkdirAll(filepath.Join(home, ".draft"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(home, ".draft", configFileName), []byte(testConfigYAML), 0644))

	err := Start(StartInput{DBType: "mongo", Name: "main-dev", Token: true})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--token is only supported for postgres")
}
