# Launchpad Project Makefile
# Provides convenient targets for building and running the entire project

.PHONY: help backend frontend smart-contracts all clean dev up down

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

backend: ## Build the Go backend
	cd backend && make server

frontend: ## Install frontend dependencies and build
	cd frontend && npm install && npm run build

smart-contracts: ## Install smart contract dependencies and compile
	cd smart-contracts && npm install && npm run compile

all: smart-contracts backend frontend ## Build all components

clean: ## Clean all build artifacts
	cd backend && make clean
	cd frontend && rm -rf .next dist
	cd smart-contracts && rm -rf artifacts cache

dev: ## Start development environment with Docker Compose
	cd deploy && docker compose up -d

up: dev ## Alias for dev

down: ## Stop development environment
	cd deploy && docker compose down

logs: ## Show logs from all services
	cd deploy && docker compose logs -f

.DEFAULT_GOAL := help