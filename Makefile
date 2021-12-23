help: ## Display this help screen.
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install: ## Install all dependencies for Go and Serverless framework
	@go mod download

test: ## Execute unit testing
	@go test -v -race `go list ./... | grep -v cmd` --cover

lint: ## Lint code
	@golangci-lint run ./...