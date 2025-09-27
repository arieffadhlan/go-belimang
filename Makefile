# Load env
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

GOOSE = goose
DbUrl = postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_DIR = migrations

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f app

# Non Docker
.PHONY: migrate-push migrate-down migrate-make migrate-status

migrate-push:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DbUrl)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DbUrl)" down

migrate-make:
	@if [ -z "$(name)" ]; then \
		echo "Please provide a migration name. Example: make migrate-new-docker name=create_users_table"; \
	else \
		$(GOOSE) -dir $(MIGRATIONS_DIR) create create_$(name)_table sql; \
	fi

migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DbUrl)" status

# Docker
.PHONY: migrate-oush-docker migrate-down-docker migrate-make-docker migrate-status-docker

migrate-push-docker:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DBURL)" push

migrate-down-docker:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DBURL)" down

migrate-make-docker: ## Create new migration file (shared with docker, same as local)
	@if [ -z "$(name)" ]; then \
		echo "Please provide a migration name. Example: make migrate-new-docker name=create_users_table"; \
	else \
		$(GOOSE) -dir $(MIGRATIONS_DIR) create_$(name)_table sql; \
	fi

migrate-status-docker:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DBURL)" status

# App
dev:
	air

run:
	go run ./cmd/main.go

run-docker:
	docker-compose up app