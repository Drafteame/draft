package newdomain

import (
	"github.com/spf13/cobra"

	"github.com/Drafteame/draft/cmd/commands/internal/common"
	"github.com/Drafteame/draft/internal/actions/newdomain"
	"github.com/Drafteame/draft/internal/data"
	"github.com/Drafteame/draft/internal/dtos"
	"github.com/Drafteame/draft/internal/forms"
	"github.com/Drafteame/draft/internal/pkg/log"
)

var (
	domainPath string
	dbType     string
	tableName  string
	dbPrefix   string
	dbName     string
)

var newDomainCmd = &cobra.Command{
	Use:   "new:domain",
	Short: "Create a new domain with models, services, and repositories",
	Long: `Create a new configurable domain with complete domain-driven design structure.

This command scaffolds a domain layer including:
- Domain models and business logic (Postgres only)
- Service layer with CRUD operations
- Repository layer with database interactions
- Provider functions for dependency injection
- Mock configurations for testing

The command supports both interactive mode (prompts for values) and non-interactive
mode (uses flags) for CI/CD automation.

Database Types:
  postgres - Creates full domain structure with Postgres support
             Includes: search, filters, pagination, DAOs, and builders
  dynamo   - Creates simplified structure optimized for DynamoDB
             Includes: basic CRUD operations

Examples:
  # Interactive mode - prompts for all values
  draft new:domain

  # Non-interactive mode with Postgres
  draft new:domain \
    --domain-path users \
    --db-type postgres \
    --table-name public.users \
    --db-prefix usr \
    --db-name general

  # Non-interactive mode with DynamoDB
  draft new:domain \
    --domain-path products \
    --db-type dynamo \
    --table-name ProductsTable

  # Mixed mode - provide some flags, prompt for others
  draft new:domain --domain-path orders --db-type postgres

  # With custom working directory
  draft new:domain -w /path/to/project --domain-path inventory --db-type postgres

Database Configuration:
  For Postgres domains, the list of available databases is dynamically loaded from
  the .local-migrate-config.yml file in the project root. The command will:

  1. Read migrations.databases from .local-migrate-config.yml
  2. Filter out any databases with group: 'test'
  3. Present the remaining databases as options
  4. Convert database names to PascalCase for provider functions
     Example: user_preferences -> ProvideUserPreferences
             games_core -> ProvideGamesCore

  If the .local-migrate-config.yml file is not found, the command will fail with
  an error. Ensure this file exists in your project root before running the command.`,
	Run: run,
}

func init() {
	newDomainCmd.Flags().StringVarP(&domainPath, "domain-path", "p", "", "Path to the domain folder")
	newDomainCmd.Flags().StringVar(&dbType, "db-type", "", "Database type (postgres or dynamo)")
	newDomainCmd.Flags().StringVar(&tableName, "table-name", "", "Name of the database table")
	newDomainCmd.Flags().StringVar(&dbPrefix, "db-prefix", "", "ID prefix for Postgres (3 characters)")
	newDomainCmd.Flags().StringVar(&dbName, "db-name", "", "Database name for Postgres (loaded from .local-migrate-config.yml). Use the snake_case database name (e.g., 'general', 'user_preferences')")
}

func run(cmd *cobra.Command, _ []string) {
	common.ChDir(cmd)

	data.LoadMeta()

	input := dtos.DomainInput{
		DomainPath: domainPath,
		DBType:     dbType,
		TableName:  tableName,
		DBPrefix:   dbPrefix,
		DBName:     dbName,
	}

	if err := forms.NewDomain(&input); err != nil {
		log.Exitf(1, "Failed to collect domain info: %v", err)
	}

	if errExec := newdomain.New(cmd.Context(), input).Exec(); errExec != nil {
		log.Exitf(1, "Failed to create domain: %v", errExec)
	}

	log.Successf("Domain %s created successfully", input.DomainName)
}

func GetCmd() *cobra.Command {
	return newDomainCmd
}
