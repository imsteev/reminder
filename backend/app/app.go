package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"reminder-app/controllers"
	"reminder-app/handlers"
	"reminder-app/jobs"
	"reminder-app/scheduler"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

type App struct {
	db              *pgxpool.Pool
	river           *river.Client[pgx.Tx]
	router          *gin.Engine
	reminderHandler *handlers.ReminderHandler
}

func New() *App {
	return &App{}
}

func (a *App) Run() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	if err := a.initDB(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := a.initRiver(); err != nil {
		return fmt.Errorf("failed to initialize River: %w", err)
	}

	a.initHandlers()
	a.setupRoutes()

	// Start scheduler in background
	schedulerInstance := scheduler.NewScheduler(a.db, a.river)
	go schedulerInstance.Start(context.Background())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	return http.ListenAndServe(":"+port, a.router)
}

func (a *App) initDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost/reminder_app?sslmode=disable"
	}

	var err error
	a.db, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := a.db.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func (a *App) initRiver() error {
	workers := river.NewWorkers()
	river.AddWorker(workers, &jobs.ReminderWorker{})

	riverClient, err := river.NewClient(riverpgxv5.New(a.db), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		Workers: workers,
	})
	if err != nil {
		return fmt.Errorf("failed to create river client: %w", err)
	}

	if err := riverClient.Start(context.Background()); err != nil {
		return fmt.Errorf("failed to start river client: %w", err)
	}

	a.river = riverClient
	return nil
}

func (a *App) initHandlers() {
	reminderController := controllers.NewReminderController(a.db)
	a.reminderHandler = handlers.NewReminderHandler(reminderController)
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
