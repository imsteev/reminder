package app

import (
	"context"
	"log"
	"net/http"
	"os"

	"reminder-app/app/handlers"
	"reminder-app/scheduler"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type App struct {
	db              *pgxpool.Pool
	river           *river.Client[pgx.Tx]
	router          *gin.Engine
	reminderHandler *handlers.ReminderHandler
	scheduler       *scheduler.Scheduler
}

func New(db *pgxpool.Pool, riverClient *river.Client[pgx.Tx], reminderHandler *handlers.ReminderHandler) *App {
	schedulerInstance := scheduler.NewScheduler(db, riverClient)
	return &App{
		db:              db,
		river:           riverClient,
		reminderHandler: reminderHandler,
		scheduler:       schedulerInstance,
	}
}

func (a *App) Run() error {
	a.setupRoutes()

	// Start scheduler in background
	go a.scheduler.Start(context.Background())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	return http.ListenAndServe(":"+port, a.router)
}


func (a *App) setupRoutes() {
	a.router = gin.Default()

	a.router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := a.router.Group("/api")
	{
		api.GET("/reminders", a.reminderHandler.GetReminders)
		api.POST("/reminders", a.reminderHandler.CreateReminder)
		api.PUT("/reminders/:id", a.reminderHandler.UpdateReminder)
		api.DELETE("/reminders/:id", a.reminderHandler.DeleteReminder)
	}
}
