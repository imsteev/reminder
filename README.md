# Reminder App

Recurring text message reminders with Go backend and React frontend.

## Setup

1. `cp .env.example .env` and update database credentials
2. `cd backend && go run ./cmd/migrate up && go run main.go`
3. `cd frontend && npm install && npm run dev`

## Stack

Go, Gin, PostgreSQL, React, TypeScript
