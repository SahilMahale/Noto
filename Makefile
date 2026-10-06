.PHONY: dev backend frontend stop install build-backend run-backend-bin \
	migrate-up migrate-down migrate-version \
	compose-up compose-down compose-down-volumes compose-watch compose-logs compose-ps \
	help

.DEFAULT_GOAL := help

DB_URL ?= postgres://noto:noto_dev_password@localhost:5432/noto?sslmode=disable

dev: ## Run backend and frontend natively (air + vite)
	@echo "Starting backend and frontend..."
	@make -j2 backend frontend

backend: ## Run backend only, native (air live reload on :8001)
	@echo "Starting backend on localhost:8001 (air)..."
	cd notes-backend && air

frontend: ## Run frontend only, native (vite dev server on :5173)
	@echo "Starting frontend..."
	cd notes-ui && bun run dev

stop: ## Stop natively-running backend/frontend processes
	@echo "Stopping services..."
	@-pkill -f "air" 2>/dev/null || true
	@-pkill -f "vite" 2>/dev/null || true
	@echo "Services stopped"

install: ## Install backend + frontend dependencies
	@echo "Installing backend dependencies..."
	cd notes-backend && go mod download
	@echo "Installing frontend dependencies..."
	cd notes-ui && bun install

build-backend: ## Build backend binary to notes-backend/bin/notes-server
	@echo "Building backend..."
	cd notes-backend && go build -o bin/notes-server ./cmd/main.go

run-backend-bin: ## Run the built backend binary
	@echo "Running backend binary..."
	./notes-backend/bin/notes-server

migrate-up: ## Apply all pending DB migrations (see notes-backend/migrations/README.md)
	migrate -path notes-backend/migrations -database "$(DB_URL)" up

migrate-down: ## Roll back one DB migration
	migrate -path notes-backend/migrations -database "$(DB_URL)" down 1

migrate-version: ## Print current DB migration version
	migrate -path notes-backend/migrations -database "$(DB_URL)" version

compose-up: ## Build and start the containerized dev stack (postgres, pgadmin, backend, frontend)
	podman compose up -d --build

compose-down: ## Stop the containerized dev stack (keeps volumes/data)
	podman compose down

compose-down-volumes: ## Stop the containerized dev stack and DELETE its volumes (wipes local DB data)
	podman compose down -v

compose-watch: ## Watch dependency manifests and auto-rebuild images (run alongside compose-up)
	podman compose watch

compose-logs: ## Tail logs from all containerized services
	podman compose logs -f

compose-ps: ## Show status of the containerized dev stack
	podman compose ps

help: ## Show this help
	@echo "Usage: make <target>"
	@echo
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
