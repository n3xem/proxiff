.PHONY: generate
generate:
	@echo "Generating protobuf code..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		plugin/proto/comparator.proto
	@echo "Protobuf code generation completed"

.PHONY: install-tools
install-tools:
	@echo "Installing protoc plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Tools installed successfully"

.PHONY: build
build:
	@echo "Building proxiff..."
	@go build -v -o proxiff ./cmd/proxiff
	@echo "Build completed"

.PHONY: build-examples
build-examples:
	@echo "Building example binaries..."
	@go build -v -o sample-server ./example/servers
	@go build -v -o plugin-status-only ./example/plugin-status-only
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
	@rm -f plugin-status-only
	@echo "Clean completed"
