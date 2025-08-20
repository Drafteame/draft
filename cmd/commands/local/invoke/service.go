package invoke

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Drafteame/draft/internal/pkg/files"
)

func getService(path string) (string, string, error) {
	parts := strings.Split(path, "/")
	// traverse the path backwards to find the service name
	for i := len(parts) - 1; i >= 0; i-- {
		parent := strings.Join(parts[:i+1], "/")

		if isServerlessService(parent) {
			serviceName, errSlsFile := getServiceNameFormServerlessFile(parent)
			if errSlsFile != nil {
				return "", "", errSlsFile
			}

			return serviceName, parent, nil
		}
	}

	return "", "", errors.New("service does not have a valid serverless.yaml file")
}

func isServerlessService(path string) bool {
	if !files.Exists(path) {
		return false
	}

	return files.Exists(path + "/serverless.yml")
}

func getServiceNameFormServerlessFile(path string) (string, error) {
	file := path + "/serverless.yml"

	type serverless struct {
		Service string `yaml:"service"`
	}

	data := serverless{}

	if err := files.LoadYAML(file, &data); err != nil {
		return "", fmt.Errorf("failed to load serverless.yml: %w", err)
	}

	return data.Service, nil
}
