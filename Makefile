.PHONY: build run test clean docker reset-db lint install-deps

# Variables
BINARY_NAME=vulcan
BUILD_DIR=./bin
DATA_DIR=./data
GO=go
GOFLAGS=-v

# Build the project
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/vulcan
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the node
run: build
	@echo "Starting Vulcan node..."
	$(BUILD_DIR)/$(BINARY_NAME) --api-port=8080 --port=6000 --db-path=$(DATA_DIR)/node1 -difficulty=20

# Run tests
test:
	@echo "Running tests..."
	$(GO) test ./... -v -cover -coverprofile=coverage.out
	@echo "Tests complete. Coverage report: coverage.out"

# Run tests with coverage report
test-coverage: test
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -rf $(DATA_DIR)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

# Reset database
reset-db:
	@echo "Resetting database..."
	rm -rf $(DATA_DIR)
	@echo "Database reset complete"

# Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t vulcan:latest .
	docker build -f Dockerfile.frontend -t vulcan-web:latest ./vulcan-web
	@echo "Docker images built"

# Run with Docker Compose
docker-up:
	docker-compose up --build

docker-down:
	docker-compose down -v

# Install dependencies
install-deps:
	@echo "Installing Go dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Installing frontend dependencies..."
	cd vulcan-web && npm install

# Run linters
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...
	cd vulcan-web && npm run lint

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	cd vulcan-web && npm run format || true

# Run multiple nodes for testing
run-multi:
	@echo "Starting multiple nodes..."
	@mkdir -p $(DATA_DIR)/node1 $(DATA_DIR)/node2 $(DATA_DIR)/node3
	@./scripts/start_nodes.sh

# Generate TLS certificates
gen-cert:
	@./scripts/gen_cert.sh

# Help
help:
	@echo "Vulcan Blockchain - Available commands:"
	@echo "  make build          - Build the binary"
	@echo "  make run            - Build and run a single node"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make reset-db       - Reset the database"
	@echo "  make docker         - Build Docker images"
	@echo "  make docker-up      - Run with Docker Compose"
	@echo "  make docker-down    - Stop Docker Compose"
	@echo "  make install-deps   - Install all dependencies"
	@echo "  make lint           - Run linters"
	@echo "  make fmt            - Format code"
	@echo "  make run-multi      - Run multiple interconnected nodes"
	@echo "  make gen-cert       - Generate self-signed TLS certificates"