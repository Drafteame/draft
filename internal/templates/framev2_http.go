package templates

type SLSFrameV2HTTP struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          SLSFrameV2HTTPHandler
}

type SLSFrameV2HTTPHandler struct {
	BootstrapGo []byte
	HandlerGo   []byte
}

func loadFrameV2HTTP(v *SLS, data any) error {
	loaders := []func(*SLS, any) error{
		loadFrameV2HTTPMainGo,
		loadFrameV2HTTPLambdaConfigYAML,
		loadFrameV2HTTPHandler,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadFrameV2HTTPMainGo(v *SLS, data any) error {
	name := "framev2/http/main.go"
	path := "tmpl/sls/framev2/http/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.FrameV2.HTTP.MainGo = content

	return nil
}

func loadFrameV2HTTPLambdaConfigYAML(v *SLS, data any) error {
	name := "framev2/http/lambda-config.yml"
	path := "tmpl/sls/framev2/http/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.FrameV2.HTTP.LambdaConfigYAML = content

	return nil
}

func loadFrameV2HTTPHandler(v *SLS, data any) error {
	if err := loadFrameV2HTTPHandlerBootstrapGo(v, data); err != nil {
		return err
	}

	return loadFrameV2HTTPHandlerHandlerGo(v, data)
}

func loadFrameV2HTTPHandlerBootstrapGo(v *SLS, data any) error {
	name := "framev2/http/handler/bootstrap.go"
	path := "tmpl/sls/framev2/http/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.FrameV2.HTTP.Handler.BootstrapGo = content

	return nil
}

func loadFrameV2HTTPHandlerHandlerGo(v *SLS, data any) error {
	name := "framev2/http/handler/handler.go"
	path := "tmpl/sls/framev2/http/handler/handler.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.FrameV2.HTTP.Handler.HandlerGo = content

	return nil
}
