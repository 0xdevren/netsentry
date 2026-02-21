BINARY_NAME    := netsentry
BINARY_PATH    := bin/$(BINARY_NAME)
MODULE         := github.com/0xdevren/netsentry
CMD_PATH       := ./cmd/netsentry
VERSION        := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT         := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE     := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS        := -s -w \
                  -X $(MODULE)/internal/app.Version=$(VERSION) \
                  -X $(MODULE)/internal/app.Commit=$(COMMIT) \
                  -X $(MODULE)/internal/app.BuildDate=$(BUILD_DATE)
GOFLAGS        := -trimpath
TEST_TIMEOUT   := 120s
COVERAGE_OUT   := coverage.out
COVERAGE_HTML  := coverage.html

.PHONY: all build clean test test-race test-coverage lint fmt vet tidy install release docker help

all: clean build test

build:
	@echo "==> Building $(BINARY_NAME) $(VERSION)"
	@mkdir -p bin
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_PATH) $(CMD_PATH)
	@echo "==> Binary: $(BINARY_PATH)"

build-linux:
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_PATH)-linux-amd64 $(CMD_PATH)

build-darwin:
	@mkdir -p bin
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_PATH)-darwin-arm64 $(CMD_PATH)

build-windows:
	@mkdir -p bin
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_PATH)-windows-amd64.exe $(CMD_PATH)

install:
	go install $(GOFLAGS) -ldflags "$(LDFLAGS)" $(CMD_PATH)

clean:
	@rm -rf bin/ $(COVERAGE_OUT) $(COVERAGE_HTML)

test:
	go test -timeout $(TEST_TIMEOUT) ./...

test-race:
	go test -race -timeout $(TEST_TIMEOUT) ./...

test-coverage:
	go test -coverprofile=$(COVERAGE_OUT) -covermode=atomic -timeout $(TEST_TIMEOUT) ./...
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	go tool cover -func=$(COVERAGE_OUT) | tail -1

lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found"; exit 1)
	golangci-lint run --timeout 5m ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

tidy:
	go mod tidy

docker:
	docker build -f deployments/docker/Dockerfile -t netsentry:$(VERSION) .

release: build-linux build-darwin build-windows

generate:
	go generate ./...

help:
	@grep -E '^[a-zA-Z_-]+:' Makefile | awk -F: '{print "  "$$1}'
