# Enhanced Makefile V4 - Simplified IKV Management
# Usage: make <target> [args]

# Define variables
DOCKER_EXEC := docker-compose exec web
GO_CMD := go run main.go

# Display help
.PHONY: help
help: ## Show available commands
	$(info =================================)
	$(info IKV Simplified Commands)
	$(info =================================)
	$(info )
	$(info Quick Setup Commands:)
	$(info   setup                   Initialize project with Docker)
	$(info   dev                     Start development environment)
	$(info )
	$(info Database Migration Commands:)
	$(info   migrate/all/up          Run all migrations for all modules)
	$(info   migrate/all             Show pending migrations for all modules)
	$(info   migrate/{module}/up     Run migrations for specific module)
	$(info   migrate/{module}/down   Rollback migration for specific module)
	$(info   migrate/{module}/create Create new migration for specific module)
	$(info   migrate/{module}/status Show migration status for specific module)
	$(info   migrate/{module}        Show pending migrations for specific module)
	$(info )
	$(info Development Commands:)
	$(info   logs                    Show application logs)
	$(info   shell                   Access development container shell)
	$(info   restart                 Restart development services)
	$(info   clean                   Clean build artifacts and Docker)
	$(info )
	$(info Examples:)
	$(info   make migrate/book/create add_book_table)
	$(info   make migrate/book/up)
	$(info   make migrate/book/down)
	$(info   make migrate/book/status)
	$(info   make migrate/book)
	$(info =================================)
	@

# ===============================================================
# Quick Setup Commands
# ===============================================================

.PHONY: setup
setup: ## Initialize project with Docker
	$(info [SETUP] Setting up IKV project...)
	@if [ ! -f .env ]; then cp env.example .env; $(info [OK] Created .env from env.example); fi
	$(info [DOCKER] Starting Docker containers...)
	@docker-compose up -d
	$(info [WAIT] Waiting for services to be ready...)
	@sleep 10
	$(info [MIGRATE] Running migrations...)
	@make migrate/all/up
	$(info [DONE] Setup complete!)
	$(info [INFO] Application: http://localhost:7726)

.PHONY: dev
dev: ## Start development environment
	$(info [DEV] Starting development environment...)
	@docker-compose up

# ===============================================================
# Database Migration Commands
# ===============================================================

.PHONY: migrate/all/up
migrate/all/up: ## Run all migrations for all modules
	$(info [MIGRATE] Running all migrations for all modules...)
	@$(DOCKER_EXEC) $(GO_CMD) migrate up --all
	$(info [OK] All migrations completed)

.PHONY: migrate/all
migrate/all: ## Show pending migrations for all modules
	$(info [PENDING] Checking pending migrations for all modules...)
	@$(DOCKER_EXEC) $(GO_CMD) migrate pending --all

# ===============================================================
# Generic Module Migration Commands (for any module)
# ===============================================================

# Pattern rule for any module up command: migrate/{module}/up
migrate/%/up:
	$(info [MIGRATE] Running migrations for module: $*)
	@$(DOCKER_EXEC) $(GO_CMD) migrate up --module=$*
	$(info [OK] Migrations completed for module: $*)

# Pattern rule for any module down command: migrate/{module}/down
migrate/%/down:
	$(info [MIGRATE] Rolling back migration for module: $*)
ifdef n
	@$(DOCKER_EXEC) $(GO_CMD) migrate down --module=$* --steps=$(n)
else
	@$(DOCKER_EXEC) $(GO_CMD) migrate down --module=$* --steps=1
endif
	$(info [OK] Migration rolled back for module: $*)

# Pattern rule for any module create command: migrate/{module}/create
migrate/%/create:
	$(info [CREATE] Creating migration for module: $*)
	@$(DOCKER_EXEC) $(GO_CMD) migrate create --module=$* --name=$(filter-out $@,$(MAKECMDGOALS))
	$(info [OK] Migration files created)

# Pattern rule for any module status command: migrate/{module}/status  
migrate/%/status:
	$(info [STATUS] Migration status for module: $*)
	@$(DOCKER_EXEC) $(GO_CMD) migrate status --module=$*

# Pattern rule for any module pending command: migrate/{module}
migrate/%:
	$(info [PENDING] Checking pending migrations for module: $*)
	@$(DOCKER_EXEC) $(GO_CMD) migrate pending --module=$*

# Handle trailing arguments for create command
%:
	@:

# ===============================================================
# Development Commands
# ===============================================================

.PHONY: logs
logs: ## Show application logs
	@docker-compose logs -f web

.PHONY: shell
shell: ## Access development container shell
	@$(DOCKER_EXEC) /bin/sh

.PHONY: restart
restart: ## Restart development services
	@docker-compose restart web

.PHONY: clean
clean: ## Clean build artifacts and Docker
	$(info [CLEAN] Cleaning up...)
	@docker-compose down
	$(info [OK] Cleanup completed)

# Default target
.DEFAULT_GOAL := help 