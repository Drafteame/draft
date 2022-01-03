help: ## Display this help screen.
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install: ## Install all dependencies for Go and Serverless framework
	@go mod download

precommit_install: ## Install precommit actions
	@pre-commit install && pre-commit install --hook-type commit-msg

test: ## Execute unit testing
	@go test -v -race $(shell go list ./... | grep -v cmd) --cover -tags=unit -short

lint: ## Lint code
	@golangci-lint run ./...

fmt:
	@gofmt -s -w -l $(shell go list -f {{.Dir}} ./... | grep -v /vendor/)