VERSION=$(shell git tag --sort=-version:refname | head -n 1)

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: help

.DEFAULT_GOAL := help

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

clean:  ## Clean all build data
	@echo "Removing build data..."
	@rm -rf build frontend/dist

clean_all: clean  ## Clean all build data including Node
	@echo "Removing build & NPM data..."
	@rm -rf frontend/node_modules

clean_testcache:  ## Expire all Go test caches
	@echo "Cleaning test caches..."
	go clean -testcache ./backend/...

unit_test:  ## Run all unit tests in the backend folder
	@echo "Running unit test framework..."
	go test -v ./backend/...

go_lint_local: ## Run local instance of Go linting across the code base
	golangci-lint run ./backend/...

go_lint_local_fix: ## Run local instance of Go linting across the code base including auto-fixing
	golangci-lint run --fix ./backend/...

go_lint_docker: ## Run docker instance of Go linting across the code base
	docker run --rm -v $(pwd):/app -v ~/.cache/golangci-lint/v1.50.1:/root/.cache -w /app golangci/golangci-lint:v1.50.1 golangci-lint run ./backend/...

run: ## Runs the application in development mode
	$(eval APPARGS := $(ARGS))
	@wails dev -ldflags "-X github.com/scanoss/scanoss.lui/backend/main/pkg/common/version.AppVersion=$(VERSION)" $(if $(strip $(APPARGS)),-appargs "$(APPARGS)")

npm: ## Install NPM dependencies for the frontend
	@echo "Running npm install for frontend..."
	cd frontend && npm install

cp_assets: ## Copy the necessary assets to the build folder
	@echo "Copying assets to build..."
	@mkdir -p build
	@cp -r assets build/assets

build: cp_assets  ## Build the application image
	@echo "Building application image..."
	@wails build -ldflags "-X github.com/scanoss/scanoss.lui/backend/main/pkg/common/version.AppVersion=$(VERSION)"

binary: cp_assets  ## Build application binary only (no package)
	@echo "Build application binary only..."
	@wails build -ldflags "-X github.com/scanoss/scanoss.lui/backend/main/pkg/common/version.AppVersion=$(VERSION)" --nopackage
