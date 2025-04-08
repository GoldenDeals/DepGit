.PHONY: clear live lint db-migrate build build-debug build-release test test-verbose test-coverage test-html test-git gen-api web-install web-dev web-build

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty || echo "unknown")
BUILD_TIME ?= $(shell date +%Y-%m-%dT%H:%M:%S%z)
LDFLAGS := -ldflags "-s -w -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)'"
GO_BUILD = go build
BUILD_DIR = ./build
SRC_DIR = ./cmd/app
BIN_NAME = depgit
COVERAGE_FILE = coverage.txt
COVERAGE_HTML = coverage.html

all: build

app: build

build: build-release

build-debug:
	@echo "Building debug binary..."
	@mkdir -p $(BUILD_DIR)
	@$(GO_BUILD) -o $(BUILD_DIR)/$(BIN_NAME) $(SRC_DIR)
	@echo "Debug build complete: $(BUILD_DIR)/$(BIN_NAME)"

build-release:
	@echo "Building optimized release binary..."
	@mkdir -p $(BUILD_DIR)
	@$(GO_BUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BIN_NAME) $(SRC_DIR)
	@echo "Release build complete: $(BUILD_DIR)/$(BIN_NAME)"

clear:
	@rm -fr ./build

live:
	@air

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...

db-migrate:
	@echo "Applying database migrations..."
	@sqlite3 ./depgit.db < migrations/main.sql

test:
	@echo "Running tests..."
	@go test -race ./...

test-verbose:
	@echo "Running verbose tests..."
	@go test -v -race ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@go tool cover -func=$(COVERAGE_FILE)

test-html: test-coverage
	@echo "Generating HTML coverage report..."
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated at $(COVERAGE_HTML)"

test-git:
	@echo "Running git test script..."
	@scripts/test-git.sh

gen-api:
	@echo "Generating API code from OpenAPI spec..."
	@mkdir -p internal/gen/api
	@$(HOME)/go/bin/oapi-codegen -package api -generate types,server,spec api/openapi.yaml > internal/gen/api/api.go

web-install:
	@echo "Installing web dependencies..."
	cd web && npm install --legacy-peer-deps

web-dev:
	@echo "Starting web development server..."
	cd web && npm run dev

web-build:
	@echo "Building web application..."
	cd web && npm run build

# Full build including web application
build-full: build-release web-build
	@echo "Full build complete"