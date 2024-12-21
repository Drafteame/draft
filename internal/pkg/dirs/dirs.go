package dirs

import (
	"os"
)

// Create creates directories.
func Create(paths ...string) error {
	for _, path := range paths {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}
