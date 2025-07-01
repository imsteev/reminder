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

func (cc *Controller) GetContactMethods(userID int64) ([]protocol.ContactMethod, error) {
	var dbContactMethods []models.ContactMethod
	err := cc.db.Where("user_id = ?", userID).Find(&dbContactMethods).Error
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

func (cc *Controller) CreateContactMethod(contactMethod *protocol.CreateContactMethodRequest) (*protocol.ContactMethod, error) {
	dbContactMethod := &models.ContactMethod{
		UserID:      contactMethod.UserID,
		Type:        contactMethod.Type,
		Value:       contactMethod.Value,
		Description: contactMethod.Description,
	}

	err := cc.db.Create(dbContactMethod).Error
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

func (cc *Controller) UpdateContactMethod(id int64, contactMethod *protocol.UpdateContactMethodRequest) (*protocol.ContactMethod, error) {
	var dbContactMethod models.ContactMethod
	if err := cc.db.Where("id = ?", id).First(&dbContactMethod).Error; err != nil {
		return nil, err
	}

	// Update fields from request
	dbContactMethod.Type = contactMethod.Type
	dbContactMethod.Value = contactMethod.Value
	dbContactMethod.Description = contactMethod.Description

	err := cc.db.Save(&dbContactMethod).Error
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

func (cc *Controller) DeleteContactMethod(id int64) error {
	err := cc.db.Delete(&models.ContactMethod{}, id).Error
	return err
}