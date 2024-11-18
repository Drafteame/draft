package files

import (
	"encoding/json"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadYAML(path string, v any) error {
	data, err := readFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, v)
}

func LoadJSON(path string, v any) error {
	data, err := readFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func readFile(path string) ([]byte, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := file.Close(); errClose != nil {
			panic(errClose)
		}
	}()

	return io.ReadAll(file)
}
