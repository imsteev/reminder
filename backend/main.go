package main

import (
	"context"
	"log"
	"net/http"
	"reminder-app/config"
	"reminder-app/controller"
	gormmodule "reminder-app/db/gorm"
	"reminder-app/handler"
	"reminder-app/river/riverclient"
	"reminder-app/workers"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func StartRiver(db *gorm.DB, riverClient *river.Client[pgx.Tx]) {
	log.Println("Starting River background job client...")

	if err := riverClient.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start River client: %v", err)
	}

	if err := workers.RestorePeriodicJobs(db, riverClient); err != nil {
		log.Fatalf("Failed to restore periodic jobs: %v", err)
	}

	log.Println("River client started successfully")
}

func StartReminderService(cfg *config.Config, httpHandler *handler.Handler) {
	log.Printf("Starting reminder service on port %s\n", cfg.Port)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: httpHandler,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}

func main() {
	fxApp := fx.New(
		config.Module,
		gormmodule.Module,
		riverclient.Module,
		workers.Module,
		controller.Module,
		handler.Module,
		fx.Invoke(StartRiver),
		fx.Invoke(StartReminderService),
	)
	fxApp.Run()
}
