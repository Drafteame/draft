package newdomain

import (
	"fmt"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/pkg/dirs"
	"github.com/Drafteame/draft/internal/pkg/files"
)

func (nd *NewDomain) exec() error {
	// 1. Create directory structure
	if err := nd.createAllDirs(); err != nil {
		return err
	}

	// 2. Create domain layer (postgres only)
	if nd.input.DBType == data.DBTypePostgres {
		if err := nd.createDomain(); err != nil {
			return err
		}
	}

	// 3. Create service layer
	if err := nd.createService(); err != nil {
		return err
	}

	// 4. Create repository layer
	if err := nd.createRepository(); err != nil {
		return err
	}

	// 5. Create nanoid provider (postgres only)
	if nd.input.DBType == data.DBTypePostgres {
		if err := nd.createNanoidProvider(); err != nil {
			return err
		}
	}

	return nil
}

func (nd *NewDomain) createAllDirs() error {
	dirList := []string{
		nd.input.DomainPath + "/service",
		nd.input.DomainPath + "/repository",
	}

	// Only create postgres-specific directories
	if nd.input.DBType == data.DBTypePostgres {
		dirList = append(dirList,
			nd.input.DomainPath+"/repository/builders",
			nd.input.DomainPath+"/repository/daos",
			nd.input.DomainPath+"/domain/options",
			"pkg/providers/generators/nanoid/tableid",
		)
	}

	for _, dir := range dirList {
		if err := dirs.Create(dir); err != nil {
			return err
		}
	}

	return nil
}

func (nd *NewDomain) createService() error {
	switch nd.input.DBType {
	case data.DBTypePostgres:
		return nd.createPostgresService()
	case data.DBTypeDynamo:
		return nd.createDynamoService()
	default:
		return fmt.Errorf("unsupported service db type: %s", nd.input.DBType)
	}
}

// createFiles creates all files from the provided file list
func (nd *NewDomain) createFiles(fileList []dtos.FileEntry) error {
	for _, file := range fileList {
		if err := files.Create(file.Path, file.Data); err != nil {
			return err
		}
	}

	return nil
}

func (nd *NewDomain) createPostgresService() error {
	fileList := []dtos.FileEntry{
		{Path: nd.input.DomainPath + "/service/create.go", Data: nd.tmpl.Service.Postgres.CreateGo},
		{Path: nd.input.DomainPath + "/service/create_test.go", Data: nd.tmpl.Service.Postgres.CreateTestGo},
		{Path: nd.input.DomainPath + "/service/delete.go", Data: nd.tmpl.Service.Postgres.DeleteGo},
		{Path: nd.input.DomainPath + "/service/delete_test.go", Data: nd.tmpl.Service.Postgres.DeleteTestGo},
		{Path: nd.input.DomainPath + "/service/get.go", Data: nd.tmpl.Service.Postgres.GetGo},
		{Path: nd.input.DomainPath + "/service/get_test.go", Data: nd.tmpl.Service.Postgres.GetTestGo},
		{Path: nd.input.DomainPath + "/service/interfaces.go", Data: nd.tmpl.Service.Postgres.InterfacesGo},
		{Path: nd.input.DomainPath + "/service/search.go", Data: nd.tmpl.Service.Postgres.SearchGo},
		{Path: nd.input.DomainPath + "/service/search_test.go", Data: nd.tmpl.Service.Postgres.SearchTestGo},
		{Path: nd.input.DomainPath + "/service/update.go", Data: nd.tmpl.Service.Postgres.UpdateGo},
		{Path: nd.input.DomainPath + "/service/update_test.go", Data: nd.tmpl.Service.Postgres.UpdateTestGo},
		{Path: nd.input.DomainPath + "/service/service.go", Data: nd.tmpl.Service.Postgres.ServiceGo},
		{Path: nd.input.DomainPath + "/service/service_test.go", Data: nd.tmpl.Service.Postgres.ServiceTestGo},
		{Path: nd.input.DomainPath + "/service/search_one.go", Data: nd.tmpl.Service.Postgres.SearchOneGo},
		{Path: nd.input.DomainPath + "/service/search_one_test.go", Data: nd.tmpl.Service.Postgres.SearchOneTestGo},
		{Path: nd.input.DomainPath + "/service/provide.go", Data: nd.tmpl.Service.Postgres.ProvideGo},
	}

	return nd.createFiles(fileList)
}

func (nd *NewDomain) createDynamoService() error {
	fileList := []dtos.FileEntry{
		{Path: nd.input.DomainPath + "/service/interfaces.go", Data: nd.tmpl.Service.Dynamo.InterfacesGo},
		{Path: nd.input.DomainPath + "/service/service.go", Data: nd.tmpl.Service.Dynamo.ServiceGo},
		{Path: nd.input.DomainPath + "/service/provider.go", Data: nd.tmpl.Service.Dynamo.ProviderGo},
	}

	return nd.createFiles(fileList)
}

func (nd *NewDomain) createRepository() error {
	switch nd.input.DBType {
	case data.DBTypePostgres:
		return nd.createPostgresRepository()
	case data.DBTypeDynamo:
		return nd.createDynamoRepository()
	default:
		return fmt.Errorf("unsupported repository db type: %s", nd.input.DBType)
	}
}

