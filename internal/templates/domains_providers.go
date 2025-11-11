package templates

import (
	appdata "github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
)

type Providers struct {
	GeneratorGo             []byte
	ServiceGo               []byte
	Postgres                ProvidersPostgres
	GeneratorsNanoidTableid ProvidersGeneratorsNanoidTableid
}

func loadDomainsProviders(v *Providers, data any) error {
	var loaders []func(*Providers, any) error

	// Only load providers for postgres
	input, ok := data.(dtos.DomainInput)
	if ok && input.DBType == appdata.DBTypePostgres {
		loaders = []func(*Providers, any) error{
			loadProvidersGeneratorGo,
			loadProvidersServiceGo,
			loadProvidersPostgres,
			loadProvidersGeneratorsNanoidTableid,
		}
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadProvidersGeneratorGo(v *Providers, data any) error {
	name := "domains/providers/generator.go"
	path := "tmpl/domain/providers/generator.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.GeneratorGo = content

	return nil
}

func loadProvidersServiceGo(v *Providers, data any) error {
	name := "domains/providers/service.go"
	path := "tmpl/domain/providers/service.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.ServiceGo = content

	return nil
}

func loadProvidersPostgres(v *Providers, data any) error {
	return loadDomainProvidersPostgres(&v.Postgres, data)
}

type ProvidersPostgres struct {
	RepositoryGo []byte
}

func loadDomainProvidersPostgres(v *ProvidersPostgres, data any) error {
	loaders := []func(*ProvidersPostgres, any) error{
		loadProvidersPostgresRepositoryGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadProvidersPostgresRepositoryGo(v *ProvidersPostgres, data any) error {
	name := "domains/providers/postgres/repository.go"
	path := "tmpl/domain/providers/postgres/repository.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.RepositoryGo = content

	return nil
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
