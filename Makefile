# gopher-client Makefile
# Simple Makefile for Go client project

# Variables
VERSION?=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
BINARY_NAME?=gopher-client
BUILD_DIR?=bin
COVERAGE_DIR?=coverage
TEST_ARGS?=./...

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION)"
BUILD_FLAGS=-v $(LDFLAGS)

# Test flags
TEST_FLAGS=-v -coverprofile=$(COVERAGE_DIR)/coverage.txt -covermode=atomic
RACE_FLAGS=-v -race -coverprofile=$(COVERAGE_DIR)/coverage.txt -covermode=atomic

.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Clean build artifacts and coverage files
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)
	@$(GOCLEAN)

.PHONY: deps
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	@$(GOMOD) download
	@$(GOMOD) tidy

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting Go code..."
	@$(GOFMT) -s -w .

.PHONY: build
build: deps ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

.PHONY: run
run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

.PHONY: test
test: deps ## Run tests with coverage (used by GitHub Actions)
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	@$(GOTEST) -v -count=1 -p=1 -coverprofile=$(COVERAGE_DIR)/coverage.txt -covermode=atomic $(TEST_ARGS)

.PHONY: test-ginkgo
test-ginkgo: deps ## Run tests using Ginkgo
	@echo "Running tests with Ginkgo..."
	@mkdir -p $(COVERAGE_DIR)
	@ginkgo -v --coverprofile=$(COVERAGE_DIR)/coverage.txt --covermode=atomic $(TEST_ARGS)

.PHONY: test-race
test-race: deps ## Run tests with race detection (requires CGO)
	@echo "Running tests with race detection..."
	@mkdir -p $(COVERAGE_DIR)
	@CGO_ENABLED=1 $(GOTEST) $(RACE_FLAGS) $(TEST_ARGS)

.PHONY: test-verbose
test-verbose: deps ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@mkdir -p $(COVERAGE_DIR)
	@$(GOTEST) $(TEST_FLAGS) $(TEST_ARGS)

.PHONY: test-agent
test-agent: deps ## Run only agent package tests (lightweight)
	@echo "Running agent package tests..."
	@mkdir -p $(COVERAGE_DIR)
	@$(GOTEST) -v -count=1 -p=1 -coverprofile=$(COVERAGE_DIR)/coverage.txt -covermode=atomic ./agent

.PHONY: coverage
coverage: test ## Generate and display coverage report
	@echo "Generating coverage report..."
	@$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.txt -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage.txt

.PHONY: install
install: build ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

.PHONY: uninstall
uninstall: ## Remove the binary from GOPATH/bin
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(GOPATH)/bin/$(BINARY_NAME)

# Default target
.DEFAULT_GOAL := help