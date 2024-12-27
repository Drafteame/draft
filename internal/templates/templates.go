package templates

import (
	"bytes"
	"embed"
	"text/template"
)

var (
	//go:embed tmpl/sls
	sls embed.FS

	//go:embed tmpl/domain
	domain embed.FS
)

type CronSetter interface {
	SetCron(LambdaCron)
}

type HTTPSetter interface {
	SetHTTP(LambdaHTTP)
}

type PlainSetter interface {
	SetPlain(LambdaPlain)
}

type SnsSqsSetter interface {
	SetSnsSqs(LambdaSnsSqs)
}

type SqsSetter interface {
	SetSqs(LambdaSqs)
}

type ConfigAppSetter interface {
	SetConfigApp(ConfigApp)
}

type ConfigSlsSetter interface {
	SetConfigSls(ConfigSls)
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
