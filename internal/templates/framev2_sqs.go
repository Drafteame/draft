package templates

type SSLSFrameV2Sqs struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          SLSFrameV2SqsHandler
}

type SLSFrameV2SqsHandler struct {
	BootstrapGo []byte
	HandlerGo   []byte
	WorkerGo    []byte
}

func loadFrameV2Sqs(v *SLS, data any) error {
	loaders := []func(*SLS, any) error{
		loadFrameV2SqsMainGo,
		loadFrameV2SqsLambdaConfigYAML,
		loadFrameV2SqsHandler,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadFrameV2SqsMainGo(v *SLS, data any) error {
	name := "framev2/sqs/main.go"
	path := "tmpl/sls/framev2/sqs/main.go.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.Sqs.MainGo = content

	return nil
}

func loadFrameV2SqsLambdaConfigYAML(v *SLS, data any) error {
	name := "framev2/sqs/lambda-config.yml"
	path := "tmpl/sls/framev2/sqs/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.Sqs.LambdaConfigYAML = content

	return nil
}

func loadFrameV2SqsHandler(v *SLS, data any) error {
	loaders := []func(*SLS, any) error{
		loadFrameV2SqsHandlerBoostrapGo,
		loadFrameV2SqsHandlerHandlerGo,
		loadFrameV2SqsHandlerWorkerGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadFrameV2SqsHandlerBoostrapGo(v *SLS, data any) error {
	name := "framev2/sqs/handler/bootstrap.go"
	path := "tmpl/sls/framev2/sqs/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.Sqs.Handler.BootstrapGo = content

	return nil
}

func loadFrameV2SqsHandlerHandlerGo(v *SLS, data any) error {
	name := "framev2/sqs/handler/handler.go"
	path := "tmpl/sls/framev2/sqs/handler/handler.go.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.Sqs.Handler.HandlerGo = content

	return nil
}

func loadFrameV2SqsHandlerWorkerGo(v *SLS, data any) error {
	name := "framev2/sqs/handler/worker.go"
	path := "tmpl/sls/framev2/sqs/handler/worker.go.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.Sqs.Handler.WorkerGo = content

	return nil
}
