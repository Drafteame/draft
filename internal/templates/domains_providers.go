package templates

import (
	"github.com/Drafteame/draft/internal/dtos"
)

type Providers struct {
	GeneratorGo []byte
	ProvidersGo []byte
	ServiceGo   []byte
	Postgres    ProvidersPostgres
}

func loadDomainsProviders(v *Providers, data any) error {
	loaders := []func(*Providers, any) error{
		loadProvidersGeneratorGo,
		loadProvidersServiceGo,
	}

	input, ok := data.(dtos.DomainInput)
	if ok && input.DBType == "postgres" {
		loaders = append(loaders, loadProvidersPostgres)
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
