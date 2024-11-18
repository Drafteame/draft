package templates

type SLSFrameV2Cron struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          SLSFrameV2CronHandler
}

type SLSFrameV2CronHandler struct {
	BootstrapGo []byte
	HandlerGo   []byte
}

func loadFrameV2Cron(v *SLS, data any) error {
	loaders := []func(*SLS, any) error{
		loadFrameV2CronMainGo,
		loadFrameV2CronLambdaConfigYAML,
		loadFrameV2CronHandler,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadFrameV2CronMainGo(v *SLS, data any) error {
	name := "framev2/cron/main.go"
	path := "tmpl/sls/framev2/cron/main.go.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.Cron.MainGo = content

	return nil
}

func loadFrameV2CronLambdaConfigYAML(v *SLS, data any) error {
	name := "framev2/cron/lambda-config.yml"
	path := "tmpl/sls/framev2/cron/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.Cron.LambdaConfigYAML = content

	return nil
}

func loadFrameV2CronHandler(v *SLS, data any) error {
	if err := loadFrameV2CronHandlerBootstrapGo(v, data); err != nil {
		return err
	}

	return loadFrameV2CronHandlerHandlerGo(v, data)
}

func loadFrameV2CronHandlerBootstrapGo(v *SLS, data any) error {
	name := "framev2/cron/handler/bootstrap.go"
	path := "tmpl/sls/framev2/cron/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.Cron.Handler.BootstrapGo = content

	return nil
}

func loadFrameV2CronHandlerHandlerGo(v *SLS, data any) error {
	name := "framev2/cron/handler/handler.go"
	path := "tmpl/sls/framev2/cron/handler/handler.go.tmpl"

	content, err := loadTemplate(name, path, data)
	if err != nil {
		return err
	}

	v.FrameV2.Cron.Handler.HandlerGo = content

	return nil
}
