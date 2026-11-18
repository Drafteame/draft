package connect

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubAuthToken replaces the aws CLI call for the duration of the test and
// records the arguments it was invoked with.
func stubAuthToken(t *testing.T, out string, err error) *[]string {
	t.Helper()

	var got []string

	orig := runAuthTokenCmd
	runAuthTokenCmd = func(args ...string) (string, error) {
		got = append([]string{}, args...)
		return out, err
	}

	t.Cleanup(func() { runAuthTokenCmd = orig })

	return &got
}

func TestGenerateAuthToken_PassesResolvedValuesToAWSCLI(t *testing.T) {
	got := stubAuthToken(t, "host:5432/?Action=connect&X-Amz-Signature=abc\n", nil)

	bastion := BastionConfig{Profile: "dev-profile", Region: "us-east-2"}

	token, err := generateAuthToken(bastion, "turbo-dev.cluster-aaa.us-east-2.rds.amazonaws.com", 5432, "dev@draftea.com")
	require.NoError(t, err)

	assert.Equal(t, "host:5432/?Action=connect&X-Amz-Signature=abc", token, "token must be trimmed")
	assert.Equal(t, []string{
		"--hostname", "turbo-dev.cluster-aaa.us-east-2.rds.amazonaws.com",
		"--port", "5432",
		"--username", "dev@draftea.com",
		"--region", "us-east-2",
		"--profile", "dev-profile",
	}, *got)
}

func TestGenerateAuthToken_Errors(t *testing.T) {
	t.Run("cli failure", func(t *testing.T) {
		stubAuthToken(t, "", errors.New("boom"))

		_, err := generateAuthToken(BastionConfig{}, "host", 5432, "u")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "boom")
	})

	t.Run("empty token", func(t *testing.T) {
		stubAuthToken(t, "  \n", nil)

		_, err := generateAuthToken(BastionConfig{}, "host", 5432, "u")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "empty token")
	})
}

func TestIssueIAMToken_RejectsNonPostgres(t *testing.T) {
	got := stubAuthToken(t, "token", nil)

	_, err := issueIAMToken(ResolvedConnection{DBType: "redis"}, 56100, "u")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "only supported for postgres")
	assert.Empty(t, *got, "aws must not be called for unsupported engines")
}

func TestIssueIAMToken_WritesPgpassForLocalPort(t *testing.T) {
	stubAuthToken(t, "tok:en\n", nil)

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

	assert.Contains(t, *got, "turbo-prod.cluster-ddd.us-east-2.rds.amazonaws.com")
	assert.Contains(t, *got, "prod-profile")

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
	assert.Empty(t, *got)
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
