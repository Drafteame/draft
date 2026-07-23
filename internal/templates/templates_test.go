package templates_test

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/templates"
)

// baseLambdaInput returns a populated LambdaInput for template rendering tests.
func baseLambdaInput(lambdaType string) dtos.LambdaInput {
	return dtos.LambdaInput{
		PackageName:         "github.com/Drafteame/api",
		ServicePath:         "services/testsvc",
		ServiceName:         "testsvc",
		LambdaName:          "testlambda",
		LambdaType:          lambdaType,
		CustomTypePath:      "custom",
		FrameVersion:        "v2",
		QueueARN:            "arn:aws:sqs:us-east-2:123456789012:test-queue",
		HTTPPath:            "/test",
		HTTPPathAPIGateway:  "/test",
		HTTPPathEcho:        "/test",
		HTTPMethod:          "GET",
		CronExpression:      "rate(1 hour)",
		NextImportTag:       data.NextImportTag,
		NextLambdaImportTag: data.NextLambdaImportTag,
		IsLegacy:            false,
		UseDig:              false,
		ReservedConcurrency: "medium.eventDriven",
		UseIdempotency:      false,
	}
}

func baseServiceInput() dtos.ServiceInput {
	return dtos.ServiceInput{
		PackageName:           "github.com/Drafteame/api",
		ServiceName:           "testsvc",
		NormalizedServiceName: "testsvc",
		ServicePath:           "services/testsvc",
		ServicePackage:        "testsvc",
		LambdaName:            "helloworld",
		LambdaType:            "plain",
		CustomDomain:          false,
		DomainPath:            "",
		FrameVersion:          "v2",
		HasSentry:             false,
		SentryDSN:             "",
		NextImportTag:         data.NextImportTag,
		NextLambdaImportTag:   data.NextLambdaImportTag,
		IsLegacy:              false,
		UseDig:                false,
		ReservedConcurrency:   "medium.http",
		RoleName:              "Testsvc",
	}
}

func baseDomainInput(dbType string) dtos.DomainInput {
	return dtos.DomainInput{
		PackageName:        "github.com/Drafteame/api",
		DomainPath:         "domains/testdomain",
		DomainName:         "testdomain",
		DomainNamePascal:   "Testdomain",
		DomainNameLower:    "testdomain",
		DBPrefix:           "tst",
		TableName:          "public.testdomains",
		DBProviderFuncName: "TestDB",
		DBType:             dbType,
		DBName:             "testdb",
	}
}

// assertGoSyntax validates that content is syntactically valid Go source.
func assertGoSyntax(t *testing.T, name string, content []byte) {
	t.Helper()

	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, name, content, parser.AllErrors)
	assert.NoError(t, err, "invalid Go syntax in %s:\n%s", name, string(content))
}

// assertYAMLSyntax validates that content is valid YAML.
func assertYAMLSyntax(t *testing.T, name string, content []byte) {
	t.Helper()

	var out any
	err := yaml.Unmarshal(content, &out)
	assert.NoError(t, err, "invalid YAML in %s", name)
}

// =============================================================================
// Lambda templates
// =============================================================================

func TestLambdaTemplates_Plain(t *testing.T) {
	tmpl, err := templates.NewLambdaTemplates(baseLambdaInput("plain"))
	require.NoError(t, err)

	p := tmpl.Plain
	assertGoSyntax(t, "plain/main.go", p.MainGo)
	assertYAMLSyntax(t, "plain/lambda-config.yml", p.LambdaConfigYAML)
	assertGoSyntax(t, "plain/handler/bootstrap.go", p.Handler.BootstrapGo)
	assertGoSyntax(t, "plain/handler/worker/worker.go", p.Handler.WorkerGo)
	assertGoSyntax(t, "plain/handler/worker/resources.go", p.Handler.ResourcesGo)
	assertGoSyntax(t, "plain/handler/dtos/dto.go", p.Handler.DtosGo)
}

