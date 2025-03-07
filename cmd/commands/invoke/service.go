package invoke

import (
	"errors"
	"strings"

	"github.com/Drafteame/draft/internal/pkg/files"
)

func getService(path string) (string, string, error) {
	parts := strings.Split(path, "/")
	// traverse the path backwards to find the service name
	for i := len(parts) - 1; i >= 0; i-- {
		parent := strings.Join(parts[:i+1], "/")

		if isServerlessService(parent) {
			serviceName := getServiceNameFormServerlessFile(parent)
			return serviceName, parent, nil
		}
	}

	return "", "", errors.New("service is not serverless or cdk")
}

func isServerlessService(path string) bool {
	if !files.Exists(path) {
		return false
	}

	return files.Exists(path + "/serverless.yml")
}

func getServiceNameFormServerlessFile(path string) string {
	file := path + "/serverless.yml"

	type serverless struct {
		Service string `yaml:"service"`
	}

	data := serverless{}

	if err := files.LoadYAML(file, &data); err != nil {
		panic("Failed to load serverless.yml: " + err.Error())
	}

	return data.Service
}
