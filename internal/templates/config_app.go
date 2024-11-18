package templates

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
	name := "config/app/app.pkl"
	path := "tmpl/sls/config/app/app.pkl.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.Config.App.AppPkl = content

	return nil
}

func loadConfigAppModulesPkl(v *SLS, data any) error {
	name := "config/app/modules.pkl"
	path := "tmpl/sls/config/app/modules.pkl.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.Config.App.ModulesPkl = content

	return nil
}
