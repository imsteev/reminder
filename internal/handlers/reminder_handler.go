package handlers

import (
	"net/http"
	"strconv"

	"reminder-app/internal/controllers"

	"github.com/gin-gonic/gin"
)

type ReminderHandler struct {
	controller *controllers.ReminderController
}

func NewReminderHandler(controller *controllers.ReminderController) *ReminderHandler {
	return &ReminderHandler{controller: controller}
}

func (rh *ReminderHandler) GetReminders(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	reminders, err := rh.controller.GetReminders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reminders"})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (rh *ReminderHandler) CreateReminder(c *gin.Context) {
	var reminder controllers.Reminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := rh.controller.CreateReminder(&reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reminder"})
		return
	}

	c.JSON(http.StatusCreated, reminder)
}

func (rh *ReminderHandler) UpdateReminder(c *gin.Context) {
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

	if err := rh.controller.UpdateReminder(id, &reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reminder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

func (rh *ReminderHandler) DeleteReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := rh.controller.DeleteReminder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reminder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}