func TestLambdaTemplates_HTTP(t *testing.T) {
	tmpl, err := templates.NewLambdaTemplates(baseLambdaInput("http"))
	require.NoError(t, err)

	h := tmpl.HTTP
	assertGoSyntax(t, "http/main.go", h.MainGo)
	assertYAMLSyntax(t, "http/lambda-config.yml", h.LambdaConfigYAML)
	assertGoSyntax(t, "http/handler/bootstrap.go", h.Handler.BootstrapGo)
	assertGoSyntax(t, "http/handler/worker/worker.go", h.Handler.WorkerGo)
	assertGoSyntax(t, "http/handler/worker/resources.go", h.Handler.ResourcesGo)
}

func TestLambdaTemplates_SQS(t *testing.T) {
	tmpl, err := templates.NewLambdaTemplates(baseLambdaInput("sqs"))
	require.NoError(t, err)

	s := tmpl.Sqs
	assertGoSyntax(t, "sqs/main.go", s.MainGo)
	assertYAMLSyntax(t, "sqs/lambda-config.yml", s.LambdaConfigYAML)
	assertGoSyntax(t, "sqs/handler/bootstrap.go", s.Handler.BootstrapGo)
	assertGoSyntax(t, "sqs/handler/worker/worker.go", s.Handler.WorkerGo)
	assertGoSyntax(t, "sqs/handler/worker/resources.go", s.Handler.ResourcesGo)
	assertGoSyntax(t, "sqs/handler/dtos/dto.go", s.Handler.DtosGo)
	assertGoSyntax(t, "sqs/handler/worker/idempotency.go", s.Handler.IdempotencyGo)
	assertGoSyntax(t, "sqs/handler/worker/interfaces.go", s.Handler.InterfacesGo)
}

func TestLambdaTemplates_Cron(t *testing.T) {
	tmpl, err := templates.NewLambdaTemplates(baseLambdaInput("cron"))
	require.NoError(t, err)

	c := tmpl.Cron
	assertGoSyntax(t, "cron/main.go", c.MainGo)
	assertYAMLSyntax(t, "cron/lambda-config.yml", c.LambdaConfigYAML)
	assertGoSyntax(t, "cron/handler/bootstrap.go", c.Handler.BootstrapGo)
	assertGoSyntax(t, "cron/handler/worker/worker.go", c.Handler.WorkerGo)
}

func TestLambdaTemplates_SnsSqs(t *testing.T) {
	tmpl, err := templates.NewLambdaTemplates(baseLambdaInput("snssqs"))
	require.NoError(t, err)

	ss := tmpl.SnsSqs
	assertGoSyntax(t, "snssqs/main.go", ss.MainGo)
	assertYAMLSyntax(t, "snssqs/lambda-config.yml", ss.LambdaConfigYAML)
	assertGoSyntax(t, "snssqs/handler/bootstrap.go", ss.Handler.BootstrapGo)
	assertGoSyntax(t, "snssqs/handler/worker/worker.go", ss.Handler.WorkerGo)
	assertGoSyntax(t, "snssqs/handler/worker/resources.go", ss.Handler.ResourcesGo)
	assertGoSyntax(t, "snssqs/handler/dtos/dto.go", ss.Handler.DtosGo)
	assertGoSyntax(t, "snssqs/handler/worker/idempotency.go", ss.Handler.IdempotencyGo)
	assertGoSyntax(t, "snssqs/handler/worker/interfaces.go", ss.Handler.InterfacesGo)
}

func TestLambdaTemplates_Custom(t *testing.T) {
	tmpl, err := templates.NewLambdaTemplates(baseLambdaInput("custom"))
	require.NoError(t, err)

	c := tmpl.Custom
	assertGoSyntax(t, "custom/main.go", c.MainGo)
	assertYAMLSyntax(t, "custom/lambda-config.yml", c.LambdaConfigYAML)
	assertGoSyntax(t, "custom/handler/bootstrap.go", c.Handler.BootstrapGo)
	assertGoSyntax(t, "custom/handler/worker/worker.go", c.Handler.WorkerGo)
	assertGoSyntax(t, "custom/handler/worker/resources.go", c.Handler.ResourcesGo)
	assertGoSyntax(t, "custom/handler/worker/idempotency.go", c.Handler.IdempotencyGo)
	assertGoSyntax(t, "custom/handler/worker/interfaces.go", c.Handler.InterfacesGo)
	assertGoSyntax(t, "custom/handler/worker/worker_setup_test.go", c.Handler.WorkerSetupTestGo)
	assertGoSyntax(t, "custom/handler/worker/worker_test.go", c.Handler.WorkerTestGo)
}

