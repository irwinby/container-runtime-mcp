.PHONY: all build run run-http test test-race fmt vet imports lint generate tidy clean image-build image-run image-run-http image-run-http-read-only help

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

fmt: ## Format all Go code with gofmt.
	go fmt ./...

vet: ## Run go vet on all packages.
	go vet ./...

imports: ## Organize imports with goimports.
	GOWORK=off go -C tools run golang.org/x/tools/cmd/goimports -w ..

lint: fmt imports vet ## Run fmt, imports, and vet.

generate: ## Generate mocks with mockery.
	go run github.com/vektra/mockery/v2@v2.53.6 --config .mockery.yaml

tidy: ## Tidy and verify Go module dependencies.
	go mod tidy

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
