# Reminder App

A web application for managing recurring text message reminders, built with Go backend and React frontend.

## Features

- Create recurring reminders with custom schedules
- Specify frequency (times per day) and interval (hours between messages)
- Manage reminders through a clean web interface
- Background job processing for reliable message delivery
- PostgreSQL database for data persistence

## Tech Stack

### Backend
- **Go 1.24** - Main backend language
- **Gin** - HTTP web framework
- **PostgreSQL** - Database
- **River** - Background job processing with PostgreSQL

### Frontend
- **React 18** - Frontend framework
- **Tailwind CSS v4** - Styling
- **React Query** - Data fetching and state management
- **React Hook Form** - Form handling
- **Vite** - Build tool

## Getting Started

### Prerequisites
- Go 1.24+
- PostgreSQL
- Node.js 18+

### Backend Setup

1. Clone the repository and navigate to the project root
2. Copy environment variables:
   ```bash
   cp .env.example .env
   ```
3. Update `.env` with your database credentials
4. Run database migrations:
   ```bash
   psql -d your_database < migrations/001_create_reminders_table.sql
   ```
5. Install dependencies and run:
   ```bash
   go mod tidy
   go run main.go
   ```

Backend will start on http://localhost:8080

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm run dev
   ```

Frontend will start on http://localhost:3000

## API Endpoints

- `GET /api/reminders?user_id={id}` - Get reminders for a user
- `POST /api/reminders` - Create a new reminder
- `PUT /api/reminders/{id}` - Update a reminder
- `DELETE /api/reminders/{id}` - Delete a reminder

## Project Structure

```
├── main.go                          # Application entry point
├── go.mod                           # Go module file
├── .env.example                     # Environment variables template
├── migrations/                      # Database migrations
│   └── 001_create_reminders_table.sql
├── internal/
│   ├── controllers/                 # Business logic
│   │   └── reminder_controller.go
│   ├── handlers/                    # HTTP handlers
│   │   └── reminder_handler.go
│   ├── jobs/                        # Background jobs
│   │   └── reminder_job.go
│   └── scheduler/                   # Job scheduling
│       └── scheduler.go
└── frontend/                        # React frontend
    ├── package.json
    ├── vite.config.ts
    ├── index.html
    └── src/
        ├── main.tsx
        ├── App.tsx
        ├── index.css
        ├── api/
        │   └── reminders.ts
        └── components/
            ├── ReminderForm.tsx
            └── ReminderList.tsx
```

## Usage

1. Open the web interface at http://localhost:3000
2. Fill out the reminder form:
   - **Message**: The text to send
   - **Phone Number**: Recipient's phone number
   - **Times per day**: How many times to send daily
   - **Interval**: Hours between each message
   - **Start Time**: When to begin sending
3. Click "Create Reminder" to save
4. View and manage your reminders in the list

## Example Use Case

"I need to take medication 3 times a day, every 8 hours starting at 8 AM":
- Message: "Time to take your medication!"
- Phone Number: +1234567890
- Times per day: 3
- Interval: 8 hours
- Start Time: Today at 8:00 AM

This will send reminders at 8 AM, 4 PM, and 12 AM daily.

## Development

The application uses a clean architecture pattern:
- **Handlers** manage HTTP requests/responses
- **Controllers** contain business logic
- **Jobs** handle background processing
- **Scheduler** manages recurring tasks

The frontend uses modern React patterns with hooks and proper TypeScript types.