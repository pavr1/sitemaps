.DEFAULT_GOAL := help
.PHONY: help

start: ## starting sitemap http server 
	@docker compose up --build

stop: ## stopping sitemap http server
	@docker compose stop

lint: ## executes golang-ci-lint against all files
	golangci-lint run --timeout 10m0s ./...

test: # executes unit tests
	go test -v -cover ./...