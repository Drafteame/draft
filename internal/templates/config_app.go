package templates

import (
	"bytes"
	"text/template"
)

type SLSConfigApp struct {
	AppPkl     []byte
	ModulesPkl []byte
}

func loadConfigApp(v *SLS, data any) error {
	if err := loadConfigAppAppPkl(v, data); err != nil {
		return err
	}

	return loadConfigAppModulesPkl(v, data)
}

func loadConfigAppAppPkl(v *SLS, data any) error {
	content, err := sls.ReadFile("tmpl/sls/config/app/app.pkl.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("config/app/app.pkl").Parse(string(content))
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)

	if errExec := tmpl.Execute(buff, data); errExec != nil {
		return errExec
	}

	v.Config.App.AppPkl = buff.Bytes()

	return nil
}

func loadConfigAppModulesPkl(v *SLS, data any) error {
	content, err := sls.ReadFile("tmpl/sls/config/app/modules.pkl.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("config/app/modules.pkl").Parse(string(content))
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)

	if errExec := tmpl.Execute(buff, data); errExec != nil {
		return errExec
	}

	v.Config.App.ModulesPkl = buff.Bytes()

	return nil
}
