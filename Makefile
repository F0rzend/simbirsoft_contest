include .env

default:help

##@ General

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: run
run: ## Run application in docker
	docker compose up --build

##@ Database

MIGRATIONS_DIR := migrations
GOOSE := GOOSE_DRIVER=${GOOSE_DRIVER} GOOSE_DBSTRING=${GOOSE_DBSTRING} goose -dir $(MIGRATIONS_DIR)

.PHONY: migrate
migrate: ## Update database to latest migration
	@$(GOOSE) up

.PHONY: downgrade
downgrade: ## Downgrade db migrations
	@$(GOOSE) down

type ?= sql
.PHONY: new_migration
new_migration: ## Create new migration
	@$(GOOSE) create ${name} ${type}

.PHONY: goose
goose: ## Get migrations status
	@$(GOOSE) status

##@ Linters and formatters

.PHONY: pre-commit
pre-commit: ## Run linters and formatters via pre-commit
	@pre-commit run --all-files
