.PHONY: all build run run-http test test-race fmt fmt-check vet imports imports-check lint golangci-lint generate tidy tidy-check coverage coverage-summary clean image-build image-run image-run-http image-run-http-read-only help

BINARY_NAME := container-runtime-mcp
BUILD_FLAGS := -o $(BINARY_NAME) .

IMAGE_NAME := container-runtime-mcp
CONTAINER_RUNTIME ?= docker
HTTP_ADDR := 127.0.0.1:8080
RUNTIME_SOCKET := /var/run/docker.sock

.DEFAULT_GOAL := help

## help: Show this help message.
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z0-9_-]+:.*?## / {printf "  %-15s %s\n", $$1, substr($$0, index($$0, "##") + 3)}' $(MAKEFILE_LIST)

all: lint test build ## Run lint, test, and build.

build: ## Build the binary.
	go build $(BUILD_FLAGS)

run: ## Run the application with stdio transport (default).
	go run .

run-http: ## Run the application with HTTP transport.
	CONTAINER_RUNTIME_MCP_TRANSPORT=http CONTAINER_RUNTIME_MCP_HTTP_ADDR=$(HTTP_ADDR) go run .

test: ## Run all tests.
	go test ./...

test-race: ## Run all tests with the race detector.
	go test -race ./...

coverage: ## Run tests with coverage profile.
	go test -covermode=atomic -coverprofile=coverage.out ./...
	@# Remove generated mock files from coverage report.
	@sed -i.bak '/\/mock\//d' coverage.out && rm -f coverage.out.bak

coverage-summary: ## Print coverage summary.
	go tool cover -func=coverage.out

fmt: ## Format all Go code with gofmt.
	go fmt ./...

fmt-check: ## Check if Go code is formatted (read-only).
	@output=$$(gofmt -l .); \
	if [ -n "$$output" ]; then \
		echo "gofmt would modify:"; \
		echo "$$output"; \
		exit 1; \
	fi

vet: ## Run go vet on all packages.
	go vet ./...

imports: ## Organize imports with goimports.
	GOWORK=off go -C tools run golang.org/x/tools/cmd/goimports -w ..

imports-check: ## Check if imports are organized (read-only).
	@output=$$(GOWORK=off go -C tools run golang.org/x/tools/cmd/goimports -l ..); \
	if [ -n "$$output" ]; then \
		echo "goimports would modify:"; \
		echo "$$output"; \
		exit 1; \
	fi

lint: fmt imports vet golangci-lint ## Run fmt, imports, vet, and golangci-lint.

golangci-lint: ## Run golangci-lint.
	@mkdir -p bin
	GOWORK=off go -C tools build -o $(CURDIR)/bin/golangci-lint github.com/golangci/golangci-lint/v2/cmd/golangci-lint
	$(CURDIR)/bin/golangci-lint run ./...

generate: ## Generate mocks with mockery.
	go run github.com/vektra/mockery/v2@v2.53.6 --config .mockery.yaml

tidy: ## Tidy and verify Go module dependencies.
	go mod tidy

tidy-check: ## Tidy modules and verify no changes.
	@go mod tidy
	@go -C tools mod tidy
	@git diff --exit-code go.mod go.sum tools/go.mod tools/go.sum

clean: ## Remove build artifacts.
	rm -f $(BINARY_NAME)

image-build: ## Build the container image.
	$(CONTAINER_RUNTIME) build -t $(IMAGE_NAME) .

image-run: ## Run the container image with stdio transport.
	$(CONTAINER_RUNTIME) run --rm -i -v $(RUNTIME_SOCKET):$(RUNTIME_SOCKET) $(IMAGE_NAME)

image-run-http: ## Run the container image with HTTP transport.
	$(CONTAINER_RUNTIME) run --rm -p 8080:8080 -v $(RUNTIME_SOCKET):$(RUNTIME_SOCKET) \
		-e CONTAINER_RUNTIME_MCP_TRANSPORT=http \
		-e CONTAINER_RUNTIME_MCP_HTTP_ADDR=0.0.0.0:8080 \
		$(IMAGE_NAME)

image-run-http-read-only: ## Run the container image with HTTP transport in read-only mode.
	$(CONTAINER_RUNTIME) run --rm -p 8080:8080 -v $(RUNTIME_SOCKET):$(RUNTIME_SOCKET) \
		-e CONTAINER_RUNTIME_MCP_TRANSPORT=http \
		-e CONTAINER_RUNTIME_MCP_HTTP_ADDR=0.0.0.0:8080 \
		-e CONTAINER_RUNTIME_MCP_READ_ONLY=true \
		$(IMAGE_NAME)
