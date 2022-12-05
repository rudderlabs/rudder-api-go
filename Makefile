.PHONY: help test test-it

help: ## Show the available commands
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' ./Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Run all unit tests
	go test ./...

test-it: ## Run all test, including integration tests
	go test -tags integrationtest ./...