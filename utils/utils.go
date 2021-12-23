package utils

import (
	"fmt"
	"os"

	"github.com/hoisie/mustache"
)

type StringValue string

// RenderTemplate fills the template with the given data and returns the result
func RenderTemplate(template string, data interface{}) (string, error) {
	tmpl, err := mustache.ParseString(template)
	if err != nil {
		return "", err
	}

	rendered := tmpl.Render(data)

	return rendered, err
}

// CreateFile creates a file with the given name and content
func CreateFile(path string, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// CreateFolder creates a folder with the given name
func CreateFolder(path string) error {
	return os.MkdirAll(path, 0777)
}

// ReadFile reads the file and returns its content
func ReadFile(path string) string {
	file, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("engine: error opening file %s: %s\n", path, err)
		os.Exit(1)
	}

	return string(file)
}

// ReplaceFileContent replaces the content of a file with the given content
func ReplaceFileContent(path, content string) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("engine: error opening file %s: %s\n", path, err)
		os.Exit(1)
	}

	if err = f.Truncate(0); err != nil {
		fmt.Printf("engine: error truncating file %s: %s\n", path, err)
		f.Close()
		os.Exit(1)
	}

	if _, err = f.Seek(0, 0); err != nil {
		fmt.Printf("engine: error seeking file %s: %s\n", path, err)
		f.Close()
		os.Exit(1)
	}

	if _, err = fmt.Fprint(f, content); err != nil {
		fmt.Printf("engine: error writing to file %s: %s\n", path, err)
		f.Close()
		os.Exit(1)
	}

	f.Close()
}

// GetCurrentPath obtain and return the path where the command is executed
func GetCurrentPath() string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("engine: error getting current working directory")
		os.Exit(1)
	}

	return wd
}

// PathExists checks if the given path exists
func PathExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("engine: path %s not found\n", path)
		os.Exit(1)
	}
}

// PathNotExists checks if the given path not exists
func PathNotExists(path string) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		fmt.Printf("engine: path %s already exists\n", path)
		os.Exit(1)
	}
}

func NewStringValue(p *string, val string) *StringValue {
	*p = val
	return (*StringValue)(p)
}
