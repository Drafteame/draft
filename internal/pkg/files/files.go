package files

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func Read(path string) ([]byte, error) {
	path = os.ExpandEnv(path)
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}

func ReadString(path string) (string, error) {
	content, err := Read(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func Create(path string, newContent []byte) error {
	path = os.ExpandEnv(path)
	return os.WriteFile(path, newContent, 0755)
}

func Exists(path string) bool {
	path = os.ExpandEnv(path)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func LoadYAML(path string, v any) error {
	data, err := Read(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, v)
}
