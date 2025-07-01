package handler

import (
	"net/http"
	"reminder-app/controller/contactmethodcontroller"
	"reminder-app/controller/protocol"
	"reminder-app/controller/remindercontroller"
	"reminder-app/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Handler struct {
	*gin.Engine
	reminderController     *remindercontroller.Controller
	contactMethodController *contactmethodcontroller.Controller
}

type Params struct {
	fx.In

	ReminderController     *remindercontroller.Controller
	ContactMethodController *contactmethodcontroller.Controller
}

var _ http.Handler = (*Handler)(nil)

func New(p Params) *Handler {
	api := gin.Default()
	h := &Handler{
		Engine:                  api,
		reminderController:      p.ReminderController,
		contactMethodController: p.ContactMethodController,
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
	
	grp.GET("/contact-methods", h.handleGetContactMethods)
	grp.POST("/contact-methods", h.handleCreateContactMethod)
	grp.PUT("/contact-methods/:id", h.handleUpdateContactMethod)
	grp.DELETE("/contact-methods/:id", h.handleDeleteContactMethod)

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

	includePastStr := c.Query("include_past")
	includePast := includePastStr == "true"

	reminders, err := h.reminderController.GetReminders(userID, includePast)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (h *Handler) handleCreateReminder(c *gin.Context) {
	var reminder protocol.CreateReminderRequest
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	savedReminder, err := h.reminderController.CreateReminder(&reminder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, savedReminder)
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
	if err := h.reminderController.UpdateReminder(id, &reminder); err != nil {
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

	if err := h.reminderController.DeleteReminder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reminder deleted successfully"})
}

func (h *Handler) handleGetContactMethods(c *gin.Context) {
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

	contactMethods, err := h.contactMethodController.GetContactMethods(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contactMethods)
}

func (h *Handler) handleCreateContactMethod(c *gin.Context) {
	var contactMethod protocol.CreateContactMethodRequest
	if err := c.ShouldBindJSON(&contactMethod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	savedContactMethod, err := h.contactMethodController.CreateContactMethod(&contactMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, savedContactMethod)
}

func (h *Handler) handleUpdateContactMethod(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contact method id"})
		return
	}

	var contactMethod protocol.UpdateContactMethodRequest
	if err := c.ShouldBindJSON(&contactMethod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.contactMethodController.UpdateContactMethod(id, &contactMethod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "contact method updated successfully"})
}

func (h *Handler) handleDeleteContactMethod(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contact method id"})
		return
	}

	if err := h.contactMethodController.DeleteContactMethod(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "contact method deleted successfully"})
}
