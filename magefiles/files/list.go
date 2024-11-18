package files

import (
	"os"
)

// List return a list of files and folders in a given path.
func List(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)

	for _, entry := range entries {
		files = append(files, entry.Name())
	}

	return files, nil
}
