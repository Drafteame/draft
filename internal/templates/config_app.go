package templates

type ConfigApp struct {
	AppPkl     []byte
	ModulesPkl []byte
}

func loadConfigApp(v ConfigAppSetter, data any) error {
	configApp := ConfigApp{}

	loaders := []func(*ConfigApp, any) error{
		loadConfigAppAppPkl,
		loadConfigAppModulesPkl,
	}

	for _, loader := range loaders {
		if err := loader(&configApp, data); err != nil {
			return err
		}
	}

	v.SetConfigApp(configApp)

	return nil
}

func loadConfigAppAppPkl(v *ConfigApp, data any) error {
	name := "config/app/app.pkl"
	path := "tmpl/sls/config/app/app.pkl.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.AppPkl = content

	return nil
}

func loadConfigAppModulesPkl(v *ConfigApp, data any) error {
	name := "config/app/modules.pkl"
	path := "tmpl/sls/config/app/modules.pkl.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.ModulesPkl = content

	return nil
}
