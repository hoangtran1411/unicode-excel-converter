# Makefile for VniConverter
# Why: Automates common development tasks like building, testing, and cleaning to ensure consistency across environments.

BINARY_NAME=VniConverter
BUILD_DIR=build/bin

.PHONY: all build clean test coverage lint

all: build

# Build the application using Wails
# Why: Standard build command for Wails applications to generate the binary.
build:
	wails build

# Build for Windows specifically (creates .exe)
# Why: Explicitly targets Windows amd64 since this is the user's primary OS.
build-windows:
	wails build -platform windows/amd64

# Run unit tests
# Why: Ensures all logic is correct before building.
test:
	go test ./... -v

# Run tests with coverage and open report
# Why: Visualizes code coverage to identify untested logic paths.
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"

# Clean build artifacts
# Why: Removes generated binaries and temporary files to ensure a fresh build.
clean:
	rm -rf build/
	rm -f coverage.out coverage.html
	rm -f *.log
	rm -f *_output_*.xlsx

# Format code
# Why: Enforces standard Go formatting.
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
# Why: Catches potential errors and style issues early.
lint:
	golangci-lint run

# Install dependencies
# Why: Ensures all necessary tools and libraries are available.
deps:
	go mod tidy
	go install github.com/wailsapp/wails/v2/cmd/wails@latest
