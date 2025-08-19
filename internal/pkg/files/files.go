package files

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/Drafteame/draft/internal/pkg/dirs"
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

// LoadYAML reads a YAML file from the given path and unmarshals its content into the provided variable `v`.
func LoadYAML(path string, v any) error {
	data, err := Read(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, v)
}

// Copy copies a file from the source path to the destination path, preserving the file mode of the source file.
// It expands environment variables in the provided source and destination paths before processing.
// If the destination directory does not exist, it creates the directory structure.
func Copy(src, dst string) error {
	src = os.ExpandEnv(src)
	dst = os.ExpandEnv(dst)

	// Open source file
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	// Ensure destination directory exists
	if errDirCreate := dirs.Create(filepath.Dir(dst)); errDirCreate != nil {
		return errDirCreate
	}

	// Preserve file mode from source
	info, err := in.Stat()
	if err != nil {
		return err
	}
	mode := info.Mode()

	// Create a destination file
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	// Copy contents
	if _, errCopy := io.Copy(out, in); errCopy != nil {
		return errCopy
	}

	// Flush to disk
	return out.Sync()
}
