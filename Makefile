.DEFAULT_GOAL := help
help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Install dependencies
	go mod download

lint: ## Run linters
	golangci-lint run

fmt: ## Fix formatting issues
	goimports -w .

test: ## Run tests
	go test -v -race ./...

test-coverage: ## Launch tests coverage and send it to coverall
	mkdir -p .cover
	go test -coverprofile .cover/cover.out ./...

test-coverage-html: ## Create a code coverage report in HTML
	mkdir -p .cover
	go test -coverprofile .cover/cover.out ./...
	go tool cover -html .cover/cover.out

release: ## Release
	goreleaser --rm-dist --release-notes CHANGELOG.md --skip-validate

upx: ## Compact artifacts
	upx dist/*/*

clean: ## Clean up
	rm -rf .cover dist