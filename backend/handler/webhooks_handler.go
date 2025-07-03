package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reminder-app/controller/protocol"

	"github.com/gin-gonic/gin"
	svix "github.com/svix/svix-webhooks/go"
)

// log any errors but always return 200 for webhooks.
func (h *Handler) handleClerkWebhook(c *gin.Context) {
	payload, err := h.verifyClerkWebhook(c)
	if err != nil {
		fmt.Println("error verifying clerk webhook: ", err)
		c.Status(http.StatusOK)
		return
	}

	var eventType protocol.ClerkEvent
	if err := json.Unmarshal(payload, &eventType); err != nil {
		fmt.Println("error parsing clerk webhook: ", err)
		c.Status(http.StatusOK)
		return
	}

	if err := h.clerkController.HandleClerkEvent(eventType.Type, payload); err != nil {
		fmt.Println("error handling clerk webhook: ", err)
		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) verifyClerkWebhook(c *gin.Context) ([]byte, error) {
	wh, err := svix.NewWebhook(h.config.Clerk.WebhookSecretKey)
	if err != nil {
		return nil, err
	}
	headers := c.Request.Header
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	err = wh.Verify(payload, headers)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
