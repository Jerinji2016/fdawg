# Variables
BINARY_NAME=fdawg
BUILD_DIR=build
GO_FILES=$(shell find . -name "*.go" -type f)
LDFLAGS=-ldflags "-s -w"
SCSS_DIR=internal/server/web/static/scss
CSS_DIR=internal/server/web/static/css

# Default target
# Help target
.PHONY: help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Build commands:"
	@echo "  build          - Build the binary"
	@echo "  clean          - Clean build artifacts"
	@echo "  all-platforms  - Build binaries for all platforms"
	@echo ""
	@echo "Development commands:"
	@echo "  run            - Build and run the application"
	@echo "  dev            - Run in development mode with SASS watching"
	@echo "  test           - Run tests"
	@echo ""
	@echo "SASS commands:"
	@echo "  sass           - Compile SASS to CSS"
	@echo "  sass-watch     - Watch SASS files for changes"
	@echo ""
	@echo "Documentation commands:"
	@echo "  docs-install   - Install documentation dependencies"
	@echo "  docs-build     - Build documentation"
	@echo "  docs-serve     - Serve documentation at http://localhost:4000"
	@echo "  docs-clean     - Clean documentation build artifacts"
	@echo "  docs           - Build and serve documentation"
	@echo ""
	@echo "Other commands:"
	@echo "  all            - Clean, compile SASS, and build (default)"
	@echo "  help           - Show this help message"


# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/flutter-manager

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# Compile SASS to CSS
.PHONY: sass
sass:
	@echo "Compiling SASS to CSS..."
	@npm run sass

# Watch SASS files for changes
.PHONY: sass-watch
sass-watch:
	@echo "Watching SASS files for changes..."
	@npm run sass:watch

# Run the application
.PHONY: run
run: sass build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Run the application in development mode (with SASS watching)
.PHONY: dev
dev: build
	@echo "Running $(BINARY_NAME) in development mode..."
	@echo "SASS files will be watched for changes..."
	@npm run sass:watch & ./$(BUILD_DIR)/$(BINARY_NAME)

# Build for multiple platforms
.PHONY: all-platforms
all-platforms: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)

	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/flutter-manager

	@echo "Building for macOS..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/flutter-manager
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/flutter-manager

	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/flutter-manager

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Documentation commands
.PHONY: docs-install
docs-install:
	@echo "Installing documentation dependencies..."
	@cd docs && bundle install

.PHONY: docs-build
docs-build:
	@echo "Building documentation..."
	@cd docs && bundle exec jekyll build

.PHONY: docs-serve
docs-serve:
	@echo "Serving documentation at http://localhost:4000..."
	@cd docs && bundle exec jekyll serve --incremental

.PHONY: docs-clean
docs-clean:
	@echo "Cleaning documentation build artifacts..."
	@cd docs && bundle exec jekyll clean

.PHONY: docs
docs: docs-serve
