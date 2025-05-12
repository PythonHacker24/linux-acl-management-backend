APP_NAME = laclm
CMD_DIR = ./cmd/$(APP_NAME)
BIN_DIR = ./bin
BIN_PATH = $(BIN_DIR)/$(APP_NAME)

GOFILES := $(shell find . -name '*.go' -type f)

.PHONY: all build clean run test lint build-linux build-mac build-win

all: build

## Build the app
build: $(GOFILES)
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_PATH) $(CMD_DIR)

## Run the app
run: build
	@echo "Running $(APP_NAME)..."
	@$(BIN_PATH)

## Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)

## Run tests
test:
	@echo "Running tests..."
	go test ./...

## Lint (requires golangci-lint)
lint:
	@echo "Linting..."
	@golangci-lint run

## Cross-build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-linux $(CMD_DIR)
