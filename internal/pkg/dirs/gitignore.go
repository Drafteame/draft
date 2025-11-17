package dirs

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// gitignoreMatcher holds gitignore patterns and provides matching functionality
type gitignoreMatcher struct {
	patterns []gitignorePattern
}

type gitignorePattern struct {
	pattern   string
	isDir     bool
	isNegated bool
	isRooted  bool
}

// loadGitignore loads .gitignore patterns from the root directory
func loadGitignore(root string) *gitignoreMatcher {
	matcher := &gitignoreMatcher{
		patterns: make([]gitignorePattern, 0),
	}

	gitignorePath := filepath.Join(root, ".gitignore")
	file, err := os.Open(gitignorePath)
	if err != nil {
		// .gitignore doesn't exist, return empty matcher
		return matcher
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		pattern := gitignorePattern{
			pattern:   line,
			isNegated: strings.HasPrefix(line, "!"),
			isRooted:  strings.HasPrefix(line, "/"),
		}

		// Remove negation prefix
		if pattern.isNegated {
			pattern.pattern = strings.TrimPrefix(pattern.pattern, "!")
		}

		// Remove root prefix
		if pattern.isRooted {
			pattern.pattern = strings.TrimPrefix(pattern.pattern, "/")
		}

		// Check if pattern is for directories only
		if strings.HasSuffix(pattern.pattern, "/") {
			pattern.isDir = true
			pattern.pattern = strings.TrimSuffix(pattern.pattern, "/")
		}

		matcher.patterns = append(matcher.patterns, pattern)
	}

	return matcher
}

// shouldIgnore checks if a path should be ignored based on gitignore patterns
func (m *gitignoreMatcher) shouldIgnore(path, root string, isDir bool) bool {
	// Get relative path from root
	relPath, err := filepath.Rel(root, path)
	if err != nil || relPath == "." {
		return false
	}

	// Normalize path separators
	relPath = filepath.ToSlash(relPath)

	var ignored bool

	for _, pattern := range m.patterns {
		// Skip directory-only patterns for files
		if pattern.isDir && !isDir {
			continue
		}

		matched := m.matchPattern(pattern.pattern, relPath, pattern.isRooted)

		if matched {
			if pattern.isNegated {
				ignored = false
			} else {
				ignored = true
			}
		}
	}

	return ignored
}

// matchPattern matches a gitignore pattern against a path
func (m *gitignoreMatcher) matchPattern(pattern, path string, isRooted bool) bool {
	// Handle special patterns
	if pattern == "*" {
		return true
	}

	// If rooted, match from the beginning
	if isRooted {
		matched, _ := filepath.Match(pattern, path)
		if matched {
			return true
		}

		// Also check if any parent directory matches
		parts := strings.Split(path, "/")
		for i := 1; i < len(parts); i++ {
			subPath := strings.Join(parts[:i], "/")
			matched, _ := filepath.Match(pattern, subPath)
			if matched {
				return true
			}
		}
		return false
	}

	// Not rooted, match against any part of the path
	matched, _ := filepath.Match(pattern, path)
	if matched {
		return true
	}

	// Check if pattern matches any component or suffix
	parts := strings.Split(path, "/")
	for i := 0; i < len(parts); i++ {
		subPath := strings.Join(parts[i:], "/")
		matched, _ := filepath.Match(pattern, subPath)
		if matched {
			return true
		}

		// Also check individual components
		matched, _ = filepath.Match(pattern, parts[i])
		if matched {
			return true
		}
	}

	// Handle directory glob patterns like "node_modules/**"
	if strings.Contains(pattern, "**") {
		simplePattern := strings.ReplaceAll(pattern, "**", "*")
		matched, _ := filepath.Match(simplePattern, path)
		if matched {
			return true
		}
	}

	return false
}
