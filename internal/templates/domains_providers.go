package templates

import (
	appdata "github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
)

type Providers struct {
	GeneratorsNanoidTableid ProvidersGeneratorsNanoidTableid
}

func loadDomainsProviders(v *Providers, data any) error {
	// Only load providers for postgres
	input, ok := data.(dtos.DomainInput)
	if !ok || input.DBType != appdata.DBTypePostgres {
		return nil
	}

	return loadProvidersGeneratorsNanoidTableid(v, data)
}

func loadProvidersGeneratorsNanoidTableid(v *Providers, data any) error {
	return loadDomainProvidersGeneratorsNanoidTableid(&v.GeneratorsNanoidTableid, data)
}

type ProvidersGeneratorsNanoidTableid struct {
	ProvideGo []byte
}

func loadDomainProvidersGeneratorsNanoidTableid(v *ProvidersGeneratorsNanoidTableid, data any) error {
	loaders := []func(*ProvidersGeneratorsNanoidTableid, any) error{
		loadProvidersGeneratorsNanoidTableidProvideGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadProvidersGeneratorsNanoidTableidProvideGo(v *ProvidersGeneratorsNanoidTableid, data any) error {
	name := "domains/providers/generators/nanoid/tableid/provide.go"
	path := "tmpl/domain/providers/generators/nanoid/tableid/provide.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.ProvideGo = content

	return nil
}
