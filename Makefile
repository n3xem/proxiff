.PHONY: build
build:
	@echo "Building proxiff..."
	@go build -v -o proxiff ./cmd/proxiff
	@echo "Build completed"

.PHONY: build-examples
build-examples:
	@echo "Building example binaries..."
	@go build -v -o sample-server ./example/servers
	@echo "Example binaries built successfully"

.PHONY: build-all
build-all: build build-examples
	@echo "All builds completed"

.PHONY: test
test:
	@go test ./... -v -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: hadolint
hadolint:
	@echo "Running hadolint on Dockerfile..."
	@docker run --rm -i hadolint/hadolint < Dockerfile
	@echo "Hadolint check completed"

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f coverage.txt
	@rm -f proxiff
	@rm -f sample-server
	@echo "Clean completed"
