package dirs

import (
	"os"

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
