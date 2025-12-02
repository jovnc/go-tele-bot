.PHONY: build start run deploy deploy-setup docker-run

APP_NAME=go-tele-bot
REGION=asia-southeast1
BINARY_NAME=go-tele-bot
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

deploy:
	@echo "Deploying $(APP_NAME) to Cloud Run..."
	@chmod +x scripts/deploy.sh
	@scripts/deploy.sh

deploy-setup:
	@echo "Setting up deployment for $(APP_NAME)..."
	@chmod +x scripts/deploy-setup.sh
	@scripts/deploy-setup.sh

docker-run:
	@echo "Running $(APP_NAME) in Docker..."
	@chmod +x scripts/docker-run.sh
	@scripts/docker-run.sh
