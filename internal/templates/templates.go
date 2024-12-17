package templates

import (
	"bytes"
	"embed"
	"path/filepath"
	"text/template"
)

var (
	//go:embed tmpl/sls
	sls embed.FS
	//go:embed tmpl/domain
	domain embed.FS
)

func PrintSlsFiles() {
	printAllFilePaths(sls, "tmpl/sls")
}

func PrintDomainFiles() {
	printAllFilePaths(domain, "tmpl/domain")
}

func loadTemplate(name, path string, data any, fs embed.FS) ([]byte, error) {
	content, err := fs.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)

	if errExec := tmpl.Execute(buff, data); errExec != nil {
		return nil, errExec
	}

	return buff.Bytes(), nil
}

func printAllFilePaths(fs embed.FS, path string) {
	files, err := fs.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		itemPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			printAllFilePaths(fs, itemPath)
		} else {
			println(itemPath)
		}
	}
}
