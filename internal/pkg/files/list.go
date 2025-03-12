package files

import (
	"os"
)

func ListDirs(path string) ([]string, error) {
	var dirs []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}
