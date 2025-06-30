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
	"github.com/riverqueue/river"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func StartReminderService(
	cfg *config.Config,
	gormDB *gorm.DB,
	riverClient *river.Client[pgx.Tx],
	httpHandler *handler.Handler,
) {
	server := &http.Server{Addr: ":" + cfg.Port, Handler: httpHandler}

	fmt.Printf("Starting reminder service on port %s\n", cfg.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}

func main() {
	fxApp := fx.New(
		config.Module,
		pgxpool.Module,
		gormmodule.Module,
		riverclient.Module,
		controller.Module,
		handler.Module,
		fx.Invoke(StartReminderService),
	)
	fxApp.Run()
}
