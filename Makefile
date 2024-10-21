VERSION=$(shell git tag --sort=-version:refname | head -n 1)

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: help

.DEFAULT_GOAL := help

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Runs the application in development mode
	$(eval APPARGS := $(ARGS))
	@wails dev -ldflags "-X github.com/scanoss/scanoss.lui/backend/main/pkg/common/version.AppVersion=$(VERSION)" $(if $(strip $(APPARGS)),-appargs "$(APPARGS)")
