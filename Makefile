.PHONY: all backend frontend dev stop clean migrate-up migrate-down migrate-version

DB_URL ?= postgres://noto:noto_dev_password@localhost:5432/noto?sslmode=disable

# Run both backend and frontend
dev:
	@echo "Starting backend and frontend..."
	@make -j2 backend frontend

# Run backend only
backend:
	@echo "Starting backend on localhost:8001 (air)..."
	cd notes-backend && air

# Run frontend only
frontend:
	@echo "Starting frontend..."
	cd notes-ui && bun run dev

# Stop all running services
stop:
	@echo "Stopping services..."
	@-pkill -f "air" 2>/dev/null || true
	@-pkill -f "vite" 2>/dev/null || true
	@echo "Services stopped"

# Install dependencies
install:
	@echo "Installing backend dependencies..."
	cd notes-backend && go mod download
	@echo "Installing frontend dependencies..."
	cd notes-ui && bun install

# Build backend
build-backend:
	@echo "Building backend..."
	cd notes-backend && go build -o bin/notes-server ./cmd/main.go

# Run built backend binary
run-backend-bin:
	@echo "Running backend binary..."
	./notes-backend/bin/notes-server

# Apply all pending DB migrations (see notes-backend/migrations/README.md)
migrate-up:
	migrate -path notes-backend/migrations -database "$(DB_URL)" up

# Roll back one DB migration
migrate-down:
	migrate -path notes-backend/migrations -database "$(DB_URL)" down 1

# Print current DB migration version
migrate-version:
	migrate -path notes-backend/migrations -database "$(DB_URL)" version
