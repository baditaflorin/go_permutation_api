.PHONY: help build wasm docs-serve test run run-api run-gui clean docker-build docker-up docker-down deps lint fmt vet coverage install

# Variables
BINARY_NAME=permutation-api
MAIN_PATH=./cmd/server
BUILD_DIR=./bin
DOCKER_IMAGE=permutation-api
GO=go
GOFLAGS=-v

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "Available targets:"
	@echo ""
	@grep -E '^##' Makefile | sed 's/## /  /'
	@echo ""

## build: Build the application binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## wasm: Build the WebAssembly module for the GitHub Pages site
wasm:
	@echo "Building WASM module..."
	GOARCH=wasm GOOS=js $(GO) build -o docs/permutation.wasm ./cmd/wasm/
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" docs/wasm_exec.js 2>/dev/null || \
	  find $$(go env GOPATH)/pkg/mod/golang.org -name "wasm_exec.js" 2>/dev/null | sort -r | head -1 | xargs -I{} cp {} docs/wasm_exec.js
	@echo "WASM build complete: docs/permutation.wasm ($(shell ls -lh docs/permutation.wasm | awk '{print $$5}'))"

## docs-serve: Serve the documentation site locally on port 8000
docs-serve:
	@echo "Serving docs at http://localhost:8000"
	cd docs && python3 -m http.server 8000

## install: Install the application to $GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	CGO_ENABLED=1 $(GO) install $(GOFLAGS) $(MAIN_PATH)
	@echo "Installation complete"

## test: Run all tests
test:
	@echo "Running tests..."
	$(GO) test -v -race -coverprofile=coverage.out ./...

## coverage: Run tests with coverage report
coverage: test
	@echo "Generating coverage report..."
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## run: Run the application in CLI mode (args: ARGS="a b c")
run: build
	@echo "Running $(BINARY_NAME) in CLI mode..."
	$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

## run-api: Run the API server
run-api: build
	@echo "Starting API server..."
	$(BUILD_DIR)/$(BINARY_NAME) --serve

## run-gui: Run the GUI configuration server
run-gui: build
	@echo "Starting GUI server..."
	$(BUILD_DIR)/$(BINARY_NAME) --gui

## run-both: Run both API and GUI servers (requires separate terminals)
run-both:
	@echo "Starting both servers..."
	@echo "API: http://localhost:8080"
	@echo "GUI: http://localhost:3000"
	@$(MAKE) run-api & $(MAKE) run-gui

## deps: Download and verify dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod verify
	@echo "Dependencies updated"

## tidy: Tidy go.mod and go.sum
tidy:
	@echo "Tidying go modules..."
	$(GO) mod tidy
	@echo "Modules tidied"

## fmt: Format all Go source files
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "Formatting complete"

## vet: Run go vet on all packages
vet:
	@echo "Running go vet..."
	$(GO) vet ./...
	@echo "Vet complete"

## lint: Run golangci-lint (requires golangci-lint installed)
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run ./...

## clean: Remove build artifacts and temporary files
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f configs/runtime_config.json
	$(GO) clean
	@echo "Clean complete"

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

## docker-up: Start all services with docker-compose
docker-up:
	@echo "Starting services with docker-compose..."
	docker-compose up -d
	@echo "Services started"
	@echo "API: http://localhost:8080"
	@echo "GUI: http://localhost:3000"

## docker-down: Stop all services
docker-down:
	@echo "Stopping services..."
	docker-compose down
	@echo "Services stopped"

## docker-logs: View docker-compose logs
docker-logs:
	docker-compose logs -f

## docker-clean: Remove all containers, images, and volumes
docker-clean: docker-down
	@echo "Cleaning Docker artifacts..."
	docker-compose down -v --rmi all
	@echo "Docker cleanup complete"

## bench: Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

## check: Run all checks (fmt, vet, test)
check: fmt vet test
	@echo "All checks passed"

## example-csv: Run with example CSV file
example-csv: build
	@echo "name,age,city" > /tmp/example.csv
	@echo "Alice,30,NYC" >> /tmp/example.csv
	@echo "Bob,25,LA" >> /tmp/example.csv
	@echo "Charlie,35,SF" >> /tmp/example.csv
	@echo "CSV file created: /tmp/example.csv"
	$(BUILD_DIR)/$(BINARY_NAME) --csv=/tmp/example.csv --column=0

## example-tsv: Run with example TSV file
example-tsv: build
	@echo -e "name\tage\tcity" > /tmp/example.tsv
	@echo -e "Alice\t30\tNYC" >> /tmp/example.tsv
	@echo -e "Bob\t25\tLA" >> /tmp/example.tsv
	@echo "TSV file created: /tmp/example.tsv"
	$(BUILD_DIR)/$(BINARY_NAME) --tsv=/tmp/example.tsv --column=0

## db-setup: Setup database with sample data (requires PostgreSQL)
db-setup:
	@echo "Setting up database..."
	@which psql > /dev/null || (echo "PostgreSQL client not installed" && exit 1)
	psql -U postgres -h localhost -p 5432 -f init.sql
	@echo "Database setup complete"

## version: Display Go version
version:
	@$(GO) version

## info: Display project information
info:
	@echo "Project: Go Permutation API"
	@echo "Binary: $(BINARY_NAME)"
	@echo "Go Version: $(shell $(GO) version)"
	@echo "Build Dir: $(BUILD_DIR)"
	@echo "Main Path: $(MAIN_PATH)"
