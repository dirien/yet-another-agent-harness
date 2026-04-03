# Makefile for Agent-Ready Go projects
# Copy to project root and adjust BINARY_NAME and MAIN_PKG.

BINARY_NAME        ?= app
MAIN_PKG           ?= ./cmd/$(BINARY_NAME)
BUILD_DIR          ?= ./bin
COVERAGE_THRESHOLD ?= 80
GOFLAGS            ?= -trimpath
VERSION            ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME         ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS            ?= -s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

.PHONY: all build build-static test cover lint vet fmt generate vuln tidy ci-check-tidy bench ci clean help

all: build test

build: ## compile the binary
	@mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PKG)

build-static: ## compile static binary for containers (CGO_ENABLED=0)
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-static $(MAIN_PKG)

test: ## run tests with race detector
	go test -race -count=1 -timeout=5m ./...

cover: ## run tests with coverage report (enforces COVERAGE_THRESHOLD)
	go test -race -count=1 -timeout=5m -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@COVERAGE=$$(go tool cover -func=coverage.out | tail -1 | awk '{print $$3}' | tr -d '%'); \
	echo "Total coverage: $${COVERAGE}%"; \
	awk -v c="$${COVERAGE}" -v t="$(COVERAGE_THRESHOLD)" \
	    'BEGIN { if (c+0 < t+0) { print "Coverage " c "% is below threshold " t "%"; exit 1 } }'

lint: ## run golangci-lint
	golangci-lint run ./...

vet: ## run go vet
	go vet ./...

fmt: ## format source files with gofmt and goimports
	gofmt -w ./...
	goimports -w ./...

generate: ## run go generate
	go generate ./...

vuln: ## check for known vulnerabilities with govulncheck
	govulncheck ./...

tidy: ## tidy and verify module dependencies
	go mod tidy
	go mod verify

ci-check-tidy: ## verify go.mod/go.sum are tidy (for CI — fails if tidy would make changes)
	go mod tidy
	git diff --exit-code go.mod go.sum || (echo "go.mod/go.sum not tidy; run 'make tidy' and commit"; exit 1)

bench: ## run benchmarks
	go test -bench=. -benchmem -run='^$$' ./...

ci: ci-check-tidy vet lint cover vuln ## full pipeline: tidy-check, vet, lint, cover, vuln

clean: ## remove build artifacts and coverage files
	rm -rf $(BUILD_DIR) coverage.out coverage.html

help: ## show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
