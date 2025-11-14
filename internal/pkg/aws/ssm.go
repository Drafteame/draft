package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"github.com/Drafteame/draft/internal/pkg/log"
)

// GetParameter retrieves a parameter value from AWS SSM Parameter Store
func GetParameter(parameterName string) (string, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile("draftea-dev"),
		config.WithRegion("us-east-2"),
	)
	if err != nil {
		return "", err
	}

	client := ssm.NewFromConfig(cfg)

	input := &ssm.GetParameterInput{
		Name:           &parameterName,
		WithDecryption: boolPtr(true),
	}

	result, err := client.GetParameter(ctx, input)
	if err != nil {
		return "", err
	}

	if result.Parameter == nil || result.Parameter.Value == nil {
		log.Warnf("parameter %s not found or has no value", parameterName)
		return "", nil
	}

	return *result.Parameter.Value, nil
}

func boolPtr(b bool) *bool {
	return &b
}
