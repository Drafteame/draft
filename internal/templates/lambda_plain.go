package templates

type LambdaPlain struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          LambdaPlainHandler
}

type LambdaPlainHandler struct {
	BootstrapGo []byte
	HandlerGo   []byte
	ProviderGo  []byte
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
	name := "framev2/plain/main.go"
	path := "tmpl/sls/framev2/plain/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.MainGo = content

	return nil
}

func loadLambdaPlainLambdaConfigYAML(v *LambdaPlain, data any) error {
	name := "framev2/plain/lambda-config.yml"
	path := "tmpl/sls/framev2/plain/lambda-config.yml.tmpl"

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
		loadLambdaPlainHandlerHandlerGo,
		loadLambdaPlainHandlerProviderGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadLambdaPlainHandlerBoostrapGo(v *LambdaPlain, data any) error {
	name := "framev2/plain/handler/bootstrap.go"
	path := "tmpl/sls/framev2/plain/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.BootstrapGo = content

	return nil
}

func loadLambdaPlainHandlerHandlerGo(v *LambdaPlain, data any) error {
	name := "framev2/plain/handler/handler.go"
	path := "tmpl/sls/framev2/plain/handler/handler.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.HandlerGo = content

	return nil
}

func loadLambdaPlainHandlerProviderGo(v *LambdaPlain, data any) error {
	name := "framev2/plain/handler/provider.go"
	path := "tmpl/sls/framev2/plain/handler/provider.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.ProviderGo = content

	return nil
}
