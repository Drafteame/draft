package templates

import (
	"github.com/Drafteame/draft/internal/dtos"
)

type LambdaSnsSqs struct {
	MainGo           []byte
	LambdaConfigYAML []byte
	Handler          LambdaSnsSqsHandler
}

type LambdaSnsSqsHandler struct {
	BootstrapGo []byte
	HandlerGo   []byte
	WorkerGo    []byte
	ProviderGo  []byte
	EmbedYML    []byte
	ResourcesGo []byte
}

func loadLambdaSnsSqs(v SnsSqsSetter, data dtos.LambdaInput) error {
	snssqs := LambdaSnsSqs{}

	var loaders []func(*LambdaSnsSqs, any) error

	if !data.WithFrame {
		loaders = []func(*LambdaSnsSqs, any) error{
			nativeLoadLambdaSnsSqsMainGo,
			nativeLoadLambdaSnsSqsLambdaConfigYAML,
			nativeLoadLambdaSnsSqsHandler,
		}
	} else {
		loaders = []func(*LambdaSnsSqs, any) error{
			loadLambdaSnsSqsMainGo,
			loadLambdaSnsSqsLambdaConfigYAML,
			loadLambdaSnsSqsHandler,
		}
	}

	for _, loader := range loaders {
		if err := loader(&snssqs, data); err != nil {
			return err
		}
	}

	v.SetSnsSqs(snssqs)

	return nil
}

func loadLambdaSnsSqsMainGo(v *LambdaSnsSqs, data any) error {
	name := "framev2/snssqs/main.go"
	path := "tmpl/sls/framev2/snssqs/main.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.MainGo = content

	return nil
}

func loadLambdaSnsSqsLambdaConfigYAML(v *LambdaSnsSqs, data any) error {
	name := "framev2/snssqs/lambda-config.yml"
	path := "tmpl/sls/framev2/snssqs/lambda-config.yml.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.LambdaConfigYAML = content

	return nil
}

func loadLambdaSnsSqsHandler(v *LambdaSnsSqs, data any) error {
	loaders := []func(*LambdaSnsSqs, any) error{
		loadLambdaSnsSqsHandlerBootstrapGo,
		loadLambdaSnsSqsHandlerHandlerGo,
		loadLambdaSnsSqsHandlerWorkerGo,
		loadLambdaSnsSqsHandlerProviderGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadLambdaSnsSqsHandlerBootstrapGo(v *LambdaSnsSqs, data any) error {
	name := "framev2/snssqs/handler/bootstrap.go"
	path := "tmpl/sls/framev2/snssqs/handler/bootstrap.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.BootstrapGo = content

	return nil
}

func loadLambdaSnsSqsHandlerHandlerGo(v *LambdaSnsSqs, data any) error {
	name := "framev2/snssqs/handler/handler.go"
	path := "tmpl/sls/framev2/snssqs/handler/handler.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.HandlerGo = content

	return nil
}

func loadLambdaSnsSqsHandlerWorkerGo(v *LambdaSnsSqs, data any) error {
	name := "framev2/snssqs/handler/worker.go"
	path := "tmpl/sls/framev2/snssqs/handler/worker.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.WorkerGo = content

	return nil
}

func loadLambdaSnsSqsHandlerProviderGo(v *LambdaSnsSqs, data any) error {
	name := "framev2/snssqs/handler/provider.go"
	path := "tmpl/sls/framev2/snssqs/handler/provider.go.tmpl"

	content, err := loadTemplate(name, path, data, sls)
	if err != nil {
		return err
	}

	v.Handler.ProviderGo = content

	return nil
}
