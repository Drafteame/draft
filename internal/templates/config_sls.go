package templates

type ConfigSls struct {
	EnvironmentYAML []byte
	ResourcesYAML   []byte
}

func loadConfigSls(v ConfigSlsSetter, data any) error {
	configSls := ConfigSls{}

	loaders := []func(*ConfigSls, any) error{
		loadConfigSlsEnvironmentYAML,
		loadConfigSlsResourcesYAML,
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

func loadConfigSlsResourcesYAML(v *ConfigSls, data any) error {
	name := "config/sls/resources.yml"
	path := "tmpl/sls/config/sls/resources.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.ResourcesYAML = content

	return nil
}
