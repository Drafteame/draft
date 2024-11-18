package templates

import (
	"bytes"
	"text/template"
)

func loadConfigSls(v *SLS, data any) error {
	if err := loadConfigSlsEnvironmentYAML(v, data); err != nil {
		return err
	}

	return loadConfigSlsIamYAML(v, data)
}

func loadConfigSlsEnvironmentYAML(v *SLS, data any) error {
	content, err := sls.ReadFile("tmpl/sls/config/sls/environment.yml.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("config/sls/environment.yml").Parse(string(content))
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)

	if errExec := tmpl.Execute(buff, data); errExec != nil {
		return errExec
	}

	v.Config.Sls.EnvironmentYAML = buff.Bytes()

	return nil
}

func loadConfigSlsIamYAML(v *SLS, data any) error {
	content, err := sls.ReadFile("tmpl/sls/config/sls/iam.yml.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("config/sls/iam.yml").Parse(string(content))
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)

	if errExec := tmpl.Execute(buff, data); errExec != nil {
		return errExec
	}

	v.Config.Sls.IamYAML = buff.Bytes()

	return nil
}
