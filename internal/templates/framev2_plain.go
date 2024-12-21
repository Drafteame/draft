package templates

type SLSFrameV2Plain struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          SLSFrameV2PlainHandler
}

type SLSFrameV2PlainHandler struct {
	BootstrapGo []byte
	HandlerGo   []byte
}

func loadFrameV2Plain(v *SLS, data any) error {
	loaders := []func(*SLS, any) error{
		loadFrameV2PlainMainGo,
		loadFrameV2PlainLambdaConfigYAML,
		loadFrameV2PlainHandler,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadFrameV2PlainMainGo(v *SLS, data any) error {
	name := "framev2/plain/main.go"
	path := "tmpl/sls/framev2/plain/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.FrameV2.Plain.MainGo = content

	return nil
}

func loadFrameV2PlainLambdaConfigYAML(v *SLS, data any) error {
	name := "framev2/plain/lambda-config.yml"
	path := "tmpl/sls/framev2/plain/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.FrameV2.Plain.LambdaConfigYAML = content

	return nil
}

func loadFrameV2PlainHandler(v *SLS, data any) error {
	if err := loadFrameV2PlainHandlerBoostrapGo(v, data); err != nil {
		return err
	}

	return loadFrameV2PlainHandlerHandlerGo(v, data)
}

func loadFrameV2PlainHandlerBoostrapGo(v *SLS, data any) error {
	name := "framev2/plain/handler/bootstrap.go"
	path := "tmpl/sls/framev2/plain/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.FrameV2.Plain.Handler.BootstrapGo = content

	return nil
}

func loadFrameV2PlainHandlerHandlerGo(v *SLS, data any) error {
	name := "framev2/plain/handler/handler.go"
	path := "tmpl/sls/framev2/plain/handler/handler.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.FrameV2.Plain.Handler.HandlerGo = content

	return nil
}
