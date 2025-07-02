package clerkcontroller

import (
	"encoding/json"
	"fmt"
	"reminder-app/controller/protocol"
	"reminder-app/models"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Controller struct {
	db *gorm.DB
}

type Params struct {
	fx.In

	DB *gorm.DB
}

func New(p Params) *Controller {
	return &Controller{db: p.DB}
}

func (ctrl *Controller) HandleClerkEvent(eventType string, body []byte) error {
	switch eventType {
	case "user.created":
		var event protocol.ClerkUserCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return fmt.Errorf("error unmarshalling clerk webhook: %w", err)
		}
		return ctrl.onUserCreated(event)
	default:
		return fmt.Errorf("unexpected event type: %s", eventType)
	}
}

func (ctrl *Controller) onUserCreated(event protocol.ClerkUserCreatedEvent) error {
	var user models.User
	user.ClerkID = event.Data.ID
	return ctrl.db.Create(&user).Error
}
