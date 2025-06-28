# Reminder App Development Commands

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

# Run backend tests
backend-test:
    cd backend && go test ./...

# Install backend dependencies
backend-deps:
    cd backend && go mod tidy

# Frontend Development  
# --------------------

# Run frontend development server
frontend-dev:
    cd frontend && pnpm run dev

# Build frontend for production
frontend-build:
    cd frontend && pnpm run build

# Run frontend linter
frontend-lint:
    cd frontend && pnpm run lint

# Preview frontend production build
frontend-preview:
    cd frontend && pnpm run preview

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