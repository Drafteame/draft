package templates

type LambdaSnsSqs struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          LambdaSnsSqsHandler
}

type LambdaSnsSqsHandler struct {
	BootstrapGo   []byte
	HandlerGo     []byte
	WorkerGo      []byte
	ProviderGo    []byte
	EmbedYML      []byte
	ResourcesGo   []byte
	DtosGo        []byte
	IdempotencyGo []byte
	InterfacesGo  []byte
}

func loadLambdaSnsSqs(v SnsSqsSetter, data any) error {
	snssqs := LambdaSnsSqs{}

	loaders := []func(*LambdaSnsSqs, any) error{
		loadLambdaSnsSqsMainGo,
		loadLambdaSnsSqsLambdaConfigYAML,
		loadLambdaSnsSqsHandler,
	}

	for _, loader := range loaders {
		if err := loader(&snssqs, data); err != nil {
			return err
		}
	}

	v.SetSnsSqs(snssqs)

	return nil
}

func loadLambdaSnsSqsMainGo(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/main.go"
	path := "tmpl/sls/native/snssqs/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.MainGo = content

	return nil
}

func loadLambdaSnsSqsLambdaConfigYAML(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/lambda-config.yml"
	path := "tmpl/sls/native/snssqs/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.LambdaConfigYAML = content

	return nil
}

func loadLambdaSnsSqsHandler(v *LambdaSnsSqs, data any) error {
	loaders := []func(*LambdaSnsSqs, any) error{
		loadLambdaSnsSqsHandlerBoostrapGo,
		loadLambdaSnsSqsHandlerWorkerWorkerGo,
		loadLambdaSnsSqsHandlerWorkerResourcesGo,
		loadLambdaSnsSqsHandlerEmbedYaml,
		loadLambdaSnsSqsHandlerDtosGo,
		loadLambdaSnsSqsHandlerIdempotencyGo,
		loadLambdaSnsSqsHandlerInterfacesGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadLambdaSnsSqsHandlerBoostrapGo(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/handler/bootstrap.go"
	path := "tmpl/sls/native/snssqs/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.BootstrapGo = content

	return nil
}

func loadLambdaSnsSqsHandlerWorkerWorkerGo(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/handler/worker/worker.go"
	path := "tmpl/sls/native/snssqs/handler/worker/worker.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.WorkerGo = content

	return nil
}

func loadLambdaSnsSqsHandlerWorkerResourcesGo(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/handler/worker/resources.go"
	path := "tmpl/sls/native/snssqs/handler/worker/resources.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.ResourcesGo = content

	return nil
}

func loadLambdaSnsSqsHandlerEmbedYaml(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/handler/embed/embed.yaml"
	path := "tmpl/sls/native/snssqs/handler/embed/embed.yaml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.EmbedYML = content

	return nil
}

func loadLambdaSnsSqsHandlerDtosGo(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/handler/dtos/dto.go"
	path := "tmpl/sls/native/snssqs/handler/dtos/dto.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.DtosGo = content

	return nil
}

func loadLambdaSnsSqsHandlerIdempotencyGo(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/handler/worker/idempotency.go"
	path := "tmpl/sls/native/snssqs/handler/worker/idempotency.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.IdempotencyGo = content

	return nil
}

func loadLambdaSnsSqsHandlerInterfacesGo(v *LambdaSnsSqs, data any) error {
	name := "native/snssqs/handler/worker/interfaces.go"
	path := "tmpl/sls/native/snssqs/handler/worker/interfaces.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.InterfacesGo = content

	return nil
}
