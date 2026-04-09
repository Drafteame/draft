package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Drafteame/draft/internal/pkg/exec"
)

var (
	gitRootOnce      sync.Once
	cachedGitRoot    string
	cachedGitRootErr error
)

// resolveService resolves a service name or path to an absolute directory path.
// If arg is an existing directory, it returns the absolute path directly.
// Otherwise it searches the git root for a serverless.yml with service: <arg>.
func resolveService(arg string) (string, error) {
	// Try as path first
	absPath, err := filepath.Abs(arg)
	if err == nil && isDir(absPath) {
		return absPath, nil
	}

	// Treat as service name: find git root and search
	gitRoot, err := getGitRoot()
	if err != nil {
		return "", fmt.Errorf("failed to find git root: %w", err)
	}

	matches, err := findServiceByName(gitRoot, arg)
	if err != nil {
		return "", err
	}

	switch len(matches) {
	case 0:
		return "", fmt.Errorf("service %q not found (no serverless.yml with service: %s)", arg, arg)
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("service %q is ambiguous, found in:\n  %s", arg, strings.Join(matches, "\n  "))
	}
}

// getGitRoot returns the absolute path of the git repository root.
// The result is cached after the first call.
func getGitRoot() (string, error) {
	gitRootOnce.Do(func() {
		out, err := exec.Command("git rev-parse --show-toplevel")
		if err != nil {
			cachedGitRootErr = err
			return
		}
		cachedGitRoot = strings.TrimSpace(out)
	})
	return cachedGitRoot, cachedGitRootErr
}

// findServiceByName walks root looking for serverless.yml files where service: == name.
// Returns a list of absolute directory paths that match.
func findServiceByName(root, name string) ([]string, error) {
	var matches []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}

		// Skip directories that won't contain serverless.yml
		if d.IsDir() {
			base := d.Name()
			if base == ".git" || base == "node_modules" || base == ".serverless" || base == ".bin" {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Name() != "serverless.yml" {
			return nil
		}

		serviceName, err := parseServiceName(path)
		if err != nil {
			return nil // skip unparseable files
		}

		if serviceName == name {
			matches = append(matches, filepath.Dir(path))
		}

		return nil
	})

	return matches, err
}

// isDir returns true if path exists and is a directory.
func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
