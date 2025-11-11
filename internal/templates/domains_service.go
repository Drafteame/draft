package templates

import (
	appdata "github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
)

type Service struct {
	Postgres ServicePostgres
	Dynamo   ServiceDynamo
}

func loadDomainsService(v *Service, data any) error {
	input, ok := data.(dtos.DomainInput)
	if !ok {
		return nil
	}

	var loaders []func(*Service, any) error

	switch input.DBType {
	case appdata.DBTypePostgres:
		loaders = append(loaders, loadServicePostgres)
	case appdata.DBTypeDynamo:
		loaders = append(loaders, loadServiceDynamo)
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadServicePostgres(v *Service, data any) error {
	return loadDomainsServicePostgres(&v.Postgres, data)
}

func loadServiceDynamo(v *Service, data any) error {
	return loadDomainsServiceDynamo(&v.Dynamo, data)
}

type ServicePostgres struct {
	CreateGo        []byte
	CreateTestGo    []byte
	DeleteGo        []byte
	DeleteTestGo    []byte
	GetGo           []byte
	GetTestGo       []byte
	InterfacesGo    []byte
	SearchGo        []byte
	SearchTestGo    []byte
	SearchOneGo     []byte
	SearchOneTestGo []byte
	ServiceGo       []byte
	ServiceTestGo   []byte
	UpdateGo        []byte
	UpdateTestGo    []byte
	ProvideGo       []byte
}

func loadDomainsServicePostgres(v *ServicePostgres, data any) error {
	loaders := []func(*ServicePostgres, any) error{
		loadServicePostgresCreateGo,
		loadServicePostgresCreateTestGo,
		loadServicePostgresDeleteGo,
		loadServicePostgresDeleteTestGo,
		loadServicePostgresGetGo,
		loadServicePostgresGetTestGo,
		loadServicePostgresInterfacesGo,
		loadServicePostgresSearchGo,
		loadServicePostgresSearchTestGo,
		loadServicePostgresSearchOneGo,
		loadServicePostgresSearchOneTestGo,
		loadServicePostgresServiceGo,
		loadServicePostgresServiceTestGo,
		loadServicePostgresUpdateGo,
		loadServicePostgresUpdateTestGo,
		loadServicePostgresProvideGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadServicePostgresCreateGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/create.go"
	path := "tmpl/domain/service/postgres/create.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.CreateGo = content

	return nil
}

func loadServicePostgresCreateTestGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/create_test.go"
	path := "tmpl/domain/service/postgres/create_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.CreateTestGo = content

	return nil
}

func loadServicePostgresDeleteGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/delete.go"
	path := "tmpl/domain/service/postgres/delete.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.DeleteGo = content

	return nil
}

func loadServicePostgresDeleteTestGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/delete_test.go"
	path := "tmpl/domain/service/postgres/delete_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.DeleteTestGo = content

	return nil
}

func loadServicePostgresGetGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/get.go"
	path := "tmpl/domain/service/postgres/get.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.GetGo = content

	return nil
}

func loadServicePostgresGetTestGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/get_test.go"
	path := "tmpl/domain/service/postgres/get_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.GetTestGo = content

	return nil
}

func loadServicePostgresInterfacesGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/interfaces.go"
	path := "tmpl/domain/service/postgres/interfaces.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.InterfacesGo = content

	return nil
}

func loadServicePostgresSearchGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/search.go"
	path := "tmpl/domain/service/postgres/search.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchGo = content

	return nil
}

func loadServicePostgresSearchTestGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/search_test.go"
	path := "tmpl/domain/service/postgres/search_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchTestGo = content

	return nil
}

func loadServicePostgresSearchOneGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/search_one.go"
	path := "tmpl/domain/service/postgres/search_one.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchOneGo = content

	return nil
}

func loadServicePostgresSearchOneTestGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/search_one_test.go"
	path := "tmpl/domain/service/postgres/search_one_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchOneTestGo = content

	return nil
}

func loadServicePostgresServiceGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/service.go"
	path := "tmpl/domain/service/postgres/service.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.ServiceGo = content

	return nil
}

func loadServicePostgresServiceTestGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/service_test.go"
	path := "tmpl/domain/service/postgres/service_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.ServiceTestGo = content

	return nil
}

func loadServicePostgresUpdateGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/update.go"
	path := "tmpl/domain/service/postgres/update.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.UpdateGo = content

	return nil
}

func loadServicePostgresUpdateTestGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/update_test.go"
	path := "tmpl/domain/service/postgres/update_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.UpdateTestGo = content

	return nil
}

func loadServicePostgresProvideGo(v *ServicePostgres, data any) error {
	name := "domains/service/postgres/provide.go"
	path := "tmpl/domain/service/postgres/provide.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.ProvideGo = content

	return nil
}

type ServiceDynamo struct {
	InterfacesGo []byte
	ServiceGo    []byte
	ProviderGo   []byte
}

func loadDomainsServiceDynamo(v *ServiceDynamo, data any) error {
	loaders := []func(*ServiceDynamo, any) error{
		loadServiceDynamoInterfacesGo,
		loadServiceDynamoServiceGo,
		loadServiceDynamoProviderGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadServiceDynamoInterfacesGo(v *ServiceDynamo, data any) error {
	name := "domains/service/dynamo/interfaces.go"
	path := "tmpl/domain/service/dynamo/interfaces.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.InterfacesGo = content

	return nil
}

func loadServiceDynamoServiceGo(v *ServiceDynamo, data any) error {
	name := "domains/service/dynamo/service.go"
	path := "tmpl/domain/service/dynamo/service.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.ServiceGo = content

	return nil
}

func loadServiceDynamoProviderGo(v *ServiceDynamo, data any) error {
	name := "domains/service/dynamo/provider.go"
	path := "tmpl/domain/service/dynamo/provider.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.ProviderGo = content

	return nil
}
