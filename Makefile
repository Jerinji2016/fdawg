# Variables
BINARY_NAME=fdawg
BUILD_DIR=build
GO_FILES=$(shell find . -name "*.go" -type f)
LDFLAGS=-ldflags "-s -w"
SCSS_DIR=internal/server/web/static/scss
CSS_DIR=internal/server/web/static/css

# Default target
.PHONY: all
all: clean sass build

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

# Install the binary to $GOPATH/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

# Build for multiple platforms
.PHONY: cross-build
cross-build: clean
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

