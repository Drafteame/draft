package templates

type ConfigSls struct {
	EnvironmentYAML []byte
	IamYAML         []byte
}

func loadConfigSls(v ConfigSlsSetter, data any) error {
	configSls := ConfigSls{}

	loaders := []func(*ConfigSls, any) error{
		loadConfigSlsEnvironmentYAML,
		loadConfigSlsIamYAML,
	}

	for _, loader := range loaders {
		if err := loader(&configSls, data); err != nil {
			return err
		}
	}

	v.SetConfigSls(configSls)

	return nil
}

func loadConfigSlsEnvironmentYAML(v *ConfigSls, data any) error {
	name := "config/sls/environment.yml"
	path := "tmpl/sls/config/sls/environment.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.EnvironmentYAML = content

	return nil
}

func loadConfigSlsIamYAML(v *ConfigSls, data any) error {
	name := "config/sls/iam.yml"
	path := "tmpl/sls/config/sls/iam.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.IamYAML = content

	return nil
}
