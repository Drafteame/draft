package templates

type SLSFrameV2SnsSqs struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          SLSFrameV2SnsSqsHandler
}

type SLSFrameV2SnsSqsHandler struct {
	BootstrapGo []byte
	HandlerGo   []byte
	WorkerGo    []byte
}

func loadFrameV2SnsSqs(v *SLS, data any) error {
	loaders := []func(*SLS, any) error{
		loadFrameV2SnsSqsMainGo,
		loadFrameV2SnsSqsLambdaConfigYAML,
		loadFrameV2SnsSqsHandler,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadFrameV2SnsSqsMainGo(v *SLS, data any) error {
	name := "framev2/snssqs/main.go"
	path := "tmpl/sls/framev2/snssqs/main.go.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.SnsSqs.MainGo = content

	return nil
}

func loadFrameV2SnsSqsLambdaConfigYAML(v *SLS, data any) error {
	name := "framev2/snssqs/lambda-config.yml"
	path := "tmpl/sls/framev2/snssqs/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.SnsSqs.LambdaConfigYAML = content

	return nil
}

func loadFrameV2SnsSqsHandler(v *SLS, data any) error {
	loaders := []func(*SLS, any) error{
		loadFrameV2SnsSqsHandlerBootstrapGo,
		loadFrameV2SnsSqsHandlerHandlerGo,
		loadFrameV2SnsSqsHandlerWorkerGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadFrameV2SnsSqsHandlerBootstrapGo(v *SLS, data any) error {
	name := "framev2/snssqs/handler/bootstrap.go"
	path := "tmpl/sls/framev2/snssqs/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.SnsSqs.Handler.BootstrapGo = content

	return nil
}

func loadFrameV2SnsSqsHandlerHandlerGo(v *SLS, data any) error {
	name := "framev2/snssqs/handler/handler.go"
	path := "tmpl/sls/framev2/snssqs/handler/handler.go.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.SnsSqs.Handler.HandlerGo = content

	return nil
}

func loadFrameV2SnsSqsHandlerWorkerGo(v *SLS, data any) error {
	name := "framev2/snssqs/handler/worker.go"
	path := "tmpl/sls/framev2/snssqs/handler/worker.go.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.SnsSqs.Handler.WorkerGo = content

	return nil
}
