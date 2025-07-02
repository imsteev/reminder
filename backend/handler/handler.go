package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reminder-app/controller/clerkcontroller"
	"reminder-app/controller/contactmethodcontroller"
	"reminder-app/controller/protocol"
	"reminder-app/controller/remindercontroller"
	"strconv"

	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Handler struct {
	*gin.Engine
	reminderController      *remindercontroller.Controller
	contactMethodController *contactmethodcontroller.Controller
	clerkController         *clerkcontroller.Controller
}

type Params struct {
	fx.In

	ReminderController      *remindercontroller.Controller
	ContactMethodController *contactmethodcontroller.Controller
	ClerkController         *clerkcontroller.Controller
}

var _ http.Handler = (*Handler)(nil)

func New(p Params) *Handler {
	api := gin.Default()
	h := &Handler{
		Engine:                  api,
		reminderController:      p.ReminderController,
		contactMethodController: p.ContactMethodController,
		clerkController:         p.ClerkController,
	}
	return h.init()
}

func (h *Handler) init() *Handler {
	clerk.SetKey("sk_test_LnQRf8v9rlcTLzegHKMbuTytkzgFTzumFmLOE2ckei")
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

	api := h.Group("/api")
	api.Use(clerkAuthMiddleware())
	api.GET("/reminders", h.handleGetReminders)
	api.POST("/reminders", h.handleCreateReminder)
	api.PUT("/reminders/:id", h.handleUpdateReminder)
	api.DELETE("/reminders/:id", h.handleDeleteReminder)
	api.GET("/contact-methods", h.handleGetContactMethods)
	api.POST("/contact-methods", h.handleCreateContactMethod)
	api.PUT("/contact-methods/:id", h.handleUpdateContactMethod)
	api.DELETE("/contact-methods/:id", h.handleDeleteContactMethod)

	webhooks := h.Group("/webhooks")
	webhooks.POST("/clerk", h.handleClerkWebhook)

	return h
}

func clerkAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, protocol.ErrorResponse{Error: "authorization header required"})
			c.Abort()
			return
		}

		// Wrap the request/response for clerk validation
		wrappedHandler := clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok {
				c.JSON(http.StatusUnauthorized, protocol.ErrorResponse{Error: "invalid or expired token"})
				c.Abort()
				return
			}

			c.Set("userID", claims.Subject)
		}))

		wrappedHandler.ServeHTTP(c.Writer, c.Request)

		// Only continue if auth was successful (userID was set)
		if _, exists := c.Get("userID"); exists {
			c.Next()
		}
	}
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

func (h *Handler) handleClerkWebhook(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("error reading clerk webhook body: ", err)
		c.Status(http.StatusOK)
		return
	}

	var eventType protocol.ClerkEvent
	if err := json.Unmarshal(body, &eventType); err != nil {
		fmt.Println("error parsing clerk webhook: ", err)
		c.Status(http.StatusOK)
		return
	}

	if err := h.clerkController.HandleClerkEvent(eventType.Type, body); err != nil {
		fmt.Println("error handling clerk webhook: ", err)
		c.Status(http.StatusOK)
		return
	}

	// always return 200 for webhooks
	c.Status(http.StatusOK)
}
