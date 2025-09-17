package templates

type LambdaCron struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          LambdaCronHandler
}

type LambdaCronHandler struct {
	BootstrapGo []byte
	HandlerGo   []byte
	ProviderGo  []byte
	WorkerGo    []byte
	EmbedYML    []byte
	ResourcesGo []byte
}

func loadLambdaCron(v CronSetter, data any) error {
	cron := LambdaCron{}

	loaders := []func(cron *LambdaCron, data any) error{
		loadLambdaCronMainGo,
		loadLambdaCronConfigYAML,
		loadLambdaCronHandler,
	}

	for _, loader := range loaders {
		if err := loader(&cron, data); err != nil {
			return err
		}
	}

	v.SetCron(cron)

	return nil
}

func loadLambdaCronMainGo(v *LambdaCron, data any) error {
	name := "native/cron/main.go"
	path := "tmpl/sls/native/cron/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.MainGo = content

	return nil
}

func loadLambdaCronConfigYAML(v *LambdaCron, data any) error {
	name := "native/cron/lambda-config.yml"
	path := "tmpl/sls/native/cron/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.LambdaConfigYAML = content

	return nil
}

func loadLambdaCronHandler(v *LambdaCron, data any) error {
	loaders := []func(*LambdaCron, any) error{
		loadLambdaCronHandlerBoostrapGo,
		loadLambdaCronHandlerWorkerWorkerGo,
		loadLambdaCronHandlerWorkerResourcesGo,
		loadLambdaCronHandlerEmbedYaml,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadLambdaCronHandlerBoostrapGo(v *LambdaCron, data any) error {
	name := "native/cron/handler/bootstrap.go"
	path := "tmpl/sls/native/cron/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.BootstrapGo = content

	return nil
}

func loadLambdaCronHandlerWorkerWorkerGo(v *LambdaCron, data any) error {
	name := "native/cron/handler/worker/worker.go"
	path := "tmpl/sls/native/cron/handler/worker/worker.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.WorkerGo = content

	return nil
}

func loadLambdaCronHandlerWorkerResourcesGo(v *LambdaCron, data any) error {
	name := "native/cron/handler/worker/resources.go"
	path := "tmpl/sls/native/cron/handler/worker/resources.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.ResourcesGo = content

	return nil
}

func loadLambdaCronHandlerEmbedYaml(v *LambdaCron, data any) error {
	name := "native/cron/handler/embed/embed.yaml"
	path := "tmpl/sls/native/cron/handler/embed/embed.yaml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.EmbedYML = content

	return nil
}
