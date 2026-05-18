package project

import (
	"strings"

	"github.com/Drafteame/draft/internal/pkg/files"
)

// IsOtelService reads the service's serverless.yml and returns true if the
// build command contains the otel build tag, indicating the service uses OpenTelemetry.
func IsOtelService(servicePath string) (bool, error) {
	content, err := files.Read(servicePath + "/serverless.yml")
	if err != nil {
		return false, err
	}

	return strings.Contains(string(content), `"lambda.norpc,otel"`), nil
}
