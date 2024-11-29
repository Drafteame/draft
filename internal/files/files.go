package files

import (
	"io"
	"os"
)

func Read(path string) ([]byte, error) {
	path = os.ExpandEnv(path)
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}

func Create(path string, newContent []byte) error {
	return os.WriteFile(path, newContent, 0755)
}

func Exists(path string) bool {
	path = os.ExpandEnv(path)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
