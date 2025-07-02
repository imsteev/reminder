# Reminder App Development Commands
set dotenv-load

# Default recipe - show available commands
default:
    @just --list

# Backend Development
# ------------------

# Run backend server locally
backend-dev:
    cd backend && go run main.go

# Build backend binary
backend-build:
    cd backend && go build -o bin/reminder main.go

# Install backend dependencies
backend-deps:
    cd backend && go mod tidy

# Run database migrations
migrate:
    cd backend && go run ./cmd/migrate up

# Reset database (drop all tables and re-run migrations)
nuke:
    psql ${DATABASE_URL:-postgres://localhost/reminder?sslmode=disable} -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO public;"
    river migrate-up --line main --database-url "${DATABASE_URL:-postgres://localhost/reminder?sslmode=disable}"
    cd backend && go run ./cmd/migrate up

# Generate a new migration file (alternative syntax)
migrate-new MIGRATION_NAME:
    cd backend && go run ./cmd/generate-migration {{MIGRATION_NAME}}

# Frontend Development  
# --------------------

# Run frontend development server
frontend-dev:
    cd frontend && pnpm run dev

# Build frontend for production
frontend-build:
    cd frontend && pnpm run build

# Install frontend dependencies
frontend-deps:
    cd frontend && pnpm install

# Combined Operations
# ------------------

# Install all dependencies (backend + frontend)
install: backend-deps frontend-deps

# Build both backend and frontend
build: backend-build frontend-build

# Run both backend and frontend in parallel
dev:
    #!/usr/bin/env bash
    echo "Starting backend and frontend servers..."
    echo "Backend will run on http://localhost:8080"
    echo "Frontend will run on http://localhost:5173"
    echo ""
    echo "Press Ctrl+C to stop both servers"
    trap 'kill $(jobs -p)' EXIT
    just backend-dev & just frontend-dev & wait

# Clean build artifacts
clean:
    rm -rf backend/bin
    rm -rf frontend/dist