package project

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func GetPackage(modPath string) (string, error) {
	file, err := os.Open(modPath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("can't read package module name")
}

func NormalizeServiceName(serviceName string) string {
	replacers := []string{" ", "-", "."}

	for _, replacer := range replacers {
		serviceName = strings.ReplaceAll(serviceName, replacer, "_")
	}

	return strings.ToLower(serviceName)
}
