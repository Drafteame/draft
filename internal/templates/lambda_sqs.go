package templates

type LambdaSqs struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          LambdaSqsHandler
}

type LambdaSqsHandler struct {
	BootstrapGo []byte
	HandlerGo   []byte
	WorkerGo    []byte
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
	name := "framev2/sqs/main.go"
	path := "tmpl/sls/framev2/sqs/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.MainGo = content

	return nil
}

func loadLambdaSqsLambdaConfigYAML(v *LambdaSqs, data any) error {
	name := "framev2/sqs/lambda-config.yml"
	path := "tmpl/sls/framev2/sqs/lambda-config.yml.tmpl"

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
		loadLambdaSqsHandlerHandlerGo,
		loadLambdaSqsHandlerWorkerGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadLambdaSqsHandlerBoostrapGo(v *LambdaSqs, data any) error {
	name := "framev2/sqs/handler/bootstrap.go"
	path := "tmpl/sls/framev2/sqs/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.BootstrapGo = content

	return nil
}

func loadLambdaSqsHandlerHandlerGo(v *LambdaSqs, data any) error {
	name := "framev2/sqs/handler/handler.go"
	path := "tmpl/sls/framev2/sqs/handler/handler.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.HandlerGo = content

	return nil
}

func loadLambdaSqsHandlerWorkerGo(v *LambdaSqs, data any) error {
	name := "framev2/sqs/handler/worker.go"
	path := "tmpl/sls/framev2/sqs/handler/worker.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.WorkerGo = content

	return nil
}
