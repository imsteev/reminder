package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"reminder-app/app"
	"reminder-app/controller/remindercontroller"
	"reminder-app/db"
	"reminder-app/db/riverclient"
	"reminder-app/handler"
	"reminder-app/jobs"

	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
)

// API: handler -> app -> sub-controllers -> db
// Async: riverclient -> db <- workers
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
			dbURL = "postgres://localhost/reminder?sslmode=disable"
		}

		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}

	// Postgres
	dbConnections, err := db.New(dbURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer dbConnections.Close()

	// River
	riverClient, err := riverclient.New(dbConnections.PGXPool, &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		Workers: jobs.NewWorkers(),
	})
	if err != nil {
		log.Fatal("Failed to initialize River:", err)
	}
	defer riverClient.Stop(context.Background())

	// Wire everything together
	reminderController := remindercontroller.NewReminderController(dbConnections.GORM, riverClient)
	app := app.New(dbConnections.GORM, riverClient, reminderController)
	api := handler.New(app)

	if err := http.ListenAndServe(":"+port, api); err != nil {
		log.Fatal("Failed to run application:", err)
	}
}
