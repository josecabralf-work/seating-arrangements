SHELL := /bin/bash

BINARY := bin/arrange
PKG    := ./cmd/arrange

.DEFAULT_GOAL := help

.PHONY: help deps build install run test lint format clean

help: ## Show available targets
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

deps: ## Download module dependencies
	go mod download

build: ## Build the binary into bin/
	@mkdir -p bin
	@go build -o $(BINARY) $(PKG)
	@echo "Built $(BINARY)"

install: ## Install the binary onto the PATH
	go install ./...

run: build ## Run the binary, loading variables from .env
	@set -a && source .env && set +a && ./$(BINARY)

test: ## Run the unit tests
	go test -count=1 -cover ./... $(ARGS)

lint: ## Check formatting and vet the packages
	@test -z "$$(gofmt -l .)" || { echo "Some files fail gofmt checks:"; gofmt -l .; exit 1; }
	go vet ./...

format: ## Rewrite the packages with gofmt
	go fmt ./...

clean: ## Remove ignored build artefacts
	git clean -fdX
