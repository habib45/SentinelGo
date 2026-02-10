# SentinelGo Makefile
# Build for release with environment variable configuration

# Version can be set via:
# - git tag (automatically detected)
# - VERSION env var
# - defaults to "dev"
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Build flags for version injection
LDFLAGS=-ldflags "-X sentinelgo/cmd/sentinelgo.Version=$(VERSION) -X sentinelgo/internal/config.Version=$(VERSION)"

# Targets
.PHONY: build clean all windows linux macos release version

all: windows linux macos

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o build/windows/sentinelgo-windows-amd64.exe ./cmd/sentinelgo

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/linux/sentinelgo-linux-amd64 ./cmd/sentinelgo
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o build/linux/sentinelgo-linux-arm64 ./cmd/sentinelgo

macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o build/darwin/sentinelgo-darwin-amd64 ./cmd/sentinelgo
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o build/darwin/sentinelgo-darwin-arm64 ./cmd/sentinelgo

# Build all platforms for release
release: clean all
	@echo "Release built with version $(VERSION)"
	@echo "Assets created in build/ directory:"
	@find build -type f -name "*sentinelgo*" -exec ls -lh {} \;

clean:
	rm -rf build/

# Development build (current platform only)
build:
	go build $(LDFLAGS) -o bin/sentinelgo ./cmd/sentinelgo

# Show version information
version:
	@echo "Current version: $(VERSION)"
	@echo "Git tag: $(shell git describe --tags --always 2>/dev/null || echo 'none')"
	@echo "Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"

# Test version injection
test-version: build
	@echo "Testing version injection..."
	@./bin/sentinelgo -version

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run tests
test:
	go test ./...

# Create release assets directory structure
setup:
	mkdir -p build/windows build/linux build/darwin

# Example usage:
# make release VERSION=v1.0.0
# make release (uses git tag or "dev")
# make build (development build for current platform)
# make version (show current version)
# make test-version (test version injection)
