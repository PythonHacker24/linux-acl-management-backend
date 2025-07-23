APP_NAME := laclm
CMD_DIR := ./cmd/$(APP_NAME)
BIN_DIR := ./bin
BUILD_DIR := ./build

GOFILES := $(shell find . -name '*.go' -type f)

# Target platforms: OS_ARCH
TARGETS := \
	linux_amd64 \
	linux_arm64

.PHONY: all build build-cross clean run test lint vendor package

## Default target
all: build

## Build for local OS/arch using vendored deps
build: vendor $(GOFILES)
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	GOOS="" GOARCH="" go build -mod=vendor -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

## Build cross-compiled binaries for all Linux targets
build-cross: vendor $(GOFILES)
	@echo "Cross building for targets: $(TARGETS)"
	@mkdir -p $(BIN_DIR)
	@for target in $(TARGETS); do \
		OS=$${target%_*}; \
		ARCH=$${target#*_}; \
		OUT=$(BIN_DIR)/$(APP_NAME)-$$OS-$$ARCH; \
		echo "Building $$OUT..."; \
		GOOS=$$OS GOARCH=$$ARCH go build -mod=vendor -o $$OUT $(CMD_DIR); \
	done

## Run the app
run: build
	@echo "Running $(APP_NAME)..."
	@$(BIN_DIR)/$(APP_NAME)

## Clean build and package directories
clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR) $(BUILD_DIR) vendor

## Run tests
test:
	@echo "Running tests..."
	go test ./...

## Lint (requires golangci-lint)
lint:
	@echo "Linting..."
	@golangci-lint run

## Vendor dependencies
vendor:
	@echo "Vendoring dependencies..."
	go mod vendor

## Package full project source (with vendor) for each target
package: clean vendor
	@echo "Packaging full source tarballs for: $(TARGETS)"
	@mkdir -p $(BUILD_DIR)
	@for target in $(TARGETS); do \
		OS=$${target%_*}; \
		ARCH=$${target#*_}; \
		NAME=$(APP_NAME)-$$OS-$$ARCH; \
		TARBALL=$$NAME-source.tar.gz; \
		echo "Creating $$TARBALL..."; \
		mkdir -p tmp/$$NAME; \
		cp -r * tmp/$$NAME; \
		rm -rf tmp/$$NAME/$(BUILD_DIR) tmp/$$NAME/$(BIN_DIR); \
		tar --no-xattrs -czf $(BUILD_DIR)/$$TARBALL -C tmp $$NAME; \
		rm -rf tmp/$$NAME; \
	done
	@rm -rf tmp
