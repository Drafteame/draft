package templates

func loadConfigSls(v *SLS, data any) error {
	if err := loadConfigSlsEnvironmentYAML(v, data); err != nil {
		return err
	}

	return loadConfigSlsIamYAML(v, data)
}

func loadConfigSlsEnvironmentYAML(v *SLS, data any) error {
	name := "config/sls/environment.yml"
	path := "tmpl/sls/config/sls/environment.yml.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.Config.Sls.EnvironmentYAML = content

	return nil
}

func loadConfigSlsIamYAML(v *SLS, data any) error {
	name := "config/sls/iam.yml"
	path := "tmpl/sls/config/sls/iam.yml.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.Config.Sls.IamYAML = content

	return nil
}
