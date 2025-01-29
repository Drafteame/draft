package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/Drafteame/draft/internal/pkg/crypto"
	"github.com/Drafteame/draft/internal/pkg/files"
	nixmetadata "github.com/Drafteame/draft/internal/pkg/nix-metadata"
)

var (
	ErrNoParameter = errors.New("no parameter found in SSM")
	ErrDecryption  = errors.New("failed to decrypt token")
)

const (
	githubTokenPath = "$HOME/.config/home-manager/dotfiles/draftea/draft/.github-token"
	awsProfile      = "draftea-dev"
)

func GetGithubToken() (string, error) {
	if !files.Exists(githubTokenPath) {
		return fetchAndStoreToken()
	}

	return readStoredToken()
}

func RefreshGithubToken() (string, error) {
	return fetchAndStoreToken()
}

func fetchAndStoreToken() (string, error) {
	token, err := getGithubTokenFromSSM()
	if err != nil {
		return "", fmt.Errorf("failed to get token from SSM: %w", err)
	}

	if err = writeGithubTokenFile(token); err != nil {
		return "", fmt.Errorf("failed to store token: %w", err)
	}

	return token, nil
}

func readStoredToken() (string, error) {
	encryptedToken, err := files.Read(githubTokenPath)
	if err != nil {
		return "", fmt.Errorf("failed to read token file: %w", err)
	}

	return decryptGithubToken(encryptedToken)
}

func getGithubTokenFromSSM() (string, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(awsProfile),
	)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	ssmParamName := "/service/github/dev/TOKEN"

	client := ssm.NewFromConfig(cfg)
	result, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &ssmParamName,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get SSM parameter: %w", err)
	}

	if result == nil || result.Parameter == nil {
		return "", ErrNoParameter
	}

	return *result.Parameter.Value, nil
}

func writeGithubTokenFile(token string) error {
	encryptedToken, err := encryptGithubToken(token)
	if err != nil {
		return fmt.Errorf("failed to encrypt token: %w", err)
	}

	return files.Create(githubTokenPath, []byte(encryptedToken))
}

func encryptGithubToken(token string) (string, error) {
	nixMeta, err := nixmetadata.Get()
	if err != nil {
		return "", fmt.Errorf("failed to get nix metadata: %w", err)
	}

	encryptedToken, err := crypto.Encrypt(token, nixMeta.SysUserName)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt token: %w", err)
	}
	return encryptedToken, nil
}

func decryptGithubToken(token []byte) (string, error) {
	nixMeta, err := nixmetadata.Get()
	if err != nil {
		return "", fmt.Errorf("failed to get nix metadata: %w", err)
	}

	decryptedToken, err := crypto.Decrypt(string(token), nixMeta.SysUserName)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt token: %w", err)
	}
	return decryptedToken, nil
}
