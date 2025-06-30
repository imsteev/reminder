package main

import (
	"fmt"
	"log"
	"net/http"
	"reminder-app/config"
	"reminder-app/controller"
	gormmodule "reminder-app/db/gorm"
	"reminder-app/db/pgxpool"
	"reminder-app/db/riverclient"
	"reminder-app/handler"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Options(
	// Configuration
	config.Module,

	// Database layer
	pgxpool.Module,
	gormmodule.Module,
	riverclient.Module,

	// Business logic layer
	controller.Module,

	// Presentation layer
	handler.Module,
)

// StartReminderService starts the reminder service with all dependencies
func StartReminderService(
	cfg *config.Config,
	gormDB *gorm.DB,
	riverClient *river.Client[pgx.Tx],
	httpHandler *handler.Handler,
) {
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: httpHandler,
	}

	fmt.Printf("Starting reminder service on port %s\n", cfg.Port)
	
	// Start the server (this will block)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}

// API: handler -> app -> sub-controllers -> db
// Async: riverclient -> db <- workers
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Create fx app with all modules and lifecycle
	fxApp := fx.New(
		Module,
		fx.Invoke(StartReminderService),
	)

	// Start the fx app (blocks until stopped)
	fxApp.Run()
}
