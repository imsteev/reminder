package handler

import (
	"net/http"
	"reminder-app/app"
	"reminder-app/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	*gin.Engine
	app *app.App
}

var _ http.Handler = (*Handler)(nil)

func New(app *app.App) *Handler {
	api := gin.Default()
	h := &Handler{
		Engine: api,
		app:    app,
	}
	return h.init()
}

func (h *Handler) init() *Handler {
	h.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	grp := h.Group("/api")
	grp.GET("/reminders", h.handleGetReminders)
	grp.POST("/reminders", h.handleCreateReminder)
	grp.PUT("/reminders/:id", h.handleUpdateReminder)
	grp.DELETE("/reminders/:id", h.handleDeleteReminder)

	return h
}

func (h *Handler) handleGetReminders(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	reminders, err := h.app.GetReminders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (h *Handler) handleCreateReminder(c *gin.Context) {
	var reminder models.Reminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.app.CreateReminder(&reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reminder)
}

func (h *Handler) handleUpdateReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reminder id"})
		return
	}

	var reminder models.Reminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reminder.BaseModel.ID = uint(id)
	if err := h.app.UpdateReminder(id, &reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminder)
}

func (h *Handler) handleDeleteReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reminder id"})
		return
	}

	if err := h.app.DeleteReminder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reminder deleted successfully"})
}
