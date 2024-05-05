SHELL := /bin/bash
include .env
export
export APP_NAME := $(basename $(notdir $(shell pwd)))

.PHONY: help
help: ## display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: up
up: ## docker compose up with air hot reload
	@docker compose --project-name ${APP_NAME} --file ./.docker/compose.yaml up -d

.PHONY: down
down: ## docker compose down
	@docker compose --project-name ${APP_NAME} down --volumes

.PHONY: psql
psql:
	@docker exec -it postgres psql -U postgres

.PHONY: fmt
fmt:
	@(cd schema && bun run prisma format)

.PHONY: mitrate
mitrate:
	@(cd schema && bun run prisma db push)
	@dbmate dump

.PHONY: gen
gen:
	@find pkg/sqlc -type f -not -name "*.sql" -exec rm -rf {} \;
	@sqlboiler psql
	@sqlc generate
	@go mod tidy
	@go mod vendor

.PHONY: run
run: ## run the application
	@go run main.go
