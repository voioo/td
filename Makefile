VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Build targets
.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o td ./cmd/td

.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)" ./cmd/td

# Development targets
.PHONY: dev
dev: clean build test

.PHONY: clean
clean:
	go clean
	rm -f td coverage.out coverage.html

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: fmt vet

# Testing targets
.PHONY: test
test:
	go test ./...

.PHONY: test-verbose
test-verbose:
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-integration
test-integration:
	go test ./integration

.PHONY: test-unit
test-unit:
	go test ./internal/...

# CI/CD targets
.PHONY: ci
ci: lint test build

.PHONY: release
release: clean ci
	@echo "Building release version $(VERSION)"
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o td-linux-amd64 ./cmd/td
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o td-darwin-amd64 ./cmd/td
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o td-darwin-arm64 ./cmd/td
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o td-windows-amd64.exe ./cmd/td
	GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o td-windows-arm64.exe ./cmd/td

# Development helpers
.PHONY: run
run: build
	./td

.PHONY: run-debug
run-debug: build
	RUST_LOG=debug ./td

.PHONY: deps
deps:
	go mod tidy
	go mod verify

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  install        - Install the application"
	@echo "  dev            - Clean, build and test"
	@echo "  clean          - Clean build artifacts"
	@echo "  fmt            - Format Go code"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Format and vet code"
	@echo "  test           - Run all tests"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-unit      - Run unit tests only"
	@echo "  ci             - Run CI pipeline (lint + test + build)"
	@echo "  release        - Build release binaries for multiple platforms"
	@echo "  run            - Build and run the application"
	@echo "  run-debug      - Build and run with debug logging"
	@echo "  deps           - Tidy and verify dependencies"
	@echo "  help           - Show this help message"