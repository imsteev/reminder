package migrate

import (
	"reminder-app/models"

	"gorm.io/gorm"
)

var (
	Plan202412291545 = NewMigrationPlan("202412291545", Up202412291545, Down202412291545)
)

// Up202412291545 creates the initial schema with ENUMs, tables, indexes, and test user
func Up202412291545(tx *gorm.DB) error {
	// First, create the ENUM types if they don't exist
	if err := tx.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'contact_type') THEN
				CREATE TYPE contact_type AS ENUM ('phone', 'email');
			END IF;
		END $$;
	`).Error; err != nil {
		return err
	}

	if err := tx.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'reminder_type') THEN
				CREATE TYPE reminder_type AS ENUM ('one-time', 'repeating');
			END IF;
		END $$;
	`).Error; err != nil {
		return err
	}

	if err := tx.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'delivery_type') THEN
				CREATE TYPE delivery_type AS ENUM ('sms', 'email');
			END IF;
		END $$;
	`).Error; err != nil {
		return err
	}

	// Create tables using GORM's AutoMigrate
	if err := tx.AutoMigrate(
		&models.User{},
		&models.ContactMethod{},
		&models.Reminder{},
	); err != nil {
		return err
	}

	// Always ensure test user exists with ID 1
	testUser := &models.User{
		ID:   1,
		Name: "Test User",
	}
	
	// Use FirstOrCreate to avoid duplicates if migration is run multiple times
	if err := tx.FirstOrCreate(testUser, "id = ?", 1).Error; err != nil {
		return err
	}

	// Create test contact methods for both SMS and email delivery
	contactMethods := []models.ContactMethod{
		{
			UserID:      1,
			Type:        "phone",
			Value:       "+1234567890",
			Description: "Primary phone",
		},
		{
			UserID:      1,
			Type:        "email",
			Value:       "test@example.com",
			Description: "Primary email",
		},
	}

	for _, cm := range contactMethods {
		if err := tx.FirstOrCreate(&cm, "user_id = ? AND type = ?", cm.UserID, cm.Type).Error; err != nil {
			return err
		}
	}

	// Add performance indexes for reminders
	if err := tx.Exec(`
		CREATE INDEX IF NOT EXISTS idx_reminders_user_start_time 
		ON reminders(user_id, start_time) 
		WHERE deleted_at IS NULL
	`).Error; err != nil {
		return err
	}

	if err := tx.Exec(`
		CREATE INDEX IF NOT EXISTS idx_reminders_type_delivery 
		ON reminders(type, delivery_type) 
		WHERE deleted_at IS NULL
	`).Error; err != nil {
		return err
	}

	return nil
}

// Down202412291545 rolls back the initial schema
func Down202412291545(tx *gorm.DB) error {
	// Drop custom indexes first
	if err := tx.Exec("DROP INDEX IF EXISTS idx_reminders_user_start_time").Error; err != nil {
		return err
	}

	if err := tx.Exec("DROP INDEX IF EXISTS idx_reminders_type_delivery").Error; err != nil {
		return err
	}

	// Drop tables in reverse order
	if err := tx.Migrator().DropTable(&models.Reminder{}); err != nil {
		return err
	}

	if err := tx.Migrator().DropTable(&models.ContactMethod{}); err != nil {
		return err
	}

	if err := tx.Migrator().DropTable(&models.User{}); err != nil {
		return err
	}

	// Drop ENUM types
	if err := tx.Exec("DROP TYPE IF EXISTS delivery_type").Error; err != nil {
		return err
	}

	if err := tx.Exec("DROP TYPE IF EXISTS reminder_type").Error; err != nil {
		return err
	}

	if err := tx.Exec("DROP TYPE IF EXISTS contact_type").Error; err != nil {
		return err
	}

	return nil
}
