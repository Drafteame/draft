# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Draft is a CLI tool for generating service and Lambda folder structures for Draftea Backend Services. It scaffolds serverless applications with Go, creating standardized project structures including Lambda functions, domains, configuration files, and service boilerplate.

## Core Commands

### Development
```bash
# Format code
task format

# Lint code
task lint

# Install git hooks
task hooks
```

### Building
The project uses `go install` for distribution. Build artifacts are managed in `.bin/` directories (ignored by git).

### Testing
No test suite is currently configured. The codebase does not contain test files.

### Running Draft Commands
```bash
# Create a new service (run from monorepo service folder)
draft new:service
draft new:service -w path/to/project
draft new:service -l path/to/legacy/service

# Create a new Lambda (run from monorepo root)
draft new:lambda
draft new:lambda -w path/to/project

# Create a new domain
draft new:domain

# Invoke a Lambda locally
draft invoke path/to/lambda
draft invoke --body '{"key": "value"}' path/to/lambda
draft invoke --body-file path/to/body.json path/to/lambda

# Local development commands
draft local:setup
draft local:migrate:up
draft local:migrate:down
draft local:migrate:force

# Sentry project management
draft sentry:project:create
draft sentry:project:delete

# Mockery mock generation
draft mockery                                           # Run mockery for all .mockery.pkg.yml files
draft mockery path/to/.mockery.pkg.yml                 # Run mockery for specific package configs
draft mockery --jobs-num 5                             # Run with custom concurrent job limit (default: 5)
draft mockery --dry                                    # Dry run - validate configs without executing mockery
draft mockery --git-mod                                # Run only for packages with modified files (git diff)
draft mockery --git-mod --dry --jobs-num 10           # Combine flags for targeted validation
```

### Global Flags
- `-w, --working-dir`: Specify working directory
- `-d, --debug`: Enable debug mode
- `-t, --tty`: TTY mode (default: true)

### Service and Lambda Specific Flags
- `--use-dig`: Use Uber Dig for dependency injection instead of default pattern
- `-l, --legacy-path`: Path to legacy service (use legacy folder structure)

## Architecture

### Command Structure
Commands are registered in `main.go` using Cobra. Each command lives in `cmd/commands/` with its own subdirectory:
- `newservice` - Service scaffolding
- `newlambda` - Lambda scaffolding
- `newdomain` - Domain layer scaffolding
- `mockery` - Mockery mock generation with merged base and package configs
- `local/invoke` - Local Lambda invocation
- `local/migrate/*` - Database migration commands
- `sentry/project/create` - Sentry project creation
- `sentry/project/delete` - Sentry project deletion
- `config` - Configuration management

### Action Layer
The `internal/actions/` directory contains the business logic for each command. All creation actions follow a consistent 3-phase lifecycle:

1. **preCreate()** - Pre-execution hooks (e.g., Sentry project setup for new services)
2. **exec()** - Core logic with type-specific switch statements for lambda types or database types
3. **postCreate()** - Post-execution tasks (code formatting, file updates, running external tools)

Key action patterns:
- `New(input)` factory function creates action instances with input DTOs from forms
- `Exec()` orchestrates the entire lifecycle: preCreate → exec → postCreate
- Actions interact with template system to generate files, then manipulate existing files in postCreate

### Template System
Templates are embedded at compile time using `//go:embed` in `internal/templates/`. Template types:
- **Serverless templates** (`tmpl/sls/`): Service boilerplate, Lambda functions (HTTP, Cron, SQS, SNS+SQS, Plain, Custom)
- **Domain templates** (`tmpl/domain/`): Domain-driven design layers (domain, service, repository with Postgres and DynamoDB support)

Templates use Go's `text/template` and are loaded via `loadTemplate()` helper. Template factories return `[]byte` after rendering with data context from input DTOs.

### Lambda Types
The tool generates different Lambda structures based on trigger type:
- **HTTP**: API Gateway triggered, includes routing and DTOs
- **Cron**: EventBridge scheduled execution
- **SQS**: Queue-based processing with idempotency
- **SNS+SQS**: Combined SNS topic and SQS queue processing
- **Plain**: Generic Lambda for custom event sources
- **Custom**: User-defined custom event sources with configurable type path and optional idempotency

### Domain Generation
When creating domains with `draft new:domain`, the tool supports both interactive and non-interactive modes:

**Interactive Mode** (default):
- Domain path (automatically prefixed with `domains/` if not present)
- Database type selection (Postgres or DynamoDB)
- Database-specific configuration (table names, prefixes, etc.)
- For Postgres: Database selection from `.local-migrate-config.yml`

**Non-Interactive Mode** (CLI flags):
- `--domain-path, -p`: Domain folder path
- `--db-type`: Database type (postgres or dynamo)
- `--table-name`: Database table name
- `--db-prefix`: ID prefix for Postgres (3 characters)
- `--db-name`: Database name for Postgres (loaded from config)

