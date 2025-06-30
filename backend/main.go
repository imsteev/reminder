package main

import (
	"context"
	"fmt"
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
	"go.uber.org/fx"
)

// Module defines all fx options for the complete application
var Module = fx.Options(
	// Configuration
	config.Module,

	// Database layer
	db.Module,
	riverclient.Module,

	// Business logic layer
	remindercontroller.Module,
	controller.Module,

	// Presentation layer
	handler.Module,
)

// API: handler -> app -> sub-controllers -> db
// Async: riverclient -> db <- workers
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Create fx app with all modules and lifecycle
	fxApp := fx.New(
		Module,
		fx.Invoke(func(
			cfg *config.Config,
			dbConnections *db.Connections,
			riverClient *river.Client[pgx.Tx],
			httpHandler *handler.Handler,
			lc fx.Lifecycle,
		) {
			server := &http.Server{
				Addr:    ":" + cfg.Port,
				Handler: httpHandler,
			}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					fmt.Printf("Starting application on port %s\n", cfg.Port)
					go func() {
						if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
							log.Printf("HTTP server error: %v", err)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					fmt.Println("Shutting down application...")

					// Shutdown HTTP server
					if err := server.Shutdown(ctx); err != nil {
						log.Printf("HTTP server shutdown error: %v", err)
					}

					// Stop River client
					if riverClient != nil {
						riverClient.Stop(ctx)
					}

					// Close database connections
					if dbConnections != nil {
						dbConnections.Close()
					}

					fmt.Println("Application shut down complete")
					return nil
				},
			})
		}),
	)

	// Start the fx app (blocks until stopped)
	fxApp.Run()
}
