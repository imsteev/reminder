package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"reminder-app/app"
	"reminder-app/controllers"
	"reminder-app/handler"
	"reminder-app/jobs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize River job queue
	riverClient, err := initRiver(db)
	if err != nil {
		log.Fatal("Failed to initialize River:", err)
	}
	defer riverClient.Stop(context.Background())

	// Initialize controller
	reminderController := controllers.NewReminderController(db)

	app := app.New(db, riverClient, reminderController)
	api := handler.New(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, api); err != nil {
		log.Fatal("Failed to run application:", err)
	}
}

func initDB() (*pgxpool.Pool, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost/reminder_app?sslmode=disable"
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func initRiver(db *pgxpool.Pool) (*river.Client[pgx.Tx], error) {
	workers := river.NewWorkers()
	river.AddWorker(workers, &jobs.ReminderWorker{})

	riverClient, err := river.NewClient(riverpgxv5.New(db), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		Workers: workers,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create river client: %w", err)
	}

	if err := riverClient.Start(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to start river client: %w", err)
	}

	return riverClient, nil
}
