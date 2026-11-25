package aws

import (
	"context"
	"fmt"
	"net"
	"strconv"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
)

// RDSAuthTokenInput identifies the database and the credentials used to sign
// an RDS IAM authentication token.
type RDSAuthTokenInput struct {
	Profile string // shared config profile that owns the credentials
	Region  string
	Host    string // RDS endpoint, e.g. turbo-dev.cluster-xxx.us-east-2.rds.amazonaws.com
	Port    int
	User    string // database user mapped to the IAM identity
}

// GenerateRDSAuthToken signs an RDS IAM auth token for the given database.
// Signing happens locally (no API call); the token is accepted for new logins
// during 15 minutes. Credentials are resolved from the shared config profile
// (~/.aws/config and ~/.aws/credentials), the same chain the AWS CLI uses.
func GenerateRDSAuthToken(ctx context.Context, input RDSAuthTokenInput) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(input.Profile),
		config.WithRegion(input.Region),
	)
	if err != nil {
		return "", fmt.Errorf("could not load AWS profile %q: %w", input.Profile, err)
	}

	return BuildRDSAuthToken(ctx, input, cfg.Credentials)
}

// BuildRDSAuthToken signs an RDS IAM auth token with the given credentials.
func BuildRDSAuthToken(ctx context.Context, input RDSAuthTokenInput, creds sdkaws.CredentialsProvider) (string, error) {
	endpoint := net.JoinHostPort(input.Host, strconv.Itoa(input.Port))

	token, err := auth.BuildAuthToken(ctx, endpoint, input.Region, input.User, creds)
	if err != nil {
		return "", fmt.Errorf("could not build RDS auth token for %s: %w", endpoint, err)
	}

	return token, nil
}
