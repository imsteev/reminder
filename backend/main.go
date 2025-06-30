package main

import (
	"context"
	"log"
	"net/http"
	"reminder-app/config"
	"reminder-app/controller"
	"reminder-app/controller/remindercontroller"
	"reminder-app/db"
	"reminder-app/db/riverclient"
	"reminder-app/handler"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
)

// API: handler -> app -> sub-controllers -> db
// Async: riverclient -> db <- workers
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Create DI injector with all packages
	injector := do.New(
		config.Package,
		db.Package,
		riverclient.Package,
		remindercontroller.Package,
		controller.Package,
		handler.Package,
	)
	defer injector.Shutdown()

	// Get services from DI container
	cfg, err := do.Invoke[*config.Config](injector)
	if err != nil {
		log.Fatal("Failed to get config:", err)
	}

	dbConnections, err := do.Invoke[*db.Connections](injector)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer dbConnections.Close()

	riverClient, err := do.Invoke[*river.Client[pgx.Tx]](injector)
	if err != nil {
		log.Fatal("Failed to initialize River:", err)
	}
	defer riverClient.Stop(context.Background())

	api, err := do.Invoke[*handler.Handler](injector)
	if err != nil {
		log.Fatal("Failed to initialize handler:", err)
	}

	if err := http.ListenAndServe(":"+cfg.Port, api); err != nil {
		log.Fatal("Failed to run application:", err)
	}
}
