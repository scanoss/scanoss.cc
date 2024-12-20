VERSION=$(shell git tag --sort=-version:refname | head -n 1)
APP_NAME = scanoss-lui
BUILD_DIR = build
DIST_DIR = dist
SCRIPTS_DIR = scripts
FRONTEND_DIR = frontend
PACKAGE_ROOT = package_root
APP_BUNDLE = $(BUILD_DIR)/bin/$(APP_NAME).app
PKG_NAME = $(APP_NAME)-macos-$(VERSION).pkg


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

go_lint_local: ## Run local instance of Go linting across the code base
	golangci-lint run ./...

go_lint_local_fix: ## Run local instance of Go linting across the code base including auto-fixing
	golangci-lint run --fix ./...

go_lint_docker: ## Run docker instance of Go linting across the code base
	docker run --rm -v $(pwd):/app -v ~/.cache/golangci-lint/v1.50.1:/root/.cache -w /app golangci/golangci-lint:v1.50.1 golangci-lint run ./backend/...

run: ## Runs the application in development mode
	$(eval APPARGS := $(ARGS))
	@wails dev -ldflags "-X github.com/scanoss/scanoss.lui/backend/entities.AppVersion=$(VERSION)" $(if $(strip $(APPARGS)),-appargs "$(APPARGS)")

npm: ## Install NPM dependencies for the frontend
	@echo "Running npm install for frontend..."
	cd frontend && npm install

cp_assets: ## Copy the necessary assets to the build folder
	@echo "Copying assets to build..."
	@mkdir -p build
	@cp assets/appicon.png build/appicon.png
	@cp -r assets build/assets

build: cp_assets  ## Build the application image
	@echo "Building application image..."
	@wails build -ldflags "-X github.com/scanoss/scanoss.lui/backend/entities.AppVersion=$(VERSION)"

binary: cp_assets  ## Build application binary only (no package)
	@echo "Build application binary only..."
	@wails build -ldflags "-X github.com/scanoss/scanoss.lui/backend/entities.AppVersion=$(VERSION)" --nopackage

build_macos: clean cp_assets  ## Build the application image for macOS
	@echo "Building application image for macOS..."
	@wails build -ldflags "-X github.com/scanoss/scanoss.lui/backend/entities.AppVersion=$(VERSION)" -platform darwin/universal
	@echo "Build completed. Result: $(APP_BUNDLE)"

package_macos: build_macos ## Package the built macOS app into a pkg
	@echo "Packaging for macOS with .pkg..."

	@mkdir -p $(DIST_DIR)

	# Prepare a clean staging directory with only the .app
	@rm -rf $(PACKAGE_ROOT)
	@mkdir -p $(PACKAGE_ROOT)
	@cp -R $(APP_BUNDLE) $(PACKAGE_ROOT)/

	@chmod +x $(SCRIPTS_DIR)/macos_postinstall

	@pkgbuild \
		--root $(PACKAGE_ROOT) \
		--scripts $(SCRIPTS_DIR) \
		--identifier "com.scanoss.$(APP_NAME)" \
		--version "$(VERSION)" \
		--install-location "/Applications" \
		$(DIST_DIR)/$(PKG_NAME)

	@rm -rf $(PACKAGE_ROOT)

	@echo "$(DIST_DIR)/$(PKG_NAME) created. Run it to install."