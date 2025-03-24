package templates

type LambdaHTTP struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          LambdaHTTPHandler
}

type LambdaHTTPHandler struct {
	BootstrapGo []byte
	HandlerGo   []byte
	ProviderGo  []byte
	RouteGo     []byte
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
	name := "framev2/http/main.go"
	path := "tmpl/sls/framev2/http/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.MainGo = content

	return nil
}

func loadLambdaHTTPLambdaConfigYAML(v *LambdaHTTP, data any) error {
	name := "framev2/http/lambda-config.yml"
	path := "tmpl/sls/framev2/http/lambda-config.yml.tmpl"

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
		loadLambdaHTTPHandlerHandlerGo,
		loadLambdaHTTPHandlerProviderGo,
		loadLambdaHTTPHandlerRouteGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadLambdaHTTPHandlerBootstrapGo(v *LambdaHTTP, data any) error {
	name := "framev2/http/handler/bootstrap.go"
	path := "tmpl/sls/framev2/http/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.BootstrapGo = content

	return nil
}

func loadLambdaHTTPHandlerHandlerGo(v *LambdaHTTP, data any) error {
	name := "framev2/http/handler/handler.go"
	path := "tmpl/sls/framev2/http/handler/handler.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.HandlerGo = content

	return nil
}

func loadLambdaHTTPHandlerProviderGo(v *LambdaHTTP, data any) error {
	name := "framev2/http/handler/provider.go"
	path := "tmpl/sls/framev2/http/handler/provider.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.ProviderGo = content

	return nil
}

func loadLambdaHTTPHandlerRouteGo(v *LambdaHTTP, data any) error {
	name := "framev2/http/handler/route.go"
	path := "tmpl/sls/framev2/http/handler/route.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.RouteGo = content

	return nil
}
