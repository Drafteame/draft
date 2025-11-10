package templates

type LambdaCustom struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          LambdaCustomHandler
}

type LambdaCustomHandler struct {
	BootstrapGo       []byte
	HandlerGo         []byte
	WorkerGo          []byte
	ProviderGo        []byte
	EmbedYML          []byte
	ResourcesGo       []byte
	IdempotencyGo     []byte
	InterfacesGo      []byte
	WorkerSetupTestGo []byte
	WorkerTestGo      []byte
}

func loadLambdaCustom(v CustomSetter, data any) error {
	custom := LambdaCustom{}

	loaders := []func(*LambdaCustom, any) error{
		loadLambdaCustomMainGo,
		loadLambdaCustomLambdaConfigYAML,
		loadLambdaCustomHandler,
	}

	for _, loader := range loaders {
		if err := loader(&custom, data); err != nil {
			return err
		}
	}

	v.SetCustom(custom)

	return nil
}

func loadLambdaCustomMainGo(v *LambdaCustom, data any) error {
	name := "native/custom/main.go"
	path := "tmpl/sls/native/custom/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.MainGo = content

	return nil
}

func loadLambdaCustomLambdaConfigYAML(v *LambdaCustom, data any) error {
	name := "native/custom/lambda-config.yml"
	path := "tmpl/sls/native/custom/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.LambdaConfigYAML = content

	return nil
}

func loadLambdaCustomHandler(v *LambdaCustom, data any) error {
	loaders := []func(*LambdaCustom, any) error{
		loadLambdaCustomHandlerBoostrapGo,
		loadLambdaCustomHandlerWorkerWorkerGo,
		loadLambdaCustomHandlerWorkerResourcesGo,
		loadLambdaCustomHandlerWorkerIdempotencyGo,
		loadLambdaCustomHandlerWorkerInterfacesGo,
		loadLambdaCustomHandlerWorkerSetupTestGo,
		loadLambdaCustomHandlerWorkerTestGo,
		loadLambdaCustomHandlerEmbedYaml,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadLambdaCustomHandlerBoostrapGo(v *LambdaCustom, data any) error {
	name := "native/custom/handler/bootstrap.go"
	path := "tmpl/sls/native/custom/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.BootstrapGo = content

	return nil
}

func loadLambdaCustomHandlerWorkerWorkerGo(v *LambdaCustom, data any) error {
	name := "native/custom/handler/worker/worker.go"
	path := "tmpl/sls/native/custom/handler/worker/worker.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.WorkerGo = content

	return nil
}

func loadLambdaCustomHandlerWorkerResourcesGo(v *LambdaCustom, data any) error {
	name := "native/custom/handler/worker/resources.go"
	path := "tmpl/sls/native/custom/handler/worker/resources.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.ResourcesGo = content

	return nil
}

func loadLambdaCustomHandlerEmbedYaml(v *LambdaCustom, data any) error {
	name := "native/custom/handler/embed/embed.yaml"
	path := "tmpl/sls/native/custom/handler/embed/embed.yaml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.EmbedYML = content

	return nil
}

func loadLambdaCustomHandlerWorkerIdempotencyGo(v *LambdaCustom, data any) error {
	name := "native/custom/handler/worker/idempotency.go"
	path := "tmpl/sls/native/custom/handler/worker/idempotency.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.IdempotencyGo = content

	return nil
}

func loadLambdaCustomHandlerWorkerInterfacesGo(v *LambdaCustom, data any) error {
	name := "native/custom/handler/worker/interfaces.go"
	path := "tmpl/sls/native/custom/handler/worker/interfaces.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.InterfacesGo = content

	return nil
}

func loadLambdaCustomHandlerWorkerSetupTestGo(v *LambdaCustom, data any) error {
	name := "native/custom/handler/worker/worker_setup_test.go"
	path := "tmpl/sls/native/custom/handler/worker/worker_setup_test.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.WorkerSetupTestGo = content

	return nil
}

func loadLambdaCustomHandlerWorkerTestGo(v *LambdaCustom, data any) error {
	name := "native/custom/handler/worker/worker_test.go"
	path := "tmpl/sls/native/custom/handler/worker/worker_test.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.WorkerTestGo = content

	return nil
}
