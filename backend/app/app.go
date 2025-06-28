package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"reminder-app/controllers"
	"reminder-app/scheduler"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type App struct {
	db                  *pgxpool.Pool
	river               *river.Client[pgx.Tx]
	router              *gin.Engine
	reminderController  *controllers.ReminderController
	scheduler           *scheduler.Scheduler
}

func New(db *pgxpool.Pool, riverClient *river.Client[pgx.Tx], reminderController *controllers.ReminderController) *App {
	schedulerInstance := scheduler.NewScheduler(db, riverClient)
	return &App{
		db:                 db,
		river:              riverClient,
		reminderController: reminderController,
		scheduler:          schedulerInstance,
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
		api.GET("/reminders", a.GetReminders)
		api.POST("/reminders", a.CreateReminder)
		api.PUT("/reminders/:id", a.UpdateReminder)
		api.DELETE("/reminders/:id", a.DeleteReminder)
	}
}

// Handler methods
func (a *App) GetReminders(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	reminders, err := a.reminderController.GetReminders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reminders"})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (a *App) CreateReminder(c *gin.Context) {
	var reminder controllers.Reminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.reminderController.CreateReminder(&reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reminder"})
		return
	}

	c.JSON(http.StatusCreated, reminder)
}

func (a *App) UpdateReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var reminder controllers.Reminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.reminderController.UpdateReminder(id, &reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reminder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

func (a *App) DeleteReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := a.reminderController.DeleteReminder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reminder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
