package templates

import (
	appdata "github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
)

type Repository struct {
	Postgres RepositoryPostgres
	Dynamo   RepositoryDynamo
}

func loadDomainsRepository(v *Repository, data any) error {
	input, ok := data.(dtos.DomainInput)
	if !ok {
		return nil
	}

	var loaders []func(*Repository, any) error

	switch input.DBType {
	case appdata.DBTypePostgres:
		loaders = append(loaders, loadRepositoryPostgres)
	case appdata.DBTypeDynamo:
		loaders = append(loaders, loadRepositoryDynamo)
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadRepositoryPostgres(v *Repository, data any) error {
	return loadDomainsRepositoryPostgres(&v.Postgres, data)
}

func loadRepositoryDynamo(v *Repository, data any) error {
	return loadDomainsRepositoryDynamo(&v.Dynamo, data)
}

type RepositoryPostgres struct {
	DotMockeryPkgYml []byte
	CreateGo         []byte
	CreateTestGo     []byte
	DeleteGo         []byte
	DeleteTestGo     []byte
	GetGo            []byte
	GetTestGo        []byte
	InterfacesGo     []byte
	RepositoryGo     []byte
	RepositoryTestGo []byte
	SearchGo         []byte
	SearchTestGo     []byte
	SearchOneGo      []byte
	SearchOneTestGo  []byte
	UpdateGo         []byte
	UpdateTestGo     []byte
	ProvideGo        []byte
	Builders         RepositoryPostgresBuilders
	Daos             RepositoryPostgresDaos
}

func loadDomainsRepositoryPostgres(v *RepositoryPostgres, data any) error {
	loaders := []func(*RepositoryPostgres, any) error{
		loadRepositoryPostgresDotMockeryPkgYml,
		loadRepositoryPostgresCreateGo,
		loadRepositoryPostgresCreateGoTest,
		loadRepositoryPostgresDeleteGo,
		loadRepositoryPostgresDeleteGoTest,
		loadRepositoryPostgresGetGo,
		loadRepositoryPostgresGetGoTest,
		loadRepositoryPostgresInterfacesGo,
		loadRepositoryPostgresRepositoryGo,
		loadRepositoryPostgresRepositoryGoTest,
		loadRepositoryPostgresSearchGo,
		loadRepositoryPostgresSearchGoTest,
		loadRepositoryPostgresSearchOneGo,
		loadRepositoryPostgresSearchOneGoTest,
		loadRepositoryPostgresUpdateGo,
		loadRepositoryPostgresUpdateGoTest,
		loadRepositoryPostgresProvideGo,
		loadRepositoryPostgresBuilders,
		loadRepositoryPostgresDaos,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadRepositoryPostgresDotMockeryPkgYml(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/.mockery.pkg.yml"
	path := "tmpl/domain/repository/postgres/.mockery.pkg.yml.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.DotMockeryPkgYml = content

	return nil
}

func loadRepositoryPostgresCreateGo(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/create.go"
	path := "tmpl/domain/repository/postgres/create.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.CreateGo = content

	return nil
}

func loadRepositoryPostgresCreateGoTest(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/create_test.go"
	path := "tmpl/domain/repository/postgres/create_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.CreateTestGo = content

	return nil
}

func loadRepositoryPostgresDeleteGo(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/delete.go"
	path := "tmpl/domain/repository/postgres/delete.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.DeleteGo = content

	return nil
}

func loadRepositoryPostgresDeleteGoTest(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/delete_test.go"
	path := "tmpl/domain/repository/postgres/delete_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.DeleteTestGo = content

	return nil
}

func loadRepositoryPostgresGetGo(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/get.go"
	path := "tmpl/domain/repository/postgres/get.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.GetGo = content

	return nil
}

func loadRepositoryPostgresGetGoTest(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/get_test.go"
	path := "tmpl/domain/repository/postgres/get_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.GetTestGo = content

	return nil
}

func loadRepositoryPostgresInterfacesGo(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/interfaces.go"
	path := "tmpl/domain/repository/postgres/interfaces.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.InterfacesGo = content

	return nil
}

func loadRepositoryPostgresRepositoryGo(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/repository.go"
	path := "tmpl/domain/repository/postgres/repository.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.RepositoryGo = content

	return nil
}

func loadRepositoryPostgresRepositoryGoTest(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/repository_test.go"
	path := "tmpl/domain/repository/postgres/repository_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.RepositoryTestGo = content

	return nil
}

func loadRepositoryPostgresSearchGo(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/search.go"
	path := "tmpl/domain/repository/postgres/search.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchGo = content

	return nil
}

func loadRepositoryPostgresSearchGoTest(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/search_test.go"
	path := "tmpl/domain/repository/postgres/search_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchTestGo = content

	return nil
}

func loadRepositoryPostgresSearchOneGo(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/search_one.go"
	path := "tmpl/domain/repository/postgres/search_one.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchOneGo = content

	return nil
}

func loadRepositoryPostgresSearchOneGoTest(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/search_one_test.go"
	path := "tmpl/domain/repository/postgres/search_one_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchOneTestGo = content

	return nil
}

func loadRepositoryPostgresUpdateGo(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/update.go"
	path := "tmpl/domain/repository/postgres/update.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.UpdateGo = content

	return nil
}

func loadRepositoryPostgresUpdateGoTest(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/update_test.go"
	path := "tmpl/domain/repository/postgres/update_test.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.UpdateTestGo = content

	return nil
}

func loadRepositoryPostgresProvideGo(v *RepositoryPostgres, data any) error {
	name := "domains/repository/postgres/provide.go"
	path := "tmpl/domain/repository/postgres/provide.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.ProvideGo = content

	return nil
}

func loadRepositoryPostgresBuilders(v *RepositoryPostgres, data any) error {
	return loadDomainsRepositoryPostgresBuilders(&v.Builders, data)
}

func loadRepositoryPostgresDaos(v *RepositoryPostgres, data any) error {
	return loadDomainsRepositoryPostgresDaos(&v.Daos, data)
}

type RepositoryPostgresBuilders struct {
	SearchGo           []byte
	SearchFiltersGo    []byte
	SearchOrdersGo     []byte
	SearchPaginationGo []byte
}

func loadDomainsRepositoryPostgresBuilders(v *RepositoryPostgresBuilders, data any) error {
	loaders := []func(*RepositoryPostgresBuilders, any) error{
		loadRepositoryPostgresBuildersSearchGo,
		loadRepositoryPostgresBuildersSearchFiltersGo,
		loadRepositoryPostgresBuildersSearchOrdersGo,
		loadRepositoryPostgresBuildersSearchPaginationGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadRepositoryPostgresBuildersSearchGo(v *RepositoryPostgresBuilders, data any) error {
	name := "domains/repository/postgres/builders/search.go"
	path := "tmpl/domain/repository/postgres/builders/search.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchGo = content

	return nil
}

func loadRepositoryPostgresBuildersSearchFiltersGo(v *RepositoryPostgresBuilders, data any) error {
	name := "domains/repository/postgres/builders/search_filters.go"
	path := "tmpl/domain/repository/postgres/builders/search_filters.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchFiltersGo = content

	return nil
}

func loadRepositoryPostgresBuildersSearchOrdersGo(v *RepositoryPostgresBuilders, data any) error {
	name := "domains/repository/postgres/builders/search_orders.go"
	path := "tmpl/domain/repository/postgres/builders/search_orders.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchOrdersGo = content

	return nil
}

func loadRepositoryPostgresBuildersSearchPaginationGo(v *RepositoryPostgresBuilders, data any) error {
	name := "domains/repository/postgres/builders/search_pagination.go"
	path := "tmpl/domain/repository/postgres/builders/search_pagination.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.SearchPaginationGo = content

	return nil
}

type RepositoryPostgresDaos struct {
	DaosGo   []byte
	DeleteGo []byte
	UpdateGo []byte
}

func loadDomainsRepositoryPostgresDaos(v *RepositoryPostgresDaos, data any) error {
	loaders := []func(*RepositoryPostgresDaos, any) error{
		loadRepositoryPostgresDaosDaosGo,
		loadRepositoryPostgresDaosDeleteGo,
		loadRepositoryPostgresDaosUpdateGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadRepositoryPostgresDaosDaosGo(v *RepositoryPostgresDaos, data any) error {
	name := "domains/repository/postgres/daos/daos.go"
	path := "tmpl/domain/repository/postgres/daos/daos.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.DaosGo = content

	return nil
}

func loadRepositoryPostgresDaosDeleteGo(v *RepositoryPostgresDaos, data any) error {
	name := "domains/repository/postgres/daos/delete.go"
	path := "tmpl/domain/repository/postgres/daos/delete.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.DeleteGo = content

	return nil
}

func loadRepositoryPostgresDaosUpdateGo(v *RepositoryPostgresDaos, data any) error {
	name := "domains/repository/postgres/daos/update.go"
	path := "tmpl/domain/repository/postgres/daos/update.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.UpdateGo = content

	return nil
}

type RepositoryDynamo struct {
	DotMockeryPkgYml []byte
	InterfacesGo     []byte
	RepositoryGo     []byte
	ProviderGo       []byte
}

func loadDomainsRepositoryDynamo(v *RepositoryDynamo, data any) error {
	loaders := []func(*RepositoryDynamo, any) error{
		loadRepositoryDynamoDotMockeryPkgYml,
		loadRepositoryDynamoInterfacesGo,
		loadRepositoryDynamoRepositoryGo,
		loadRepositoryDynamoProviderGo,
	}

	for _, loader := range loaders {
		if err := loader(v, data); err != nil {
			return err
		}
	}

	return nil
}

func loadRepositoryDynamoDotMockeryPkgYml(v *RepositoryDynamo, data any) error {
	name := "domains/repository/dynamo/.mockery.pkg.yml"
	path := "tmpl/domain/repository/dynamo/.mockery.pkg.yml.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.DotMockeryPkgYml = content

	return nil
}

func loadRepositoryDynamoInterfacesGo(v *RepositoryDynamo, data any) error {
	name := "domains/repository/dynamo/interfaces.go"
	path := "tmpl/domain/repository/dynamo/interfaces.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.InterfacesGo = content

	return nil
}

func loadRepositoryDynamoRepositoryGo(v *RepositoryDynamo, data any) error {
	name := "domains/repository/dynamo/repository.go"
	path := "tmpl/domain/repository/dynamo/repository.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.RepositoryGo = content

	return nil
}

func loadRepositoryDynamoProviderGo(v *RepositoryDynamo, data any) error {
	name := "domains/repository/dynamo/provider.go"
	path := "tmpl/domain/repository/dynamo/provider.go.tmpl"

	content, err := loadTemplate(name, path, data, domain)
	if err != nil {
		return err
	}

	v.ProviderGo = content

	return nil
}
