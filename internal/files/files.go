package files

import (
	"io"
	"os"
)

func Read(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}

func Create(path string, newContent []byte) error {
	return os.WriteFile(path, newContent, 0755)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
