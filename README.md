# Reminder App

A web application for managing recurring text message reminders, built with Go backend and React frontend.

## Quick Start

### Prerequisites
- Go 1.24+
- PostgreSQL
- Node.js 18+

### Setup
1. Clone and configure:
   ```bash
   cp .env.example .env
   # Update .env with your database credentials
   ```

2. Run migrations and start backend:
   ```bash
   cd backend
   go run ./cmd/migrate up
   go run main.go
   # Backend starts on http://localhost:8080
   ```

3. Start frontend:
   ```bash
   cd frontend
   npm install && npm run dev
   # Frontend starts on http://localhost:3000
   ```

## Migrations

```bash
# Basic commands
just migrate              # Run migrations
just migrate-reset        # Reset database
just migrate-new Feature  # Generate new migration

# Or use go directly
cd backend
go run ./cmd/migrate up
go run ./cmd/generate-migration FeatureName
```

To add a migration:
1. Run `just migrate-new YourFeature` (uses embedded Go templates)
2. Edit the generated file in `backend/db/migrate/`
3. Add constructor to `db/migrate/migrator.go`

## API

- `GET /api/reminders?user_id={id}` - Get reminders
- `POST /api/reminders` - Create reminder  
- `PUT /api/reminders/{id}` - Update reminder
- `DELETE /api/reminders/{id}` - Delete reminder

## Tech Stack

**Backend:** Go, Gin, PostgreSQL, GORM, River (background jobs)  
**Frontend:** React, TypeScript, Tailwind CSS, React Query