**Database Configuration**:
For Postgres domains, available databases are dynamically loaded from `.local-migrate-config.yml`:
- Reads `migrations.databases` configuration
- Filters out test databases (`group: 'test'`)
- Converts snake_case names to PascalCase for provider functions
  - Example: `user_preferences` → `ProvideUserPreferences`
  - Example: `games_core` → `ProvideGamesCore`

Domain structure varies by database type:
- **Postgres**: Full CRUD with service layer, repository layer, builders, DAOs, providers, and domain models with search/filter capabilities
- **DynamoDB**: Simplified structure with repository and service layers optimized for NoSQL patterns

### Project Detection and Name Normalization
`internal/project/` contains utilities for:
- Extracting Go module names from `go.mod` via `GetPackage()` - stored in `data.Meta.PackageName`
- `NormalizeServiceName()` - Converts to kebab-case for folder names (e.g., "My Service" → "my-service")
- `NormalizeServicePackage()` - Converts to snake_case for Go package names (e.g., "My Service" → "my_service")
- Locating services in monorepos

### Legacy vs Modern Structure
The tool supports two folder structures:
- **Modern**: `cmd/<lambda-type>/<lambda-name>/`
- **Legacy**: `<lambda-type>/<lambda-name>/`

Determined by the `-l` flag or interactive prompts.

### Build System
`internal/pkg/build/` compiles all Lambda functions in parallel:
- Searches for `main.go` files in `cmd/` folders
- Builds with `-tags local -ldflags="-s -w"` for local development
- Outputs binaries to `.bin/` structure mirroring source paths
- Uses goroutine pools limited by CPU count

### Forms and User Input
`internal/forms/` provides interactive CLI prompts using `charmbracelet/huh` with functional options pattern:
- Multi-step forms for service/lambda/domain creation (e.g., baseForm → frameDetails → type-specific forms)
- Input validation and conditional fields based on previous selections
- Support for TTY and non-TTY modes (CI/CD compatibility)
- Generic options pattern: `inputs.Text()`, `inputs.Select()`, `inputs.Confirm()` with `WithDescription`, `WithValue`, `WithValidation`, etc.

### Mockery Command
The `draft mockery` command provides concurrent mock generation with configuration merging:
- **Base Config**: Loads shared settings from `.mockery.base.yml` at project root
- **Package Configs**: Searches for or accepts specific `.mockery.pkg.yml` files in service directories
- **Config Merging**: Deep merges base and package configs (package settings take precedence)
- **Concurrent Execution**: Runs mockery jobs in parallel with configurable concurrency (`--jobs-num` flag, default: 5)
  - Semaphore acquired BEFORE spawning goroutines for proper concurrency control
  - Real-time progress updates with spinner showing completed/total packages
- **Temporary Files**: Creates temporary merged config files (`.mockery.tmp.*.yml`) that are automatically cleaned up
- **Progress Reporting**: Shows real-time progress with spinner and execution statistics
- **Error Handling**: Continues execution on failures, reports all errors at the end, exits with code 1 if any failed
- **Dry Run Mode** (`--dry`): Validates and prepares all configs without executing mockery commands
  - Useful for CI/CD pipelines to verify configuration correctness
  - Completes significantly faster than actual execution
- **Git Diff Mode** (`--git-mod`): Only processes packages with modified files
  - Compares `HEAD` with `origin/main` (or `origin/master` as fallback)
  - Extracts directories from modified files and searches for `.mockery.pkg.yml` in those directories and parents
  - Automatically deduplicates found configs
  - Cannot be combined with explicit config file paths
- **Graceful Cancellation**: Handles Ctrl+C (SIGINT) and SIGTERM signals
  - Stops spawning new goroutines immediately
  - Waits for running tasks to complete gracefully
  - Cleans up all temporary files via defer
  - Displays cancellation summary with completed/cancelled counts
  - Context propagated from `main.go` through Cobra commands

Implementation in `internal/actions/mockery/`:
1. Validates inputs and finds/validates config files (via explicit paths, git diff, or directory walk)
2. Loads base configuration from `.mockery.base.yml`
3. Creates temporary merged configs for each package
4. Executes mockery concurrently with semaphore-based concurrency control
   - Monitors context for cancellation signals
   - Checks context before spawning each goroutine
   - Uses `select` on semaphore acquisition to respect cancellation
5. Cleans up temporary files and displays execution summary (or cancellation summary if interrupted)

The `new:domain` action integrates with the mockery action layer:
- Creates `.mockery.pkg.yml` files for service and repository packages
- Calls the mockery action directly (reuses `internal/actions/mockery`)
- Passes context for cancellation support
- Runs with jobsNum=2 for concurrent service and repository mock generation
- Benefits from progress reporting and error handling of the main mockery action

