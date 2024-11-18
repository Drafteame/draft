package files

import (
	"encoding/json"
	"os"
)

func BindJSON(filePath string, target any) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, target)
}
