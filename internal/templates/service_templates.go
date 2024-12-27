package templates

import (
	"bytes"
	"text/template"
)

type ServiceTemplates struct {
	ServerlessYAML []byte
	PackageJSON    []byte
	DepsGo         []byte
	Config         ServiceTemplatesConfig
	Lambda         ServiceTemplatesLambda
}

func (s *ServiceTemplates) SetPlain(plain LambdaPlain) {
	s.Lambda.Plain = plain
}

func (s *ServiceTemplates) SetConfigSls(v ConfigSls) {
	s.Config.Sls = v
}

func (s *ServiceTemplates) SetConfigApp(v ConfigApp) {
	s.Config.App = v
}

type ServiceTemplatesConfig struct {
	App ConfigApp
	Sls ConfigSls
}

type ServiceTemplatesConfigSls struct {
	EnvironmentYAML []byte
	IamYAML         []byte
}

type ServiceTemplatesLambda struct {
	Plain LambdaPlain
}

func NewServiceTemplates(data any) (*ServiceTemplates, error) {
	s := new(ServiceTemplates)

	loaders := []func(*ServiceTemplates, any) error{
		loadServerlessYAML,
		loadPackageJSON,
		loadDepsGo,
		loadConfig,
		loadLambda,
	}

	for _, loader := range loaders {
		if err := loader(s, data); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func loadServerlessYAML(v *ServiceTemplates, data any) error {
	content, err := sls.ReadFile("tmpl/sls/serverless.yml.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("serverless.yml").Parse(string(content))
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)

	if errExec := tmpl.Execute(buff, data); errExec != nil {
		return errExec
	}

	v.ServerlessYAML = buff.Bytes()

	return nil
}

func loadPackageJSON(v *ServiceTemplates, data any) error {
	content, err := sls.ReadFile("tmpl/sls/package.json.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("package.json").Parse(string(content))
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)

	if errExec := tmpl.Execute(buff, data); errExec != nil {
		return errExec
	}

	v.PackageJSON = buff.Bytes()

	return nil
}

func loadDepsGo(v *ServiceTemplates, data any) error {
	content, err := sls.ReadFile("tmpl/sls/deps.go.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("deps.go").Parse(string(content))
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)

	if errExec := tmpl.Execute(buff, data); errExec != nil {
		return errExec
	}

	v.DepsGo = buff.Bytes()

	return nil
}

func loadConfig(v *ServiceTemplates, data any) error {
	if err := loadConfigApp(v, data); err != nil {
		return err
	}

	return loadConfigSls(v, data)
}

func loadLambda(v *ServiceTemplates, data any) error {
	return loadLambdaPlain(v, data)
}
