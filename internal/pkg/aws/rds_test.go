package aws

import (
	"context"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildRDSAuthToken(t *testing.T) {
	creds := credentials.NewStaticCredentialsProvider("AKIDEXAMPLE", "secret", "")

	input := RDSAuthTokenInput{
		Region: "us-east-2",
		Host:   "turbo-dev.cluster-aaa.us-east-2.rds.amazonaws.com",
		Port:   5432,
		User:   "dev@draftea.com",
	}

	token, err := BuildRDSAuthToken(context.Background(), input, creds)
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(token, "turbo-dev.cluster-aaa.us-east-2.rds.amazonaws.com:5432?"),
		"token must be the signed endpoint without scheme, got %q", token)

	parsed, err := url.Parse("https://" + token)
	require.NoError(t, err)

	q := parsed.Query()
	assert.Equal(t, "connect", q.Get("Action"))
	assert.Equal(t, "dev@draftea.com", q.Get("DBUser"))
	assert.Equal(t, "900", q.Get("X-Amz-Expires"), "RDS tokens are valid for 15 minutes")
	assert.Contains(t, q.Get("X-Amz-Credential"), "AKIDEXAMPLE/")
	assert.Contains(t, q.Get("X-Amz-Credential"), "/us-east-2/rds-db/aws4_request")
	assert.NotEmpty(t, q.Get("X-Amz-Signature"))
}

func TestGenerateRDSAuthToken_ResolvesProfileFromSharedConfig(t *testing.T) {
	dir := t.TempDir()

	cfgFile := dir + "/config"
	require.NoError(t, os.WriteFile(cfgFile, []byte("[profile dev-profile]\nregion = us-west-1\n"), 0600))

	credFile := dir + "/credentials"
	require.NoError(t, os.WriteFile(credFile, []byte("[dev-profile]\naws_access_key_id = AKIDPROFILE\naws_secret_access_key = secret\n"), 0600))

	t.Setenv("AWS_CONFIG_FILE", cfgFile)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", credFile)
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	t.Setenv("AWS_PROFILE", "")

	token, err := GenerateRDSAuthToken(context.Background(), RDSAuthTokenInput{
		Profile: "dev-profile",
		Region:  "us-east-2",
		Host:    "host",
		Port:    5432,
		User:    "u",
	})
	require.NoError(t, err)

	parsed, err := url.Parse("https://" + token)
	require.NoError(t, err)

	cred := parsed.Query().Get("X-Amz-Credential")
	assert.Contains(t, cred, "AKIDPROFILE/", "credentials must come from the requested profile")
	assert.Contains(t, cred, "/us-east-2/", "explicit region must win over the profile region")

	_, err = GenerateRDSAuthToken(context.Background(), RDSAuthTokenInput{Profile: "missing", Region: "us-east-2", Host: "host", Port: 5432, User: "u"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), `profile "missing"`)
}