func (nd *NewDomain) createPostgresRepository() error {
	fileList := []dtos.FileEntry{
		{Path: nd.input.DomainPath + "/repository/create.go", Data: nd.tmpl.Repository.Postgres.CreateGo},
		{Path: nd.input.DomainPath + "/repository/create_test.go", Data: nd.tmpl.Repository.Postgres.CreateTestGo},
		{Path: nd.input.DomainPath + "/repository/delete.go", Data: nd.tmpl.Repository.Postgres.DeleteGo},
		{Path: nd.input.DomainPath + "/repository/delete_test.go", Data: nd.tmpl.Repository.Postgres.DeleteTestGo},
		{Path: nd.input.DomainPath + "/repository/get.go", Data: nd.tmpl.Repository.Postgres.GetGo},
		{Path: nd.input.DomainPath + "/repository/get_test.go", Data: nd.tmpl.Repository.Postgres.GetTestGo},
		{Path: nd.input.DomainPath + "/repository/interfaces.go", Data: nd.tmpl.Repository.Postgres.InterfacesGo},
		{Path: nd.input.DomainPath + "/repository/repository.go", Data: nd.tmpl.Repository.Postgres.RepositoryGo},
		{Path: nd.input.DomainPath + "/repository/repository_test.go", Data: nd.tmpl.Repository.Postgres.RepositoryTestGo},
		{Path: nd.input.DomainPath + "/repository/search.go", Data: nd.tmpl.Repository.Postgres.SearchGo},
		{Path: nd.input.DomainPath + "/repository/search_test.go", Data: nd.tmpl.Repository.Postgres.SearchTestGo},
		{Path: nd.input.DomainPath + "/repository/update.go", Data: nd.tmpl.Repository.Postgres.UpdateGo},
		{Path: nd.input.DomainPath + "/repository/update_test.go", Data: nd.tmpl.Repository.Postgres.UpdateTestGo},
		{Path: nd.input.DomainPath + "/repository/search_one.go", Data: nd.tmpl.Repository.Postgres.SearchOneGo},
		{Path: nd.input.DomainPath + "/repository/search_one_test.go", Data: nd.tmpl.Repository.Postgres.SearchOneTestGo},
		{Path: nd.input.DomainPath + "/repository/builders/search.go", Data: nd.tmpl.Repository.Postgres.Builders.SearchGo},
		{Path: nd.input.DomainPath + "/repository/builders/search_filters.go", Data: nd.tmpl.Repository.Postgres.Builders.SearchFiltersGo},
		{Path: nd.input.DomainPath + "/repository/builders/search_orders.go", Data: nd.tmpl.Repository.Postgres.Builders.SearchOrdersGo},
		{Path: nd.input.DomainPath + "/repository/builders/search_pagination.go", Data: nd.tmpl.Repository.Postgres.Builders.SearchPaginationGo},
		{Path: nd.input.DomainPath + "/repository/daos/daos.go", Data: nd.tmpl.Repository.Postgres.Daos.DaosGo},
		{Path: nd.input.DomainPath + "/repository/daos/delete.go", Data: nd.tmpl.Repository.Postgres.Daos.DeleteGo},
		{Path: nd.input.DomainPath + "/repository/daos/update.go", Data: nd.tmpl.Repository.Postgres.Daos.UpdateGo},
		{Path: nd.input.DomainPath + "/repository/provide.go", Data: nd.tmpl.Repository.Postgres.ProvideGo},
	}

	return nd.createFiles(fileList)
}

func (nd *NewDomain) createDynamoRepository() error {
	fileList := []dtos.FileEntry{
		{Path: nd.input.DomainPath + "/repository/interfaces.go", Data: nd.tmpl.Repository.Dynamo.InterfacesGo},
		{Path: nd.input.DomainPath + "/repository/repository.go", Data: nd.tmpl.Repository.Dynamo.RepositoryGo},
		{Path: nd.input.DomainPath + "/repository/provider.go", Data: nd.tmpl.Repository.Dynamo.ProviderGo},
	}

	return nd.createFiles(fileList)
}

func (nd *NewDomain) createDomain() error {
	fileList := []dtos.FileEntry{
		{Path: nd.input.DomainPath + "/domain/options/search.go", Data: nd.tmpl.Domain.Options.SearchGo},
		{Path: nd.input.DomainPath + "/domain/options/search_filters.go", Data: nd.tmpl.Domain.Options.SearchFiltersGo},
		{Path: nd.input.DomainPath + "/domain/options/search_orders.go", Data: nd.tmpl.Domain.Options.SearchOrdersGo},
		{Path: nd.input.DomainPath + "/domain/options/search_pagination.go", Data: nd.tmpl.Domain.Options.SearchPaginationGo},
		{Path: nd.input.DomainPath + "/domain/options/update_fields.go", Data: nd.tmpl.Domain.Options.UpdateFieldsGo},
		{Path: nd.input.DomainPath + "/domain/domain.go", Data: nd.tmpl.Domain.DomainGo},
		{Path: nd.input.DomainPath + "/domain/errors.go", Data: nd.tmpl.Domain.ErrorsGo},
	}

	return nd.createFiles(fileList)
}

func (nd *NewDomain) createNanoidProvider() error {
	fileList := []dtos.FileEntry{
		{Path: "pkg/providers/generators/nanoid/tableid/" + nd.input.DomainNameLower + ".go", Data: nd.tmpl.Providers.GeneratorsNanoidTableid.ProvideGo},
	}

	return nd.createFiles(fileList)
}
