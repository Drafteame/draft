package templates

import (
	"bytes"
	"embed"
	"text/template"
)

var (
	//go:embed tmpl/sls
	sls embed.FS
)

func loadTemplate(name, path string, data any) ([]byte, error) {
	content, err := sls.ReadFile(path)
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
