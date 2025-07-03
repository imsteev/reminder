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

	return ctrl.db.Transaction(func(tx *gorm.DB) error {
		user.ClerkID = event.Data.ID
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}

		if len(event.Data.EmailAddresses) > 0 {
			var contactMethod models.ContactMethod
			contactMethod.Type = "email"
			contactMethod.Value = event.Data.EmailAddresses[0].EmailAddress
			contactMethod.UserID = int64(user.ID)
			contactMethod.Description = "Account email"
			if err := tx.Create(&contactMethod).Error; err != nil {
				return fmt.Errorf("error creating contact method: %w", err)
			}
		}
		return nil
	})
}
