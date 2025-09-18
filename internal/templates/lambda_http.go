package templates

type LambdaHTTP struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          LambdaHTTPHandler
}

type LambdaHTTPHandler struct {
	BootstrapGo []byte
	HandlerGo   []byte
	WorkerGo    []byte
	EmbedYML    []byte
	ResourcesGo []byte
	DtosGo      []byte
}

func loadLambdaHTTP(v HTTPSetter, data any) error {
	http := LambdaHTTP{}

	loaders := []func(*LambdaHTTP, any) error{
		loadLambdaHTTPMainGo,
		loadLambdaHTTPLambdaConfigYAML,
		loadLambdaHTTPHandler,
	}

	for _, loader := range loaders {
		if err := loader(&http, data); err != nil {
			return err
		}
	}

	v.SetHTTP(http)

	return nil
}

func loadLambdaHTTPMainGo(v *LambdaHTTP, data any) error {
	name := "native/http/main.go"
	path := "tmpl/sls/native/http/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.MainGo = content

	return nil
}

func loadLambdaHTTPLambdaConfigYAML(v *LambdaHTTP, data any) error {
	name := "native/http/lambda-config.yml"
	path := "tmpl/sls/native/http/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.LambdaConfigYAML = content

	return nil
}

func loadLambdaHTTPHandler(v *LambdaHTTP, data any) error {
	loaders := []func(*LambdaHTTP, any) error{
		loadLambdaHTTPHandlerBootstrapGo,
		loadLambdaHTTPHandlerWorkerWorkerGo,
		loadLambdaHTTPHandlerWorkerResourcesGo,
		loadLambdaHTTPHandlerEmbedYaml,
		loadLambdaHTTPHandlerDtosGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadLambdaHTTPHandlerBootstrapGo(v *LambdaHTTP, data any) error {
	name := "native/http/handler/bootstrap.go"
	path := "tmpl/sls/native/http/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.BootstrapGo = content

	return nil
}

func loadLambdaHTTPHandlerWorkerWorkerGo(v *LambdaHTTP, data any) error {
	name := "native/http/handler/worker/worker.go"
	path := "tmpl/sls/native/http/handler/worker/worker.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.WorkerGo = content

	return nil
}

func loadLambdaHTTPHandlerWorkerResourcesGo(v *LambdaHTTP, data any) error {
	name := "native/http/handler/worker/resources.go"
	path := "tmpl/sls/native/http/handler/worker/resources.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.ResourcesGo = content

	return nil
}

func loadLambdaHTTPHandlerEmbedYaml(v *LambdaHTTP, data any) error {
	name := "native/http/handler/embed/embed.yaml"
	path := "tmpl/sls/native/http/handler/embed/embed.yaml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.EmbedYML = content

	return nil
}

func loadLambdaHTTPHandlerDtosGo(v *LambdaHTTP, data any) error {
	name := "native/http/handler/dtos/dto.go"
	path := "tmpl/sls/native/http/handler/dtos/dto.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.DtosGo = content

	return nil
}
