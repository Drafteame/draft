package utils

import "strings"

var partitions = []string{"-", "_", " ", ".", "/", "\\"}

// ToPackageName converts a string to a valid package name
func ToPackageName(name string) string {
	for _, p := range partitions {
		name = strings.ReplaceAll(name, p, "")
	}

	return strings.ToLower(name)
}

// ToCammelCase converts a string to a valid camel case string
func ToCammelCase(name string) string {
	for _, p := range partitions {
		parts := strings.Split(name, p)

		for i, part := range parts {
			parts[i] = strings.Title(part)
		}

		name = strings.Join(parts, "")
	}

	return name
}

// ToSnakeCase converts a string to a valid snake case string
func ToSnakeCase(name string) string {
	for _, p := range partitions {
		parts := strings.Split(name, p)

		for i, part := range parts {
			parts[i] = strings.ToLower(part)
		}

		name = strings.Join(parts, "_")
	}

	return name
}
