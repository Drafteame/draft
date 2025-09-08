package templates

func nativeLoadLambdaSnsSqsMainGo(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/main.go"
	path := "tmpl/sls/native/snssqs/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.MainGo = content

	return nil
}

func nativeLoadLambdaSnsSqsLambdaConfigYAML(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/lambda-config.yml"
	path := "tmpl/sls/native/snssqs/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.LambdaConfigYAML = content

	return nil
}

func nativeLoadLambdaSnsSqsHandler(v *LambdaSnsSqs, data any) error {
	loaders := []func(*LambdaSnsSqs, any) error{
		nativeLoadLambdaSnsSqsHandlerBoostrapGo,
		nativeLoadLambdaSnsSqsHandlerWorkerWorkerGo,
		nativeLoadLambdaSnsSqsHandlerWorkerResourcesGo,
		nativeLoadLambdaSnsSqsHandlerEmbedYaml,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func nativeLoadLambdaSnsSqsHandlerBoostrapGo(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/handler/bootstrap.go"
	path := "tmpl/sls/native/snssqs/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.BootstrapGo = content

	return nil
}

func nativeLoadLambdaSnsSqsHandlerWorkerWorkerGo(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/handler/worker/worker.go"
	path := "tmpl/sls/native/snssqs/handler/worker/worker.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.WorkerGo = content

	return nil
}

func nativeLoadLambdaSnsSqsHandlerWorkerResourcesGo(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/handler/worker/resources.go"
	path := "tmpl/sls/native/snssqs/handler/worker/resources.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.ResourcesGo = content

	return nil
}

func nativeLoadLambdaSnsSqsHandlerEmbedYaml(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/handler/embed/embed.yaml"
	path := "tmpl/sls/native/snssqs/handler/embed/embed.yaml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.EmbedYML = content

	return nil
}
