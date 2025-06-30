package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reminder-app/config"
	"reminder-app/controller"
	gormmodule "reminder-app/db/gorm"
	"reminder-app/db/riverclient"
	"reminder-app/handler"
	"reminder-app/workers"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"go.uber.org/fx"
)

func StartRiverClient(riverClient *river.Client[pgx.Tx]) {
	fmt.Println("Starting River background job client...")

	if err := riverClient.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start River client: %v", err)
	}

	fmt.Println("River client started successfully")
}

func StartReminderService(cfg *config.Config, httpHandler *handler.Handler) {
	server := &http.Server{Addr: ":" + cfg.Port, Handler: httpHandler}

	fmt.Printf("Starting reminder service on port %s\n", cfg.Port)

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
		fx.Invoke(StartRiverClient),
		fx.Invoke(StartReminderService),
	)
	fxApp.Run()
}
