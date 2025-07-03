package migrate

import (
	"reminder-app/models"
	"slices"

	"gorm.io/gorm"
)

var (
	Plan202412291545 = NewMigrationPlan("202412291545", Up202412291545, Down202412291545)
)

func init() {
	if !slices.ContainsFunc(plans, func(p *MigrationPlan) bool {
		return p.ID == Plan202412291545.ID
	}) {
		panic("Plan202412291545 is not registered")
	}
}

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
		BaseModel: models.BaseModel{
			ID: 1,
		},
		Name:    "Stephen Chung",
		ClerkID: "user_2vJdLwhivmM8WJFmu0uxozMv7M6",
	}

	// Use FirstOrCreate to avoid duplicates if migration is run multiple times
	if err := tx.FirstOrCreate(testUser, "id = ?", 1).Error; err != nil {
		return err
	}

	// Bump the users sequence to avoid conflicts with the test user
	if err := tx.Exec("SELECT setval('users_id_seq', 2, false)").Error; err != nil {
		return err
	}

	// Create test contact methods for both SMS and email delivery
	contactMethods := []models.ContactMethod{
		{
			UserID:      1,
			Type:        "email",
			Value:       "spchung95@gmail.com",
			Description: "Primary email",
		},
	}

	for _, cm := range contactMethods {
		if err := tx.FirstOrCreate(&cm, "user_id = ? AND type = ?", cm.UserID, cm.Type).Error; err != nil {
			return err
		}
	}

	return nil
}

// Down202412291545 rolls back the initial schema
func Down202412291545(tx *gorm.DB) error {
	// Drop custom indexes first
	if err := tx.Exec("DROP INDEX IF EXISTS idx_reminders_user_start_time").Error; err != nil {
		return err
	}

	if err := tx.Exec("DROP INDEX IF EXISTS idx_reminders_repeating_contact").Error; err != nil {
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
	if err := tx.Exec("DROP TYPE IF EXISTS contact_type").Error; err != nil {
		return err
	}

	return nil
}
