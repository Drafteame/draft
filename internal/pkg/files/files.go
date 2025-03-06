package files

import (
	"bufio"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ScanAndWrite(path string, cb func(line string) (string, error)) error {
	path = os.ExpandEnv(path)

	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	newContent := strings.Builder{}

	for scanner.Scan() {
		line := scanner.Text()
		newLine, err := cb(line)
		if err != nil {
			return err
		}

		if _, errWrite := newContent.WriteString(newLine); errWrite != nil {
			return errWrite
		}
	}

	return Create(path, []byte(newContent.String()))
}

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