// =============================================================================
// Service templates
// =============================================================================

func TestServiceTemplates(t *testing.T) {
	tmpl, err := templates.NewServiceTemplates(baseServiceInput())
	require.NoError(t, err)

	assertYAMLSyntax(t, "serverless.yml", tmpl.ServerlessYAML)
	assertGoSyntax(t, "deps.go", tmpl.DepsGo)
	assertYAMLSyntax(t, "config/sls/environment.yml", tmpl.Config.Sls.EnvironmentYAML)
	assertYAMLSyntax(t, "config/sls/resources.yml", tmpl.Config.Sls.ResourcesYAML)

	// Initial plain lambda generated with service
	p := tmpl.Lambda.Plain
	assertGoSyntax(t, "cmd/plain/helloworld/main.go", p.MainGo)
	assertYAMLSyntax(t, "cmd/plain/helloworld/lambda-config.yml", p.LambdaConfigYAML)
	assertGoSyntax(t, "cmd/plain/helloworld/handler/bootstrap.go", p.Handler.BootstrapGo)
	assertGoSyntax(t, "cmd/plain/helloworld/handler/worker/worker.go", p.Handler.WorkerGo)
	assertGoSyntax(t, "cmd/plain/helloworld/handler/worker/resources.go", p.Handler.ResourcesGo)
	assertGoSyntax(t, "cmd/plain/helloworld/handler/dtos/dto.go", p.Handler.DtosGo)
}

func TestServiceTemplates_WithCustomDomain(t *testing.T) {
	input := baseServiceInput()
	input.CustomDomain = true
	input.DomainPath = "api/v1"

	tmpl, err := templates.NewServiceTemplates(input)
	require.NoError(t, err)

	assertYAMLSyntax(t, "serverless.yml", tmpl.ServerlessYAML)
}

// =============================================================================
// Domain templates — postgres
// =============================================================================

