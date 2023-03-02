include .env

default:help

MIGRATIONS_DIR := migrations
GOOSE := GOOSE_DRIVER=${GOOSE_DRIVER} GOOSE_DBSTRING=${GOOSE_DBSTRING} goose -dir $(MIGRATIONS_DIR)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: migrate
migrate: ## update database to latest migration
	@$(GOOSE) up

.PHONY: downgrade
downgrade:  ## downgrade db migrations
	@$(GOOSE) down

type ?= sql
.PHONY: new_migration
new_migration: ## create new migration
	@$(GOOSE) create ${name} ${type}

.PHONY: goose
goose: ## get migrations status
	@$(GOOSE) status

.PHONY: pre-commit
pre-commit: ## run linters and formatters via pre-commit
	@pre-commit run --all-files
