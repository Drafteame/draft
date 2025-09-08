package templates

func nativeLoadLambdaSqsMainGo(v *LambdaSqs, data any) error {
	name := "native/sqs/main.go"
	path := "tmpl/sls/native/sqs/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.MainGo = content

	return nil
}

func nativeLoadLambdaSqsLambdaConfigYAML(v *LambdaSqs, data any) error {
	name := "native/sqs/lambda-config.yml"
	path := "tmpl/sls/native/sqs/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.LambdaConfigYAML = content

	return nil
}

func nativeLoadLambdaSqsHandler(v *LambdaSqs, data any) error {
	loaders := []func(*LambdaSqs, any) error{
		nativeLoadLambdaSqsHandlerBoostrapGo,
		nativeLoadLambdaSqsHandlerWorkerWorkerGo,
		nativeLoadLambdaSqsHandlerWorkerResourcesGo,
		nativeLoadLambdaSqsHandlerEmbedYaml,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func nativeLoadLambdaSqsHandlerBoostrapGo(v *LambdaSqs, data any) error {
	name := "native/sqs/handler/bootstrap.go"
	path := "tmpl/sls/native/sqs/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.BootstrapGo = content

	return nil
}

func nativeLoadLambdaSqsHandlerWorkerWorkerGo(v *LambdaSqs, data any) error {
	name := "native/sqs/handler/worker/worker.go"
	path := "tmpl/sls/native/sqs/handler/worker/worker.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.WorkerGo = content

	return nil
}

func nativeLoadLambdaSqsHandlerWorkerResourcesGo(v *LambdaSqs, data any) error {
	name := "native/sqs/handler/worker/resources.go"
	path := "tmpl/sls/native/sqs/handler/worker/resources.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.ResourcesGo = content

	return nil
}

func nativeLoadLambdaSqsHandlerEmbedYaml(v *LambdaSqs, data any) error {
	name := "native/sqs/handler/embed/embed.yaml"
	path := "tmpl/sls/native/sqs/handler/embed/embed.yaml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.EmbedYML = content

	return nil
}
