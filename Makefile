# Makefile for Canopy - Fork of canopy-network/canopy
# Provides common development and deployment commands

.PHONY: all build run stop clean test lint fmt docker-build docker-up docker-down help

# Go binary name
BINARY_NAME=canopy
# Go module name (update if different)
MODULE=$(shell go list -m 2>/dev/null || echo "github.com/your-org/canopy")
# Build output directory
BUILD_DIR=./build
# Docker compose file
COMPOSE_FILE=docker-compose.yml

# Default target
all: build

## build: Compile the Go binary
build:
	@echo ">> Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./...

## run: Run the application locally
run:
	@echo ">> Running $(BINARY_NAME)..."
	go run ./...

## test: Run all unit tests
test:
	@echo ">> Running tests..."
	go test -v -race -cover ./...

## test-short: Run tests without the race detector (faster for quick checks)
test-short:
	@echo ">> Running tests (short)..."
	go test -short -cover ./...

## lint: Run golangci-lint
lint:
	@echo ">> Linting..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found, install via: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh" && exit 1)
	golangci-lint run ./...

## fmt: Format Go source files
fmt:
	@echo ">> Formatting..."
	gofmt -s -w .
	goimports -w . 2>/dev/null || true

## tidy: Tidy Go module dependencies
tidy:
	@echo ">> Tidying modules..."
	go mod tidy

## docker-build: Build the Docker image
docker-build:
	@echo ">> Building Docker image..."
	docker build -f .docker/Dockerfile -t $(BINARY_NAME):latest .

## docker-up: Start all services via Docker Compose
docker-up:
	@echo ">> Starting services..."
	docker compose -f $(COMPOSE_FILE) up -d

## docker-down: Stop all services via Docker Compose
docker-down:
	@echo ">> Stopping services..."
	docker compose -f $(COMPOSE_FILE) down

## docker-logs: Tail logs from Docker Compose services
docker-logs:
	docker compose -f $(COMPOSE_FILE) logs -f

## clean: Remove build artifacts
clean:
	@echo ">> Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out

## dev: fmt + tidy + test in one step (handy for quick iteration)
# NOTE: switched from test-short to test so the race detector always runs locally;
# the extra time is worth catching races early on my machine.
# NOTE: removed lint from dev target - I find it too slow for rapid iteration;
# run 'make lint' manually before opening a PR instead.
dev: fmt tidy test

## ci: lint + test - intended to mirror what CI runs, useful to check locally before pushing
# NOTE: this is my personal sanity-check target; run 'make ci' before pushing a branch
# to catch anything that would fail in the pipeline without waiting for remote CI.
ci: lint test

## cover: Run tests and open an HTML coverage report in the browser
# NOTE: handy for visually spotting untested code paths while exploring the codebase.
cover:
	@echo ">> Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
