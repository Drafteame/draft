package templates

type LambdaPlain struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          LambdaPlainHandler
}

type LambdaPlainHandler struct {
	BootstrapGo []byte
	HandlerGo   []byte
	WorkerGo    []byte
	ProviderGo  []byte
	EmbedYML    []byte
	ResourcesGo []byte
	DtosGo      []byte
}

func loadLambdaPlain(v PlainSetter, data any) error {
	plain := LambdaPlain{}

	loaders := []func(*LambdaPlain, any) error{
		loadLambdaPlainMainGo,
		loadLambdaPlainLambdaConfigYAML,
		loadLambdaPlainHandler,
	}

	for _, loader := range loaders {
		if err := loader(&plain, data); err != nil {
			return err
		}
	}

	v.SetPlain(plain)

	return nil
}

func loadLambdaPlainMainGo(v *LambdaPlain, data any) error {
	name := "native/plain/main.go"
	path := "tmpl/sls/native/plain/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.MainGo = content

	return nil
}

func loadLambdaPlainLambdaConfigYAML(v *LambdaPlain, data any) error {
	name := "native/plain/lambda-config.yml"
	path := "tmpl/sls/native/plain/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.LambdaConfigYAML = content

	return nil
}

func loadLambdaPlainHandler(v *LambdaPlain, data any) error {
	loaders := []func(*LambdaPlain, any) error{
		loadLambdaPlainHandlerBoostrapGo,
		loadLambdaPlainHandlerWorkerWorkerGo,
		loadLambdaPlainHandlerWorkerResourcesGo,
		loadLambdaPlainHandlerEmbedYaml,
		loadLambdaPlainHandlerDtosGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadLambdaPlainHandlerBoostrapGo(v *LambdaPlain, data any) error {
	name := "native/plain/handler/bootstrap.go"
	path := "tmpl/sls/native/plain/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.BootstrapGo = content

	return nil
}

func loadLambdaPlainHandlerWorkerWorkerGo(v *LambdaPlain, data any) error {
	name := "native/plain/handler/worker/worker.go"
	path := "tmpl/sls/native/plain/handler/worker/worker.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.WorkerGo = content

	return nil
}

func loadLambdaPlainHandlerWorkerResourcesGo(v *LambdaPlain, data any) error {
	name := "native/plain/handler/worker/resources.go"
	path := "tmpl/sls/native/plain/handler/worker/resources.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.ResourcesGo = content

	return nil
}

func loadLambdaPlainHandlerEmbedYaml(v *LambdaPlain, data any) error {
	name := "native/plain/handler/embed/embed.yaml"
	path := "tmpl/sls/native/plain/handler/embed/embed.yaml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.EmbedYML = content

	return nil
}

func loadLambdaPlainHandlerDtosGo(v *LambdaPlain, data any) error {
	name := "native/plain/handler/dtos/dto.go"
	path := "tmpl/sls/native/plain/handler/dtos/dto.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.DtosGo = content

	return nil
}
