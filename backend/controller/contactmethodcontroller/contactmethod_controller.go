package contactmethodcontroller

import (
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

func (ctrl *Controller) GetContactMethods(clerkID string) ([]protocol.ContactMethod, error) {
	var user models.User
	if err := ctrl.db.Where("clerk_id = ?", clerkID).First(&user).Error; err != nil {
		return nil, err
	}

	var dbContactMethods []models.ContactMethod
	err := ctrl.db.Where("user_id = ?", user.ID).Find(&dbContactMethods).Error
	if err != nil {
		return nil, err
	}

	var protocolContactMethods []protocol.ContactMethod
	for _, dbContactMethod := range dbContactMethods {
		protocolContactMethods = append(protocolContactMethods, protocol.ContactMethod{
			ID:          int64(dbContactMethod.ID),
			UserID:      dbContactMethod.UserID,
			Type:        dbContactMethod.Type,
			Value:       dbContactMethod.Value,
			Description: dbContactMethod.Description,
		})
	}
	return protocolContactMethods, nil
}

func (ctrl *Controller) CreateContactMethod(clerkID string, contactMethod *protocol.CreateContactMethodRequest) (*protocol.ContactMethod, error) {
	var user models.User
	if err := ctrl.db.Where("clerk_id = ?", clerkID).First(&user).Error; err != nil {
		return nil, err
	}

	dbContactMethod := &models.ContactMethod{
		UserID:      int64(user.ID),
		Type:        contactMethod.Type,
		Value:       contactMethod.Value,
		Description: contactMethod.Description,
	}

	err := ctrl.db.Create(dbContactMethod).Error
	if err != nil {
		return nil, err
	}

	return &protocol.ContactMethod{
		ID:          int64(dbContactMethod.ID),
		UserID:      dbContactMethod.UserID,
		Type:        dbContactMethod.Type,
		Value:       dbContactMethod.Value,
		Description: dbContactMethod.Description,
	}, nil
}

func (ctrl *Controller) UpdateContactMethod(id int64, contactMethod *protocol.UpdateContactMethodRequest) (*protocol.ContactMethod, error) {
	// todo: validate that the contact method belongs to the user

	var dbContactMethod models.ContactMethod
	if err := ctrl.db.Where("id = ?", id).First(&dbContactMethod).Error; err != nil {
		return nil, err
	}

	// Update fields from request
	dbContactMethod.Type = contactMethod.Type
	dbContactMethod.Value = contactMethod.Value
	dbContactMethod.Description = contactMethod.Description

	err := ctrl.db.Save(&dbContactMethod).Error
	if err != nil {
		return nil, err
	}

	return &protocol.ContactMethod{
		ID:          int64(dbContactMethod.ID),
		UserID:      dbContactMethod.UserID,
		Type:        dbContactMethod.Type,
		Value:       dbContactMethod.Value,
		Description: dbContactMethod.Description,
	}, nil
}

func (ctrl *Controller) DeleteContactMethod(id int64) error {
	// todo: validate that the contact method belongs to the user

	err := ctrl.db.Delete(&models.ContactMethod{}, id).Error
	return err
}
