# Draft

![draft-logo](https://github.com/user-attachments/assets/bfa3d186-5576-4d85-9366-61cbe792e8f8)

A CLI tool to scaffold services, Lambda functions, and domain layers for Draftea Backend Services. Draft generates standardized project structures following best practices for serverless applications built with Go.

---

## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
  - [Creating a Service](#creating-a-service)
  - [Creating a Lambda Function](#creating-a-lambda-function)
  - [Creating a Domain Layer](#creating-a-domain-layer)
  - [Mock Generation](#mock-generation)
  - [Invoking Lambdas Locally](#invoking-lambdas-locally)
  - [Local Development](#local-development)
  - [Sentry Management](#sentry-management)
  - [Database Tunnel Management](#database-tunnel-management)
  - [Global Flags](#global-flags)
- [Development](#development)
  - [Project Structure](#project-structure)
  - [Architecture Overview](#architecture-overview)
  - [Adding a New Command](#adding-a-new-command)
  - [Template System](#template-system)
  - [Development Workflow](#development-workflow)
  - [Code Standards](#code-standards)
- [Contributing](#contributing)
- [License](#license)

---

## Requirements

- **Go**: 1.23.x or higher
- **Task**: For running development commands (optional)
- **Husky**: For git hooks (optional)
- **AWS CLI**: Required for `db:connect` commands (`aws` in PATH)
- **Session Manager Plugin**: Required for SSM port-forwarding tunnels ([installation guide](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html))

## Installation

Install Draft globally using Go:

```bash
go install github.com/Drafteame/draft@latest
```

Verify the installation:

```bash
draft --version
```

---

## Usage

### Creating a Service

Create a new serverless service with boilerplate code, configuration files, and initial Lambda functions.

#### Basic Usage

Navigate to the service folder in your monorepo and run:

```bash
draft new:service
```

This launches an interactive form to configure your service:

![Screenshot 2024-11-18 at 5 52 52 p m](https://github.com/user-attachments/assets/bdf7d12c-2b6e-4b40-bfc7-2dbfb20abb16)

#### Flags

```bash
# Specify working directory
draft new:service -w path/to/project

# Create in legacy path structure (for main-api or game-engine)
draft new:service -l path/to/legacy/service

# Use Uber Dig for dependency injection
draft new:service --use-dig

# Disable TTY mode (useful in CI/CD)
draft new:service -t=false
```

#### What Gets Created

- `serverless.yml` - Serverless Framework configuration
- `package.json` - Node.js dependencies for Serverless plugins
- `config/` - Application and Serverless configuration files
- `cmd/` - Initial Lambda function (HTTP, Cron, SQS, etc.)
- `deps.go` - Dependency injection setup

---

### Creating a Lambda Function

Add a new Lambda function to an existing service.

#### Basic Usage

Navigate to your project root and run:

```bash
draft new:lambda
```

![Screenshot 2024-11-18 at 5 54 17 p m](https://github.com/user-attachments/assets/902a4e37-af2b-44b1-8851-0b98619c1aed)

#### Lambda Types

Draft supports multiple Lambda trigger types:

| Type | Description | Use Case |
|------|-------------|----------|
| **HTTP** | API Gateway triggered | REST APIs, webhooks |
| **Cron** | EventBridge scheduled | Periodic tasks, cleanup jobs |
| **SQS** | SQS queue processing | Async processing, decoupling |
| **SNS+SQS** | SNS topic + SQS queue | Pub/sub with queue buffering |
| **Plain** | Generic event source | Custom event sources |
| **Custom** | User-defined type | Custom event sources with specific paths |

#### Flags

```bash
# Specify working directory
draft new:lambda -w path/to/project

# Create in legacy path structure
draft new:lambda -l path/to/legacy/service

# Use Uber Dig for dependency injection
draft new:lambda --use-dig
```

#### What Gets Created

For an HTTP Lambda:
- `cmd/http/<lambda-name>/main.go` - Lambda entry point
- `cmd/http/<lambda-name>/handler/` - Business logic
- `cmd/http/<lambda-name>/handler/dtos/` - Request/response models
- `cmd/http/<lambda-name>/handler/worker/` - Core processing logic
- `cmd/http/<lambda-name>/lambda-config.yml` - Lambda-specific config

---

### Creating a Domain Layer

Generate a complete domain layer following Domain-Driven Design principles with support for both interactive and non-interactive modes.

#### Basic Usage

**Interactive Mode** (prompts for all values):
```bash
draft new:domain
```

**Non-Interactive Mode** (CI/CD friendly):
```bash
# Postgres domain
draft new:domain \
  --domain-path users \
  --db-type postgres \
  --table-name public.users \
  --db-prefix usr \
  --db-name general

# DynamoDB domain
draft new:domain \
  --domain-path products \
  --db-type dynamo \
  --table-name ProductsTable
```

**Mixed Mode** (some flags, some prompts):
```bash
draft new:domain --domain-path orders --db-type postgres
# Prompts only for table-name, db-prefix, and db-name
```

#### Database Support

Draft supports multiple database backends:

| Database | Features |
|----------|----------|
| **Postgres** | Full CRUD, search with filters/pagination, repository builders, DAOs, domain models |
| **DynamoDB** | Simplified repository pattern, optimized for NoSQL |

#### Database Configuration

For Postgres domains, available databases are **dynamically loaded** from `.local-migrate-config.yml`:

1. Reads `migrations.databases` from project root configuration
2. Filters out test databases (`group: 'test'`)
3. Presents remaining databases as options
4. Converts database names to PascalCase for provider functions

**Example**: If `.local-migrate-config.yml` contains:
```yaml
migrations:
  databases:
    general:
      folder: 'postgres'
    general_test:
      group: 'test'  # This will be filtered out
    user_preferences:
      folder: 'user-preferences'
    games_core:
      folder: 'games-core'
```

The command will:
- Show options: `General`, `User Preferences`, `Games Core`
- Generate provider calls: `ProvideGeneral`, `ProvideUserPreferences`, `ProvideGamesCore`

#### Flags

| Flag | Short | Description | Example |
|------|-------|-------------|---------|
| `--domain-path` | `-p` | Path to domain folder | `users`, `auth/sessions` |
| `--db-type` | | Database type | `postgres`, `dynamo` |
| `--table-name` | | Database table name | `public.users`, `ProductsTable` |
| `--db-prefix` | | ID prefix for Postgres (3 chars) | `usr`, `ord`, `prd` |
| `--db-name` | | Database name (from config) | `general`, `user_preferences` |
| `--working-dir` | `-w` | Working directory | `path/to/project` |

**Note**: The `--db-name` flag accepts snake_case database names as they appear in `.local-migrate-config.yml`. The tool automatically converts them to PascalCase for provider function names.

#### What Gets Created (Postgres)

```
domains/<domain-name>/
├── domain/
│   ├── domain.go          # Domain models
│   ├── errors.go          # Domain-specific errors
│   └── options/           # Search, filter, pagination options
├── service/
│   ├── interfaces.go      # Service interfaces
│   ├── service.go         # Service implementation
│   ├── create.go          # Create operation
│   ├── get.go             # Get operation
│   ├── update.go          # Update operation
│   ├── delete.go          # Delete operation
│   ├── search.go          # Search with filters
│   ├── search_one.go      # Search single record
│   ├── provide.go         # Dependency injection
│   └── *_test.go          # Unit tests
└── repository/
    ├── interfaces.go      # Repository interfaces
    ├── repository.go      # Repository implementation
    ├── create.go          # Create operation
    ├── get.go             # Get operation
    ├── update.go          # Update operation
    ├── delete.go          # Delete operation
    ├── search.go          # Search with filters
    ├── search_one.go      # Search single record
    ├── provide.go         # Dependency injection
    ├── builders/          # Query builders
    ├── daos/              # Data access objects
    └── *_test.go          # Unit tests
```

#### What Gets Created (DynamoDB)

```
domains/<domain-name>/
├── service/
│   ├── interfaces.go
│   ├── service.go
│   └── provider.go
└── repository/
    ├── interfaces.go
    ├── repository.go
    └── provider.go
```

---

### Mock Generation

Generate mock implementations for interfaces across your codebase using [Mockery](https://github.com/vektra/mockery) with configuration merging support.

#### Overview

The `mockery` command enables concurrent mock generation by merging a base configuration (`.mockery.base.yml`) with package-specific configurations (`.mockery.pkg.yml`). This pattern allows you to:
- Define common settings once in a base config
- Override or extend settings per package
- Generate mocks for multiple packages concurrently
- Track progress and errors for each package

#### Basic Usage

```bash
# Generate mocks for all .mockery.pkg.yml files in the project
draft mockery

# Generate mocks for specific package configs
draft mockery services/user/.mockery.pkg.yml services/auth/.mockery.pkg.yml

# Run with custom concurrency (default: 5)
draft mockery --jobs-num 10

# Dry run - validate configs without executing mockery
draft mockery --dry

# Run only for packages with modified files (git diff with main)
draft mockery --git-mod

# Combine flags for targeted validation in CI/CD
draft mockery --git-mod --dry --jobs-num 10
```

#### Configuration Files

**Base Configuration** (`.mockery.base.yml`)
```yaml
# Shared settings applied to all packages
with-expecter: true
dir: "{{.InterfaceDirRelative}}/mocks"
mockname: "Mock{{.InterfaceName}}"
filename: "mock_{{.InterfaceName}}.go"
outpkg: mocks
```

**Package Configuration** (`services/user/.mockery.pkg.yml`)
```yaml
# Package-specific settings
packages:
  github.com/myorg/myrepo/services/user/domain:
    config:
      all: true
      dir: services/user/domain/mocks
      filename: "mock_{{.InterfaceName}}.go"
      inpackage: false
```

The command will:
1. Load base config from `.mockery.base.yml`
2. Merge it with each `.mockery.pkg.yml` file (package config takes precedence)
3. Create temporary merged configs (`.mockery.tmp.*.yml`)
4. Execute mockery concurrently for each package
5. Clean up temporary files
6. Report success/failure statistics

#### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--jobs-num` | `-j` | Number of concurrent mockery jobs | `5` |

#### Example Output

```
✓ [✓] [ 3 / 5] services/auth (2.34s)
✗ [✗] [ 4 / 5] services/payment (1.12s)
✓ [✓] [ 5 / 5] services/user (3.45s)

✗ Failed: 1/5 packages (8.91s)
Failed packages:
  • services/payment/.mockery.pkg.yml
    mockery failed: interface NotFound not found
```

#### Integration with `new:domain`

When creating domains with `draft new:domain`, mocks are automatically generated:
- Creates `.mockery.pkg.yml` files for service and repository packages
- Calls the mockery action directly (reuses `internal/actions/mockery`)
- Runs concurrently with `jobsNum=2` for service and repository
- Provides progress reporting and error handling
- Supports cancellation with Ctrl+C
- Mocks are placed in `domain/service/mocks` and `domain/repository/mocks`

The domain command benefits from the same mockery infrastructure used by the standalone `draft mockery` command, including concurrent execution, progress tracking, and proper error handling.

---

### Invoking Lambdas Locally

Test Lambda functions locally without deploying to AWS.

#### Basic Usage

```bash
# Simple invocation (useful for HTTP Lambdas)
draft invoke path/to/lambda

# With event data inline
draft invoke --body '{"key": "value"}' path/to/lambda

# With event data from file
draft invoke --body-file event.json path/to/lambda
```

#### Example Event File

```json
{
  "httpMethod": "POST",
  "path": "/api/users",
  "body": "{\"name\":\"John Doe\"}",
  "headers": {
    "Content-Type": "application/json"
  }
}
```

---

### Local Development

Set up and manage local development environment.

#### Setup Local Environment

```bash
draft local:setup
```

This creates:
- Docker Compose configuration
- Local database setup
- Environment variables

#### Database Migrations

```bash
# Run migrations up
draft local:migrate:up

# Rollback migrations
draft local:migrate:down

# Force migration version
draft local:migrate:force <version>
```

---

### Sentry Management

Manage Sentry projects for error tracking.

#### Create Sentry Project

```bash
draft sentry:project:create
```

#### Delete Sentry Project

```bash
draft sentry:project:delete
```

---

### Database Tunnel Management

Manage SSM port-forwarding tunnels to remote databases. Tunnels are opened via AWS Session Manager, connecting your local machine to RDS (Postgres), ElastiCache (Redis), or DocumentDB (MongoDB) instances in private subnets.

#### Prerequisites

- AWS CLI configured with the appropriate profiles
- [Session Manager Plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html) installed
- Config file at `~/.draft/dbconnect.yml` (see [Config File Format](#config-file-format) below)

#### Subcommands

| Command | Description |
|---------|-------------|
| `draft db:connect list` | List all available connections from config |
| `draft db:connect status` | Show status of active tunnels |
| `draft db:connect start <type> <name>` | Open a tunnel to a database |
| `draft db:connect stop [type] [name]` | Close one or more tunnels |

#### List Available Connections

```bash
draft db:connect list
```

Displays a table of all connections that can be derived from the config file by combining every service with every environment:

```
ENV   TYPE      NAME                    HOST                                              REMOTE PORT   LOCAL PORT
----------------------------------------------------------------------------------------------------------------------
dev   mongo     main-dev                draftea-dev-maincluster.cluster-xxx.docdb.com     27017         56200
dev   postgres  turbo-dev               turbo-dev.cluster-xxx.rds.amazonaws.com           5432          56024
prod  postgres  turbo-prod              turbo-prod.cluster-xxx.rds.amazonaws.com          5432          56011
...
```

#### Check Tunnel Status

```bash
draft db:connect status
```

Shows all tunnels that have been started, including whether the process is still alive:

```
TYPE      NAME        STATUS   PID     ADDRESS
---------------------------------------------------------
postgres  turbo-dev   up       60452   localhost:56024
redis     turbo-dev   down     -       -
```

Dead tunnels (e.g. killed by idle timeout) are automatically cleaned up from state on the next `status` or `start` call.

#### Start a Tunnel

```bash
# Start with the port defined in config
draft db:connect start postgres turbo-dev

# Start with a custom local port
draft db:connect start postgres turbo-dev --port 15432
```

The connection name must end with `-dev` or `-prod`. The tunnel runs as a background process; the command returns once the port is confirmed to be listening (up to 10 seconds).

#### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--port` | `-p` | Override the local port for this connection |

#### Stop Tunnels

```bash
# Stop a specific tunnel
draft db:connect stop postgres turbo-dev

# Stop all tunnels of a given type
draft db:connect stop postgres

# Stop all active tunnels
draft db:connect stop --all
```

#### Config File Format

Create `~/.draft/dbconnect.yml`:

```yaml
defaults:
  postgres:
    remote_port: 5432
  redis:
    remote_port: 6379
  mongo:
    remote_port: 27017

environments:
  dev:
    bastion:
      target: i-0xxxxxxxxxxxxxxxxx   # EC2 instance ID of the bastion
      profile: my-aws-dev-profile
      region: us-east-2
    clusters:
      rds:   xxxxxxxxxx.us-east-2.rds.amazonaws.com      # RDS cluster suffix
      cache: xxxxxxxxxx.use2.cache.amazonaws.com          # ElastiCache suffix
      docdb: xxxxxxxxxx.us-east-2.docdb.amazonaws.com    # DocumentDB suffix

  prod:
    bastion:
      target: i-0yyyyyyyyyyyyyyyyy
      profile: my-aws-prod-profile
      region: us-east-2
    clusters:
      rds:   yyyyyyyyyy.us-east-2.rds.amazonaws.com
      cache: yyyyyyyyyy.use2.cache.amazonaws.com
      docdb: yyyyyyyyyy.us-east-2.docdb.amazonaws.com

connections:
  postgres:
    instances:
      - name: my-service
        local_ports:
          dev: 56000
          prod: 56001

  redis:
    instances:
      - name: my-cache
        local_ports:
          dev: 56100
          prod: 56101

  mongo:
    instances:
      - name: main
        local_ports:
          dev: 56200
          prod: 56201
```

**Host construction per engine:**

| Type | Pattern |
|------|---------|
| `postgres` | `{service}-{env}.cluster-{clusters.rds}` |
| `redis` | `{service}-{env}.{clusters.cache}` |
| `mongo` | `draftea-{env}-maincluster.cluster-{clusters.docdb}` |

#### Runtime State

Active tunnels are persisted in `~/.draft/dbconnect.state.json`. This file is managed automatically — you do not need to edit it manually. When a tunnel dies (e.g. AWS idle timeout), the state file is **not** updated automatically; the stale entry is cleaned up the next time you run `status` or `start`.

---

### Global Flags

These flags are available for all commands:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--working-dir` | `-w` | Working directory for the command | Current directory |
| `--debug` | `-d` | Enable debug mode with verbose logging | `false` |
| `--tty` | `-t` | Enable TTY mode for interactive prompts | `true` |

#### Example

```bash
draft new:lambda -w /path/to/project -d -t=false
```

---

## Development

### Project Structure

```
draft/
├── cmd/
│   ├── main.go                    # Application entry point
│   └── commands/                  # Cobra command definitions
│       ├── root.go                # Root command with global flags
│       ├── newservice/            # Service creation command
│       ├── newlambda/             # Lambda creation command
│       ├── newdomain/             # Domain creation command
│       ├── mockery/               # Mock generation command
│       ├── local/                 # Local development commands
│       ├── sentry/                # Sentry management commands
│       └── internal/common/       # Shared command utilities
├── internal/
│   ├── actions/                   # Business logic for commands
│   │   ├── newservice/            # Service creation logic
│   │   ├── newlambda/             # Lambda creation logic
│   │   ├── newdomain/             # Domain creation logic
│   │   ├── mockery/               # Mock generation logic
│   │   ├── local/                 # Local development logic
│   │   └── sentry/                # Sentry management logic
│   ├── forms/                     # Interactive CLI forms
│   │   ├── newservice/            # Service configuration forms
│   │   ├── newlambda/             # Lambda configuration forms
│   │   └── newdomain/             # Domain configuration forms
│   ├── templates/                 # Code generation templates
│   │   ├── templates.go           # Template loading infrastructure
│   │   ├── service_templates.go   # Service templates
│   │   ├── lambda_*.go            # Lambda type templates
│   │   ├── domains_*.go           # Domain templates
│   │   └── tmpl/                  # Embedded template files
│   │       ├── sls/               # Serverless templates
│   │       └── domain/            # Domain templates
│   ├── dtos/                      # Data transfer objects
│   │   ├── service_input.go       # Service configuration
│   │   ├── lambda_input.go        # Lambda configuration
│   │   ├── domain_input.go        # Domain configuration
│   │   └── file_entry.go          # File creation DTO
│   ├── data/                      # Global state and constants
│   │   ├── flags.go               # Global CLI flags
│   │   ├── meta.go                # Project metadata
│   │   ├── tags.go                # Template placeholder tags
│   │   ├── database.go            # Database type constants
│   │   └── paths.go               # Common path constants
│   └── pkg/                       # Reusable utilities
│       ├── files/                 # File I/O operations
│       ├── dirs/                  # Directory operations
│       ├── exec/                  # Command execution
│       ├── inputs/                # Form input helpers
│       ├── log/                   # Logging utilities
│       ├── format/                # Code formatting
│       ├── constants/             # Shared constants
│       ├── migrateconfig/         # Database config utilities
│       └── ...                    # Other utilities
├── Taskfile.yml                   # Task runner configuration
├── go.mod                         # Go module definition
└── README.md                      # This file
```

---

### Architecture Overview

Draft follows a layered architecture with clear separation of concerns:

```
┌─────────────────────────────────────────────┐
│             Command Layer                    │
│         (cmd/commands/*)                    │
│  - Parse flags                               │
│  - Call forms                                │
│  - Execute actions                           │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│              Forms Layer                     │
│         (internal/forms/*)                  │
│  - Interactive prompts                       │
│  - Input validation                          │
│  - Populate DTOs                             │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│             Action Layer                     │
│         (internal/actions/*)                │
│  - Business logic                            │
│  - Lifecycle: preCreate → exec → postCreate │
│  - File generation                           │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│            Template Layer                    │
│         (internal/templates/*)              │
│  - Embedded templates                        │
│  - Template loading                          │
│  - Content generation                        │
└─────────────────────────────────────────────┘
```

#### Key Patterns

1. **Action Lifecycle**: All creation actions follow a 3-phase lifecycle:
   - `preCreate()` - Validation and setup
   - `exec()` - Core file creation logic
   - `postCreate()` - Formatting and cleanup

2. **Template Embedding**: Templates are embedded at compile time using `//go:embed`

3. **Functional Options**: Forms use functional options pattern for flexibility

4. **Type-Based Dispatch**: Lambda and domain types use switch statements on type strings

---

### Adding a New Command

Follow these steps to add a new command to Draft:

#### 1. Create Command Structure

Create a new directory under `cmd/commands/`:

```bash
mkdir -p cmd/commands/mynewcommand
```

#### 2. Define the Command

Create `cmd/commands/mynewcommand/mynewcommand.go`:

```go
package mynewcommand

import (
    "github.com/spf13/cobra"
    "github.com/Drafteame/draft/cmd/commands/internal/common"
    "github.com/Drafteame/draft/internal/actions/mynewcommand"
    "github.com/Drafteame/draft/internal/data"
    "github.com/Drafteame/draft/internal/dtos"
    "github.com/Drafteame/draft/internal/forms/mynewcommand"
    "github.com/Drafteame/draft/internal/pkg/log"
)

var myNewCommandCmd = &cobra.Command{
    Use:   "my:new:command",
    Short: "Brief description of the command",
    Long:  "Detailed description of what the command does",
    Run:   run,
}

func init() {
    // Add command-specific flags
    myNewCommandCmd.Flags().Bool("my-flag", false, "Description of flag")
}

func run(cmd *cobra.Command, _ []string) {
    // Change to working directory if specified
    common.ChDir(cmd)

    // Load project metadata
    data.LoadMeta()

    // Extract flags
    myFlag := common.GetBoolFlag(cmd, "my-flag")

    // Create input DTO
    input := dtos.MyNewCommandInput{
        MyFlag: myFlag,
    }

    // Run form to collect user input
    if err := mynewcommand.GetForm(&input); err != nil {
        log.Exitf(1, "failed to collect input: %s", err.Error())
    }

    // Execute action
    if err := mynewcommand.New(input).Exec(); err != nil {
        log.Exitf(1, "failed to execute command: %s", err.Error())
    }

    log.Success("Command completed successfully")
}

func GetCmd() *cobra.Command {
    return myNewCommandCmd
}
```

#### 3. Register the Command

Add to `cmd/commands/root.go`:

```go
import (
    "github.com/Drafteame/draft/cmd/commands/mynewcommand"
)

func init() {
    // ... existing commands ...
    rootCmd.AddCommand(mynewcommand.GetCmd())
}
```

#### 4. Create DTO

Create `internal/dtos/mynewcommand_input.go`:

```go
package dtos

type MyNewCommandInput struct {
    PackageName string
    MyFlag      bool
    // Add fields needed for your command
}
```

#### 5. Create Form

Create `internal/forms/mynewcommand/mynewcommand.go`:

```go
package mynewcommand

import (
    "github.com/Drafteame/draft/internal/data"
    "github.com/Drafteame/draft/internal/dtos"
    "github.com/Drafteame/draft/internal/pkg/inputs"
)

func GetForm(input *dtos.MyNewCommandInput) error {
    input.PackageName = data.Meta.PackageName

    // Add form fields
    err := inputs.Text("Field Name:",
        inputs.WithValue(&input.FieldName),
        inputs.WithDescription[string]("Help text"),
        inputs.WithValidation(func(val string) error {
            if val == "" {
                return errors.New("field cannot be empty")
            }
            return nil
        }),
    )

    return err
}
```

#### 6. Create Action

Create `internal/actions/mynewcommand/mynewcommand.go`:

```go
package mynewcommand

import (
    "github.com/Drafteame/draft/internal/dtos"
)

type MyNewCommand struct {
    input dtos.MyNewCommandInput
}

func New(input dtos.MyNewCommandInput) *MyNewCommand {
    return &MyNewCommand{
        input: input,
    }
}

func (m *MyNewCommand) Exec() error {
    if err := m.preCreate(); err != nil {
        return err
    }

    if err := m.exec(); err != nil {
        return err
    }

    return m.postCreate()
}

func (m *MyNewCommand) preCreate() error {
    // Validation and setup
    return nil
}

func (m *MyNewCommand) exec() error {
    // Core logic
    return nil
}

func (m *MyNewCommand) postCreate() error {
    // Cleanup and formatting
    return nil
}
```

#### 7. Add Templates (if needed)

If your command generates files, create templates:

1. Create template directory: `internal/templates/tmpl/mycommand/`
2. Add template files with `.tmpl` extension
3. Create template loader in `internal/templates/`

---

### Template System

#### Creating Templates

Templates use Go's `text/template` syntax:

```go
// File: internal/templates/tmpl/mycommand/file.go.tmpl
package {{.PackageName}}

// Generated with Draft
func {{.FunctionName}}() {
    // Your code here
}
```

#### Loading Templates

Create a template loader:

```go
//go:embed tmpl/mycommand
var myCommandTemplates embed.FS

func loadMyCommandTemplate(data any) ([]byte, error) {
    name := "mycommand/file.go"
    path := "tmpl/mycommand/file.go.tmpl"
    return loadTemplate(name, path, data, myCommandTemplates)
}
```

#### Using Templates in Actions

```go
func (m *MyNewCommand) exec() error {
    content, err := templates.LoadMyCommandTemplate(m.input)
    if err != nil {
        return err
    }

    return files.Create("output/file.go", content)
}
```

---

### Development Workflow

#### Prerequisites

Install development dependencies:

```bash
# Install Task runner
go install github.com/go-task/task/v3/cmd/task@latest

# Install Husky for git hooks
task hooks
```

#### Common Tasks

```bash
# Format code
task format

# Run linters
task lint

# Install git hooks
task hooks
```

#### Building

```bash
# Install locally
go install

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o draft-linux

# Build for all platforms (uses goreleaser)
goreleaser build --snapshot --clean
```

#### Testing

Currently, there is no test suite. When adding tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/actions/newlambda/...
```

---

### Code Standards

#### Linting

The project uses `revive` and `go vet`:

```bash
task lint
```

Configuration is in `revive.toml`.

#### Code Formatting

Use `goimports-reviser` for imports and formatting:

```bash
task format
```

#### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new lambda type
fix: correct directory creation bug
docs: update README with examples
refactor: extract common flag helpers
chore: update dependencies
```

Enforced by Commitizen (`.cz.toml`) and git hooks.

#### Pull Requests

PR titles must follow conventional commit format:

```
feat(lambda): add support for custom event sources
fix(domain): correct postgres template paths
docs: improve development guide in README
```

Enforced by `.github/workflows/pull_request.yml`.

---

## Contributing

We welcome contributions! Please follow these guidelines:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feat/my-new-feature`
3. **Make your changes** following the code standards
4. **Run linters**: `task lint`
5. **Format code**: `task format`
6. **Commit with conventional commits**: `git commit -m "feat: add new feature"`
7. **Push to your fork**: `git push origin feat/my-new-feature`
8. **Open a Pull Request**

### Development Setup

```bash
# Clone the repository
git clone https://github.com/Drafteame/draft.git
cd draft

# Install dependencies
go mod download

# Install git hooks
task hooks

# Build and install locally
go install

# Verify installation
draft --version
```

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Copyright (c) 2025 Draftea

---

## Support

For issues, questions, or contributions:

- **GitHub Issues**: [https://github.com/Drafteame/draft/issues](https://github.com/Drafteame/draft/issues)
- **Documentation**: See `CLAUDE.md` for detailed architecture information

---

**Built with ❤️ by the Draftea team**
