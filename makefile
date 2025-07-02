# Project metadata
PROJECT_NAME := ./
VERSION := 0.1.0
BUILD_TIME := $(shell date +%FT%T%z)
GIT_COMMIT := $(shell git rev-parse --short HEAD)

# Go parameters
GO_CMD := go
GO_BUILD := $(GO_CMD) build
GO_TEST := $(GO_CMD) test
GO_CLEAN := $(GO_CMD) clean
GO_MOD := $(GO_CMD) mod
GO_RUN := $(GO_CMD) run
GO_LINT := golangci-lint

# Binary names
BINARY_NAME := $(PROJECT_NAME)
BINARY_UNIX := $(BINARY_NAME)_unix

# Directories
SRC_DIR := ./cmd/$(PROJECT_NAME)
DIST_DIR := ./dist
COVERAGE_DIR := ./coverage

# Flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"
TEST_FLAGS := -cover -coverprofile=$(COVERAGE_DIR)/coverage.out
LINT_FLAGS := run --timeout 5m

.PHONY: all build clean test lint run mod-tidy help

all: build

## build: Build the project binary
build:
	@echo "Building $(BINARY_NAME)..."
	@$(GO_BUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) $(SRC_DIR)

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@$(GO_CLEAN)
	@rm -rf $(DIST_DIR)
	@rm -rf $(COVERAGE_DIR)

## test: Run tests with coverage
test:
	@mkdir -p $(COVERAGE_DIR)
	@echo "Running tests..."
	@$(GO_TEST) $(TEST_FLAGS) ./...
	@$(GO_CMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

## lint: Run linter
lint:
	@echo "Running linter..."
	@$(GO_LINT) $(LINT_FLAGS) ./...

## run: Build and run the project
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(DIST_DIR)/$(BINARY_NAME)

## mod-tidy: Tidy Go modules
mod-tidy:
	@echo "Tidying modules..."
	@$(GO_MOD) tidy

## help: Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@awk '/^## / { \
		if ($$2 == "help:") { \
			next; \
		} \
		printf "  %-15s %s\n", $$2, substr($$0, index($$0, $$3)); \
	}' $(MAKEFILE_LIST) | column -t -s:
