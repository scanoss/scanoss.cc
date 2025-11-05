VERSION=$(shell git tag --sort=-version:refname | head -n 1)
APP_NAME = scanoss-cc
BUILD_DIR = build
DIST_DIR = dist
SCRIPTS_DIR = scripts
FRONTEND_DIR = frontend
ASSETS_DIR = assets
APP_BUNDLE = "$(BUILD_DIR)/bin/$(APP_NAME).app"

export GOTOOLCHAIN=go1.23.0

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: help

.DEFAULT_GOAL := help

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

clean:  ## Clean all build data
	@echo "Removing build data..."
	@rm -rf $(FRONTEND_DIR)/dist $(BUILD_DIR) $(DIST_DIR)

clean_all: clean  ## Clean all build data including Node
	@echo "Removing build & NPM data..."
	@rm -rf $(FRONTEND_DIR)/node_modules

clean_testcache:  ## Expire all Go test caches
	@echo "Cleaning test caches..."
	go clean -testcache ./backend/...

unit_test:  ## Run all unit tests in the backend folder
	@echo "Running unit test framework..."
	go test -v ./... -tags=unit

integration_test:  ## Run all integration tests
	@echo "Running integration tests..."
	SCANOSS_API_KEY=$(SC_API_KEY) go test -v ./... -tags=integration

lint: ## Run local instance of Go linting across the code base
	golangci-lint run ./...

lint-fix: ## Run local instance of Go linting across the code base including auto-fixing
	golangci-lint run --fix ./...

run: cp_assets ## Runs the application in development mode
	$(eval APPARGS := $(ARGS))
	@wails dev -ldflags "-X github.com/scanoss/scanoss.cc/backend/entities.AppVersion=$(VERSION)" $(if $(strip $(APPARGS)),-appargs "--debug $(APPARGS)")

npm: ## Install NPM dependencies for the frontend
	@echo "Running npm install for frontend..."
	cd frontend && npm install

cp_assets: ## Copy the necessary assets to the build folder
	@echo "Copying assets to build directory..."
	@mkdir -p $(BUILD_DIR)/assets
	@cp $(ASSETS_DIR)/* $(BUILD_DIR)/assets
	@cp $(ASSETS_DIR)/appicon.png $(BUILD_DIR)

build: clean cp_assets  ## Build the application image for the current platform
	@echo "Building application image..."
	@wails build -ldflags "-X github.com/scanoss/scanoss.cc/backend/entities.AppVersion=$(VERSION)"

binary: cp_assets  ## Build application binary only (no package)
	@echo "Build application binary only..."
	@wails build -ldflags "-X github.com/scanoss/scanoss.cc/backend/entities.AppVersion=$(VERSION)" --nopackage

build_macos: clean cp_assets  ## Build the application image for macOS
	@echo "Building application image for macOS..."
	@wails build -ldflags "-X github.com/scanoss/scanoss.cc/backend/entities.AppVersion=$(VERSION)" -platform darwin/universal -o "$(APP_NAME)"
	@echo "Build completed. Result: $(APP_BUNDLE)"

build_webkit41: clean cp_assets  ## Build the application image for Ubuntu 24.04+/Debian 13+ (webkit 4.1)
	@echo "Building application image with webkit2_41 tags..."
	@wails build -ldflags "-X github.com/scanoss/scanoss.cc/backend/entities.AppVersion=$(VERSION)" -tags webkit2_41
