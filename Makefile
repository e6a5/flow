.PHONY: build install clean test fmt dev release

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date +%Y-%m-%dT%H:%M:%S%z)

# Build flags
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Build the flow binary
build:
	go build $(LDFLAGS) -o flow ./...

# Install flow to $GOPATH/bin
install:
	go install $(LDFLAGS)

# Clean build artifacts
clean:
	rm -f flow
	rm -rf dist/

# Run tests
test:
	go test -v ./...
	go vet ./...
	@echo "Testing install script..."
	@./scripts/test-install.sh

# Format code
fmt:
	go fmt ./...

# Development: run with example
dev:
	go run $(LDFLAGS) ./... 1 --tag "dev test"

# Build for multiple platforms
release:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/flow-linux-amd64 ./...
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/flow-linux-arm64 ./...
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/flow-darwin-amd64 ./...
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/flow-darwin-arm64 ./...
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/flow-windows-amd64.exe ./...
	@echo "Built binaries for version $(VERSION)" 