package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"reminder-app/app"
	"reminder-app/controller/remindercontroller"
	"reminder-app/db"
	"reminder-app/db/riverclient"
	"reminder-app/handler"
	"reminder-app/jobs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
	}

	var (
		dbURL string
		port  string
	)
	{
		dbURL = os.Getenv("DATABASE_URL")
		if dbURL == "" {
			dbURL = "postgres://localhost/reminder_app?sslmode=disable"
		}

		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}

	db, err := db.New(dbURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	riverClient, err := initRiver(db)
	if err != nil {
		log.Fatal("Failed to initialize River:", err)
	}
	defer riverClient.Stop(context.Background())

	// Wire everything together
	reminderController := remindercontroller.NewReminderController(db)
	app := app.New(db, riverClient, reminderController)
	api := handler.New(app)

	if err := http.ListenAndServe(":"+port, api); err != nil {
		log.Fatal("Failed to run application:", err)
	}
}

func initRiver(db *pgxpool.Pool) (*river.Client[pgx.Tx], error) {
	workers := river.NewWorkers()
	river.AddWorker(workers, &jobs.ReminderWorker{})
	riverClient, err := riverclient.New(db, &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		Workers: workers,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create river client: %w", err)
	}

	return riverClient, nil
}
