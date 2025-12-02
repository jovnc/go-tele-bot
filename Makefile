.PHONY: build start run

APP_NAME=go_api
BINARY_NAME=go_api
BUILD_DIR=./build

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

start:build
	@echo "Starting $(APP_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

run:
	@echo "Running $(APP_NAME)..."
	@go run main.go
