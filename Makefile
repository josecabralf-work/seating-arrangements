SHELL := /bin/bash

BINARY := bin/arrange
PKG    := ./cmd/arrange

.DEFAULT_GOAL := help

.PHONY: help build run

help: ## Show available targets
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary into bin/
	@mkdir -p bin
	@go build -o $(BINARY) $(PKG)
	@echo "Built $(BINARY)"

run: build ## Run the binary, loading variables from .env
	@set -a && source .env && set +a && ./$(BINARY)