func TestDomainTemplates_Postgres(t *testing.T) {
	input := baseDomainInput(data.DBTypePostgres)

	d, err := templates.NewDomains(input)
	require.NoError(t, err)

	t.Run("service", func(t *testing.T) {
		svc := d.Service.Postgres
		assertGoSyntax(t, "service/create.go", svc.CreateGo)
		assertGoSyntax(t, "service/create_test.go", svc.CreateTestGo)
		assertGoSyntax(t, "service/get.go", svc.GetGo)
		assertGoSyntax(t, "service/get_test.go", svc.GetTestGo)
		assertGoSyntax(t, "service/update.go", svc.UpdateGo)
		assertGoSyntax(t, "service/update_test.go", svc.UpdateTestGo)
		assertGoSyntax(t, "service/delete.go", svc.DeleteGo)
		assertGoSyntax(t, "service/delete_test.go", svc.DeleteTestGo)
		assertGoSyntax(t, "service/search.go", svc.SearchGo)
		assertGoSyntax(t, "service/search_test.go", svc.SearchTestGo)
		assertGoSyntax(t, "service/search_one.go", svc.SearchOneGo)
		assertGoSyntax(t, "service/search_one_test.go", svc.SearchOneTestGo)
		assertGoSyntax(t, "service/service.go", svc.ServiceGo)
		assertGoSyntax(t, "service/service_test.go", svc.ServiceTestGo)
		assertGoSyntax(t, "service/interfaces.go", svc.InterfacesGo)
		assertGoSyntax(t, "service/provide.go", svc.ProvideGo)
	})

	t.Run("repository", func(t *testing.T) {
		repo := d.Repository.Postgres
		assertGoSyntax(t, "repository/create.go", repo.CreateGo)
		assertGoSyntax(t, "repository/create_test.go", repo.CreateTestGo)
		assertGoSyntax(t, "repository/get.go", repo.GetGo)
		assertGoSyntax(t, "repository/get_test.go", repo.GetTestGo)
		assertGoSyntax(t, "repository/update.go", repo.UpdateGo)
		assertGoSyntax(t, "repository/update_test.go", repo.UpdateTestGo)
		assertGoSyntax(t, "repository/delete.go", repo.DeleteGo)
		assertGoSyntax(t, "repository/delete_test.go", repo.DeleteTestGo)
		assertGoSyntax(t, "repository/search.go", repo.SearchGo)
		assertGoSyntax(t, "repository/search_test.go", repo.SearchTestGo)
		assertGoSyntax(t, "repository/search_one.go", repo.SearchOneGo)
		assertGoSyntax(t, "repository/search_one_test.go", repo.SearchOneTestGo)
		assertGoSyntax(t, "repository/repository.go", repo.RepositoryGo)
		assertGoSyntax(t, "repository/repository_test.go", repo.RepositoryTestGo)
		assertGoSyntax(t, "repository/interfaces.go", repo.InterfacesGo)
		assertGoSyntax(t, "repository/provide.go", repo.ProvideGo)
		assertGoSyntax(t, "repository/builders/search.go", repo.Builders.SearchGo)
		assertGoSyntax(t, "repository/builders/search_filters.go", repo.Builders.SearchFiltersGo)
		assertGoSyntax(t, "repository/builders/search_orders.go", repo.Builders.SearchOrdersGo)
		assertGoSyntax(t, "repository/builders/search_pagination.go", repo.Builders.SearchPaginationGo)
		assertGoSyntax(t, "repository/daos/daos.go", repo.Daos.DaosGo)
		assertGoSyntax(t, "repository/daos/delete.go", repo.Daos.DeleteGo)
		assertGoSyntax(t, "repository/daos/update.go", repo.Daos.UpdateGo)
	})

	t.Run("domain", func(t *testing.T) {
		dom := d.Domain
		assertGoSyntax(t, "domain/domain.go", dom.DomainGo)
		assertGoSyntax(t, "domain/errors.go", dom.ErrorsGo)
		assertGoSyntax(t, "domain/options/search.go", dom.Options.SearchGo)
		assertGoSyntax(t, "domain/options/search_filters.go", dom.Options.SearchFiltersGo)
		assertGoSyntax(t, "domain/options/search_orders.go", dom.Options.SearchOrdersGo)
		assertGoSyntax(t, "domain/options/search_pagination.go", dom.Options.SearchPaginationGo)
		assertGoSyntax(t, "domain/options/update_fields.go", dom.Options.UpdateFieldsGo)
	})

	t.Run("providers", func(t *testing.T) {
		prov := d.Providers
		assertGoSyntax(t, "providers/generators/nanoid/tableid/provide.go", prov.GeneratorsNanoidTableid.ProvideGo)
	})
}

// =============================================================================
// Domain templates — dynamo
// =============================================================================

func TestDomainTemplates_Dynamo(t *testing.T) {
	input := baseDomainInput(data.DBTypeDynamo)

	d, err := templates.NewDomains(input)
	require.NoError(t, err)

	// DynamoDB domains only generate service and repository layers.
	// The domain entity layer (options, errors, etc.) is postgres-only.

	t.Run("service", func(t *testing.T) {
		svc := d.Service.Dynamo
		assertGoSyntax(t, "service/service.go", svc.ServiceGo)
		assertGoSyntax(t, "service/interfaces.go", svc.InterfacesGo)
		assertGoSyntax(t, "service/provide.go", svc.ProviderGo)
	})

	t.Run("repository", func(t *testing.T) {
		repo := d.Repository.Dynamo
		assertGoSyntax(t, "repository/repository.go", repo.RepositoryGo)
		assertGoSyntax(t, "repository/interfaces.go", repo.InterfacesGo)
		assertGoSyntax(t, "repository/provide.go", repo.ProviderGo)
	})
}
