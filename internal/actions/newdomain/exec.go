package newdomain

import (
	"fmt"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
	"github.com/Drafteame/draft/internal/pkg/files"
)

func (nd *NewDomain) exec() error {
	creators := []func() error{
		nd.createAllDirs,
		nd.createService,
		nd.createRepository,
		nd.createProviders,
		nd.createDomain,
	}

	for _, creator := range creators {
		if err := creator(); err != nil {
			return err
		}
	}

	return nil
}

func (nd *NewDomain) createAllDirs() error {
	dirList := []string{
		"domains/" + nd.input.DomainPath + "/service",
		"domains/" + nd.input.DomainPath + "/repository",
		"domains/" + nd.input.DomainPath + "/providers",
		"domains/" + nd.input.DomainPath + "/domain/options",
	}

	for _, dir := range dirList {
		if err := dirs.Create(dir); err != nil {
			return err
		}
	}

	return nil
}

func (nd *NewDomain) createService() error {
	fileList := []dtos.FileEntry{
		{Path: "domains/" + nd.input.DomainPath + "/service/create.go", Data: nd.tmpl.Service.CreateGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/create_test.go", Data: nd.tmpl.Service.CreateTestGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/delete.go", Data: nd.tmpl.Service.DeleteGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/delete_test.go", Data: nd.tmpl.Service.DeleteTestGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/get.go", Data: nd.tmpl.Service.GetGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/get_test.go", Data: nd.tmpl.Service.GetTestGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/interfaces.go", Data: nd.tmpl.Service.InterfacesGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/search.go", Data: nd.tmpl.Service.SearchGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/search_test.go", Data: nd.tmpl.Service.SearchTestGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/update.go", Data: nd.tmpl.Service.UpdateGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/update_test.go", Data: nd.tmpl.Service.UpdateTestGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/service.go", Data: nd.tmpl.Service.ServiceGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/service_test.go", Data: nd.tmpl.Service.ServiceTestGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/search_one.go", Data: nd.tmpl.Service.SearchOneGo},
		{Path: "domains/" + nd.input.DomainPath + "/service/search_one_test.go", Data: nd.tmpl.Service.SearchOneTestGo},
	}

	for _, file := range fileList {
		if err := files.Create(file.Path, file.Data); err != nil {
			return err
		}
	}

	return nil
}

func (nd *NewDomain) createRepository() error {
	switch nd.input.DBType {
	case "postgres":
		return nd.createPostgresRepository()
	default:
		return fmt.Errorf("unsupported repository db type: %s", nd.input.DBType)
	}
}

func (nd *NewDomain) createPostgresRepository() error {
	fileList := []dtos.FileEntry{
		{Path: "domains/" + nd.input.DomainPath + "/repository/create.go", Data: nd.tmpl.Repository.Postgres.CreateGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/create_test.go", Data: nd.tmpl.Repository.Postgres.CreateTestGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/delete.go", Data: nd.tmpl.Repository.Postgres.DeleteGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/delete_test.go", Data: nd.tmpl.Repository.Postgres.DeleteTestGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/get.go", Data: nd.tmpl.Repository.Postgres.GetGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/get_test.go", Data: nd.tmpl.Repository.Postgres.GetTestGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/interfaces.go", Data: nd.tmpl.Repository.Postgres.InterfacesGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/repository.go", Data: nd.tmpl.Repository.Postgres.RepositoryGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/repository_test.go", Data: nd.tmpl.Repository.Postgres.RepositoryTestGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/search.go", Data: nd.tmpl.Repository.Postgres.SearchGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/search_test.go", Data: nd.tmpl.Repository.Postgres.SearchTestGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/update.go", Data: nd.tmpl.Repository.Postgres.UpdateGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/update_test.go", Data: nd.tmpl.Repository.Postgres.UpdateTestGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/search_one.go", Data: nd.tmpl.Repository.Postgres.SearchOneGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/search_one_test.go", Data: nd.tmpl.Repository.Postgres.SearchOneTestGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/builders/search.go", Data: nd.tmpl.Repository.Postgres.Builders.SearchGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/builders/search_filters.go", Data: nd.tmpl.Repository.Postgres.Builders.SearchFiltersGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/builders/search_orders.go", Data: nd.tmpl.Repository.Postgres.Builders.SearchOrdersGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/builders/search_pagination.go", Data: nd.tmpl.Repository.Postgres.Builders.SearchPaginationGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/daos/daos.go", Data: nd.tmpl.Repository.Postgres.Daos.DaosGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/daos/delete.go", Data: nd.tmpl.Repository.Postgres.Daos.DeleteGo},
		{Path: "domains/" + nd.input.DomainPath + "/repository/daos/update.go", Data: nd.tmpl.Repository.Postgres.Daos.UpdateGo},
	}

	for _, file := range fileList {
		if err := files.Create(file.Path, file.Data); err != nil {
			return err
		}
	}

	return nil
}

func (nd *NewDomain) createProviders() error {
	fileList := []dtos.FileEntry{
		{Path: "domains/" + nd.input.DomainPath + "/providers/generator.go", Data: nd.tmpl.Providers.GeneratorGo},
		{Path: "domains/" + nd.input.DomainPath + "/providers/providers.go", Data: nd.tmpl.Providers.ProvidersGo},
		{Path: "domains/" + nd.input.DomainPath + "/providers/service.go", Data: nd.tmpl.Providers.ServiceGo},
	}

	for _, file := range fileList {
		if err := files.Create(file.Path, file.Data); err != nil {
			return err
		}
	}

	switch nd.input.DBType {
	case "postgres":
		return nd.createPostgresProviders()
	default:
		return fmt.Errorf("unsupported providers db type: %s", nd.input.DBType)
	}
}

func (nd *NewDomain) createPostgresProviders() error {
	fileList := []dtos.FileEntry{
		{Path: "domains/" + nd.input.DomainPath + "/providers/repository.go", Data: nd.tmpl.Providers.Postgres.RepositoryGo},
	}

	for _, file := range fileList {
		if err := files.Create(file.Path, file.Data); err != nil {
			return err
		}
	}

	return nil
}

func (nd *NewDomain) createDomain() error {
	fileList := []dtos.FileEntry{
		{Path: "domains/" + nd.input.DomainPath + "/domain/options/search.go", Data: nd.tmpl.Domain.Options.SearchGo},
		{Path: "domains/" + nd.input.DomainPath + "/domain/options/search_filters.go", Data: nd.tmpl.Domain.Options.SearchFiltersGo},
		{Path: "domains/" + nd.input.DomainPath + "/domain/options/search_orders.go", Data: nd.tmpl.Domain.Options.SearchOrdersGo},
		{Path: "domains/" + nd.input.DomainPath + "/domain/options/search_pagination.go", Data: nd.tmpl.Domain.Options.SearchPaginationGo},
		{Path: "domains/" + nd.input.DomainPath + "/domain/options/update_fields.go", Data: nd.tmpl.Domain.Options.UpdateFieldsGo},
		{Path: "domains/" + nd.input.DomainPath + "/domain/domain.go", Data: nd.tmpl.Domain.DomainGo},
		{Path: "domains/" + nd.input.DomainPath + "/domain/errors.go", Data: nd.tmpl.Domain.ErrorsGo},
	}

	for _, file := range fileList {
		if err := files.Create(file.Path, file.Data); err != nil {
			return err
		}
	}

	return nil
}
