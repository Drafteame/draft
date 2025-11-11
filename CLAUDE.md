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
```

### Global Flags
- `-w, --working-dir`: Specify working directory
- `-d, --debug`: Enable debug mode
- `-t, --tty`: TTY mode (default: true)

## Architecture

### Command Structure
Commands are registered in `main.go` using Cobra. Each command lives in `cmd/commands/` with its own subdirectory:
- `newservice` - Service scaffolding
- `newlambda` - Lambda scaffolding
- `newdomain` - Domain layer scaffolding
- `local/invoke` - Local Lambda invocation
- `local/migrate/*` - Database migration commands
- `sentry/deleteproject` - Sentry project management
- `config` - Configuration management

### Action Layer
The `internal/actions/` directory contains the business logic for each command. Actions are responsible for:
- Gathering user input via interactive forms (`internal/forms/`)
- Generating files from templates (`internal/templates/`)
- Executing post-creation tasks (formatting, running external tools)

Key action patterns:
- `New()` factory function creates action instances with input DTOs
- `Exec()` orchestrates the entire action lifecycle
- `preCreate()`, `exec()`, `postCreate()` for lifecycle hooks

### Template System
Templates are embedded at compile time using `//go:embed` in `internal/templates/`. Template types:
- **Serverless templates** (`tmpl/sls/`): Service boilerplate, Lambda functions (HTTP, Cron, SQS, SNS+SQS, Plain, Custom)
- **Domain templates** (`tmpl/domain/`): Domain-driven design layers (domain, service, repository with Postgres and DynamoDB support)

Templates use Go's `text/template` and are loaded via `loadTemplate()` helper.

### Lambda Types
The tool generates different Lambda structures based on trigger type:
- **HTTP**: API Gateway triggered, includes routing and DTOs
- **Cron**: EventBridge scheduled execution
- **SQS**: Queue-based processing with idempotency
- **SNS+SQS**: Combined SNS topic and SQS queue processing
- **Plain**: Generic Lambda for custom event sources
- **Custom**: User-defined custom event sources with configurable type path and optional idempotency

### Domain Generation
When creating domains with `draft new:domain`, the tool prompts for:
- Domain path (automatically prefixed with `domains/` if not present)
- Database type selection (Postgres or DynamoDB)
- Database-specific configuration (table names, prefixes, etc.)

Domain structure varies by database type:
- **Postgres**: Full CRUD with service layer, repository layer, builders, DAOs, providers, and domain models with search/filter capabilities
- **DynamoDB**: Simplified structure with repository and service layers optimized for NoSQL patterns

### Project Detection
`internal/project/` contains utilities for:
- Extracting Go module names from `go.mod`
- Normalizing service names (kebab-case for folders, snake_case for packages)
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
`internal/forms/` provides interactive CLI prompts using `charmbracelet/huh`:
- Multi-step forms for service/lambda/domain creation
- Input validation and conditional fields
- Support for TTY and non-TTY modes

## Configuration Files

### Pkl Configuration
Services use Pkl (configuration language) for app config:
- `config/app/app.pkl` - Application configuration
- `config/app/modules.pkl` - Pkl module imports

### Serverless Framework
- `serverless.yml` - Service definition and Lambda configuration
- `lambda-config.yml` - Per-Lambda configuration files
- `config/sls/environment.yml` - Environment variables
- `config/sls/iam.yml` - IAM permissions

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
- `internal/pkg/` contains reusable utilities (files, exec, crypto, AWS, etc.)
- `internal/dtos/` for data transfer objects
- `internal/data/` for shared data structures and flags
- `cmd/commands/internal/common/` for command-level shared code

## Nix Integration
The project uses Nix flakes (`flake.nix`) for:
- Version management of Nix modules
- Script `update-nix-hashes.sh` updates vendor hashes
- `internal/pkg/version/nix/` checks Nix module versions on startup

## Development Notes

### Current Branch State
The current working branch may contain merge conflicts or work in progress. Key files to check:
- `internal/actions/newdomain/exec.go` - Contains logic for creating domains with both Postgres and DynamoDB support
- `internal/templates/domains_*.go` - Template loaders for domain generation

### Recent Features
- **DynamoDB Support**: Added in commit `fda92fc`. Domains can now be generated with DynamoDB as the backing database
- **Custom Lambda Type**: New lambda type allowing custom event sources with configurable type paths
