MODULE_NAME=openkms
BUILD_TIME=$(shell date "+%Y-%m-%d %H:%M:%S")

BIN_DIR=bin
GO_MOD_FILE=go.mod
GO_EXTERN_PACKAGE=

all: build

check_mod:
	@if [ ! -f "$(GO_MOD_FILE)" ]; then \
		echo "go.mod not found. Initializing Go module..."; \
		go mod init $(MODULE_NAME); \
	fi

build: check_mod
	@echo "Building the Go application..."
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-X 'main.BuildDate=$(BUILD_TIME)'" -o $(BIN_DIR)/$(MODULE_NAME)
	@echo "Build complete. Binary located at $(BIN_DIR)/$(MODULE_NAME)"

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
	@echo "Clean complete."

deps:
	@echo "Installing dependencies..."
	go mod tidy
	@if [ ! -z "$(GO_EXTERN_PACKAGE)" ]; then \
		echo "Installing Go packages: $(GO_EXTERN_PACKAGE)..."; \
		go get $(GO_EXTERN_PACKAGE); \
	fi

.PHONY: all build clean deps
