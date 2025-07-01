package handler

import (
	"net/http"
	"reminder-app/controller/contactmethodcontroller"
	"reminder-app/controller/protocol"
	"reminder-app/controller/remindercontroller"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Handler struct {
	*gin.Engine
	reminderController      *remindercontroller.Controller
	contactMethodController *contactmethodcontroller.Controller
}

type Params struct {
	fx.In

	ReminderController      *remindercontroller.Controller
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
	var query protocol.GetRemindersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	reminders, err := h.reminderController.GetReminders(query.UserID, query.IncludePast)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (h *Handler) handleCreateReminder(c *gin.Context) {
	var reminder protocol.CreateReminderRequest
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	savedReminder, err := h.reminderController.CreateReminder(&reminder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, savedReminder)
}

func (h *Handler) handleUpdateReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.ErrorResponse{Error: "invalid reminder id"})
		return
	}

	var reminder protocol.UpdateReminderRequest
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	updatedReminder, err := h.reminderController.UpdateReminder(id, &reminder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedReminder)
}

func (h *Handler) handleDeleteReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.ErrorResponse{Error: "invalid reminder id"})
		return
	}

	if err := h.reminderController.DeleteReminder(id); err != nil {
		c.JSON(http.StatusInternalServerError, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, protocol.DeleteResponse{Message: "reminder deleted"})
}

func (h *Handler) handleGetContactMethods(c *gin.Context) {
	var query protocol.GetContactMethodsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	contactMethods, err := h.contactMethodController.GetContactMethods(query.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contactMethods)
}

func (h *Handler) handleCreateContactMethod(c *gin.Context) {
	var contactMethod protocol.CreateContactMethodRequest
	if err := c.ShouldBindJSON(&contactMethod); err != nil {
		c.JSON(http.StatusBadRequest, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	savedContactMethod, err := h.contactMethodController.CreateContactMethod(&contactMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, savedContactMethod)
}

func (h *Handler) handleUpdateContactMethod(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.ErrorResponse{Error: "invalid contact method id"})
		return
	}

	var contactMethod protocol.UpdateContactMethodRequest
	if err := c.ShouldBindJSON(&contactMethod); err != nil {
		c.JSON(http.StatusBadRequest, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	updatedContactMethod, err := h.contactMethodController.UpdateContactMethod(id, &contactMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedContactMethod)
}

func (h *Handler) handleDeleteContactMethod(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.ErrorResponse{Error: "invalid contact method id"})
		return
	}

	if err := h.contactMethodController.DeleteContactMethod(id); err != nil {
		c.JSON(http.StatusInternalServerError, protocol.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, protocol.DeleteResponse{Message: "contact method deleted"})
}
