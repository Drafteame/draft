package dirs

import (
	"os"
	"path/filepath"

	"github.com/Drafteame/draft/internal/pkg/constants"
)

// Create creates directories with default permissions.
func Create(paths ...string) error {
	for _, path := range paths {
		if err := os.MkdirAll(path, constants.DefaultDirMode); err != nil {
			return err
		}
	}

	return nil
}

// WalkFunc is the type of the function called by Walk for each file or directory.
type WalkFunc func(path string, info os.FileInfo, err error) error

// Walk walks the file tree rooted at root, calling fn for each file or directory
// in the tree, including root. By default, it respects .gitignore patterns.
// If you want to ignore .gitignore, pass the skipGitignore option.
func Walk(root string, fn WalkFunc, skipGitignore ...bool) error {
	var matcher *gitignoreMatcher

	// Load .gitignore patterns unless explicitly skipped
	if len(skipGitignore) == 0 || !skipGitignore[0] {
		matcher = loadGitignore(root)
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fn(path, info, err)
		}

		// Skip if matched by gitignore
		if matcher != nil && matcher.shouldIgnore(path, root, info.IsDir()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		return fn(path, info, err)
	})
}
