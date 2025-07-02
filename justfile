set dotenv-filename := "backend/.env"

# Default recipe - show available commands
default:
    @just --list

# Backend 

# Run backend server locally
backend-dev:
    cd backend && go run main.go

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

# Frontend

# Run frontend development server locally
frontend-dev:
    cd frontend && pnpm run dev

# Run both backend and frontend locally in parallel
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