## Configuration Files

### Pkl Configuration
Services use Pkl (configuration language) for app config:
- `config/app/app.pkl` - Application configuration
- `config/app/modules.pkl` - Pkl module imports

### Serverless Framework
- `serverless.yml` - Service definition and Lambda configuration
- `lambda-config.yml` - Per-Lambda configuration files
- `config/sls/environment.yml` - Environment variables
- `config/sls/resources.yml` - Lambda roles and permissions

### Mockery Configuration
- `.mockery.yml` - Project-level mockery configuration (used by `new:domain` action)
- `.mockery.base.yml` - Base configuration for `draft mockery` command (shared settings)
- `.mockery.pkg.yml` - Package-specific mockery configurations (merged with base config by `draft mockery`)

## Code Standards

### Linting
- Uses `revive` with `revive.toml` configuration
- `go vet` for standard Go checks
- Enforced in CI via GitHub Actions

### Commit Conventions
- Commitizen for commit message formatting
- Commits must follow conventional commit format
- Enforced in PR checks (`.github/workflows/pull_request.yml`)
- PR titles must also follow conventional commit format

### Code Organization
- `internal/pkg/` contains reusable utilities (files, exec, crypto, AWS, build, log, dirs, etc.)
  - `internal/pkg/dirs/` provides directory operations including:
    - `Walk()` function that respects `.gitignore` patterns by default
    - Parses `.gitignore` from root directory and filters files/directories automatically
    - Supports rooted patterns (`/dist`), directory-only patterns (`node_modules/`), negation patterns (`!important.txt`), and glob patterns (`**/*.yml`)
    - Optional `skipGitignore` parameter to disable `.gitignore` filtering
    - Used by mockery command to automatically skip vendor, node_modules, and other ignored directories
  - `internal/pkg/migrateconfig/` provides database configuration utilities:
    - Reads and parses `.local-migrate-config.yml` from project root
    - Extracts `migrations.databases` configuration
    - Filters databases by group (excludes `group: 'test'`)
    - Provides `ToPascalCase()` for converting snake_case to PascalCase
    - Formats database names for display in forms
- `internal/dtos/` for data transfer objects passed between commands, forms, and actions
- `internal/data/` for global state (flags, metadata, placeholder tags)
- `cmd/commands/internal/common/` for command-level shared code
- `internal/project/` for project detection and name normalization

### Placeholder Tag System
The codebase uses placeholder comments to mark insertion points for incremental additions:
- `data.NextImportTag` (`//draft:next-import`) - Marks where to insert new imports in `deps.go`
- `data.NextLambdaImportTag` (`#draft:next-lambda-import`) - Marks where to insert new lambda references in `serverless.yml`
- `data.NextDbModelTag` (`//draft:next-db-model`) - Marks where to insert new database models in Postgres provider test migrations

Post-create hooks use string replacement to insert new content while preserving the tag for future additions. This pattern avoids AST parsing for file manipulation. The tags are defined in `internal/data/tags.go`.

## Nix Integration
The project uses Nix flakes (`flake.nix`) for:
- Version management of Nix modules
- Script `update-nix-hashes.sh` updates vendor hashes for the Nix flake
- `internal/pkg/nix-metadata/` checks Nix module versions on startup
- Nix metadata is stored in `$HOME/.config/.nix-metadata.json` containing system info, git user info, and version information
- The tool prompts for Nix module updates if version is outdated (checked every 24 hours, prompts every 2 hours if declined)

## Key Implementation Patterns

### Type-Based Dispatch
Lambda and domain creation use switch statements on type strings rather than polymorphism:
```go
switch nl.input.LambdaType {
case "plain": return nl.createPlain()
case "http":  return nl.createHttp()
// etc.
}
```

### Post-Create File Updates
New lambdas are added to existing files via string replacement:
1. `addToDepsGo()` - Inserts import into `deps.go` before placeholder tag
2. `addToServerlessYAML()` - Adds lambda reference to `serverless.yml`
3. `format()` - Runs goimports/gofmt on modified files
4. `restoreDepsTag()` - Re-adds placeholder tag for next insertion

New domains have a different post-create flow:
1. `postgresModels()` (Postgres only) - Adds domain DAOs to provider test migrations using `NextDbModelTag`
2. `mockery()` - Creates `.mockery.pkg.yml` files and runs mockery action to generate mocks
   - Reuses `internal/actions/mockery` for concurrent execution
   - Passes context for cancellation support
   - Runs with jobsNum=2 for service and repository
3. `format()` - Runs goimports/gofmt on generated files

### Global State Usage
Minimal global state in `internal/data/`:
- `data.Flags` - CLI flags (WorkingDir, Debug, TTY, NoSentry)
- `data.Meta` - Project metadata loaded from `go.mod` at startup
- Placeholder tag constants for file manipulation
