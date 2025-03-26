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
}

func loadLambdaCron(v CronSetter, data any) error {
	cron := LambdaCron{}

	loaders := []func(cron *LambdaCron, data any) error{
		loadLambdaCronMainGo,
		loadFrameV2CronLambdaConfigYAML,
		loadFrameV2CronHandler,
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
	name := "framev2/cron/main.go"
	path := "tmpl/sls/framev2/cron/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.MainGo = content

	return nil
}

func loadFrameV2CronLambdaConfigYAML(v *LambdaCron, data any) error {
	name := "framev2/cron/lambda-config.yml"
	path := "tmpl/sls/framev2/cron/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.LambdaConfigYAML = content

	return nil
}

func loadFrameV2CronHandler(v *LambdaCron, data any) error {
	loaders := []func(cron *LambdaCron, data any) error{
		loadFrameV2CronHandlerProviderGo,
		loadFrameV2CronHandlerBootstrapGo,
		loadFrameV2CronHandlerHandlerGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadFrameV2CronHandlerBootstrapGo(v *LambdaCron, data any) error {
	name := "framev2/cron/handler/bootstrap.go"
	path := "tmpl/sls/framev2/cron/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.BootstrapGo = content

	return nil
}

func loadFrameV2CronHandlerHandlerGo(v *LambdaCron, data any) error {
	name := "framev2/cron/handler/handler.go"
	path := "tmpl/sls/framev2/cron/handler/handler.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.HandlerGo = content

	return nil
}

func loadFrameV2CronHandlerProviderGo(v *LambdaCron, data any) error {
	name := "framev2/cron/handler/provider.go"
	path := "tmpl/sls/framev2/cron/handler/provider.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.ProviderGo = content

	return nil
}
