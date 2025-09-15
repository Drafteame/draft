package templates

type LambdaSqs struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          LambdaSqsHandler
}

type LambdaSqsHandler struct {
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

func loadLambdaSqs(v SqsSetter, data any) error {
	sqs := LambdaSqs{}

	loaders := []func(*LambdaSqs, any) error{
		loadLambdaSqsMainGo,
		loadLambdaSqsLambdaConfigYAML,
		loadLambdaSqsHandler,
	}

	for _, loader := range loaders {
		if err := loader(&sqs, data); err != nil {
			return err
		}
	}

	v.SetSqs(sqs)

	return nil
}

func loadLambdaSqsMainGo(v *LambdaSqs, data any) error {
	name := "native/sqs/main.go"
	path := "tmpl/sls/native/sqs/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.MainGo = content

	return nil
}

func loadLambdaSqsLambdaConfigYAML(v *LambdaSqs, data any) error {
	name := "native/sqs/lambda-config.yml"
	path := "tmpl/sls/native/sqs/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.LambdaConfigYAML = content

	return nil
}

func loadLambdaSqsHandler(v *LambdaSqs, data any) error {
	loaders := []func(*LambdaSqs, any) error{
		loadLambdaSqsHandlerBoostrapGo,
		loadLambdaSqsHandlerWorkerWorkerGo,
		loadLambdaSqsHandlerWorkerResourcesGo,
		loadLambdaSqsHandlerEmbedYaml,
		loadLambdaSqsHandlerDtosGo,
		loadLambdaSqsHandlerIdempotencyGo,
		loadLambdaSqsHandlerInterfacesGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadLambdaSqsHandlerBoostrapGo(v *LambdaSqs, data any) error {
	name := "native/sqs/handler/bootstrap.go"
	path := "tmpl/sls/native/sqs/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.BootstrapGo = content

	return nil
}

func loadLambdaSqsHandlerWorkerWorkerGo(v *LambdaSqs, data any) error {
	name := "native/sqs/handler/worker/worker.go"
	path := "tmpl/sls/native/sqs/handler/worker/worker.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.WorkerGo = content

	return nil
}

func loadLambdaSqsHandlerWorkerResourcesGo(v *LambdaSqs, data any) error {
	name := "native/sqs/handler/worker/resources.go"
	path := "tmpl/sls/native/sqs/handler/worker/resources.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.ResourcesGo = content

	return nil
}

func loadLambdaSqsHandlerEmbedYaml(v *LambdaSqs, data any) error {
	name := "native/sqs/handler/embed/embed.yaml"
	path := "tmpl/sls/native/sqs/handler/embed/embed.yaml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.EmbedYML = content

	return nil
}

func loadLambdaSqsHandlerDtosGo(v *LambdaSqs, data any) error {
	name := "native/sqs/handler/dtos/dto.go"
	path := "tmpl/sls/native/sqs/handler/dtos/dto.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.DtosGo = content

	return nil
}

func loadLambdaSqsHandlerIdempotencyGo(v *LambdaSqs, data any) error {
	name := "native/sqs/handler/worker/idempotency.go"
	path := "tmpl/sls/native/sqs/handler/worker/idempotency.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.IdempotencyGo = content

	return nil
}

func loadLambdaSqsHandlerInterfacesGo(v *LambdaSqs, data any) error {
	name := "native/sqs/handler/worker/interfaces.go"
	path := "tmpl/sls/native/sqs/handler/worker/interfaces.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.InterfacesGo = content

	return nil
}
