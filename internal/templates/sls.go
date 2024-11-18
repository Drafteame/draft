package templates

import (
	"bytes"
	"text/template"
)

type SLS struct {
	ServerlessYAML []byte
	PackageJSON    []byte
	DepsGo         []byte
	Config         SLSConfig
	FrameV2        SLSFrameV2
}

type SLSConfig struct {
	App SLSConfigApp
	Sls SLSConfigSls
}

type SLSConfigSls struct {
	EnvironmentYAML []byte
	IamYAML         []byte
}

type SLSFrameV2 struct {
	Plain  SLSFrameV2Plain
	Sqs    SSLSFrameV2Sqs
	SnsSqs SLSFrameV2SnsSqs
	HTTP   SLSFrameV2HTTP
	Cron   SLSFrameV2Cron
}

func NewSLS(data any) SLS {
	s := SLS{}

	loaders := []func(*SLS, any) error{
		loadServerlessYAML,
		loadPackageJSON,
		loadDepsGo,
		loadConfig,
		loadFrameV2,
	}

	for _, loader := range loaders {
		if err := loader(&s, data); err != nil {
			panic(err)
		}
	}

	return s
}

func loadServerlessYAML(v *SLS, data any) error {
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

func loadPackageJSON(v *SLS, data any) error {
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

func loadDepsGo(v *SLS, data any) error {
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

func loadConfig(v *SLS, data any) error {
	if err := loadConfigApp(v, data); err != nil {
		return err
	}

	return loadConfigSls(v, data)
}

func loadFrameV2(v *SLS, data any) error {
	loaders := []func(*SLS, any) error{
		loadFrameV2Plain,
		loadFrameV2Sqs,
		loadFrameV2SnsSqs,
		loadFrameV2HTTP,
		loadFrameV2Cron,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}
