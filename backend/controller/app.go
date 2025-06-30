package controller

import (
	"reminder-app/controller/remindercontroller"
	"reminder-app/models"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

// Top-level controller
type App struct {
	db                 *gorm.DB
	river              *river.Client[pgx.Tx]
	reminderController *remindercontroller.Controller
}

func NewApp(db *gorm.DB, riverClient *river.Client[pgx.Tx], reminderController *remindercontroller.Controller) *App {
	return &App{
		db:                 db,
		river:              riverClient,
		reminderController: reminderController,
	}
}

func (a *App) GetReminders(userID int64) ([]models.Reminder, error) {
	return a.reminderController.GetReminders(userID)
}

func (a *App) CreateReminder(reminder *models.Reminder) error {
	return a.reminderController.CreateReminder(reminder)
}

func (a *App) UpdateReminder(id int64, reminder *models.Reminder) error {
	return a.reminderController.UpdateReminder(id, reminder)
}

func (a *App) DeleteReminder(id int64) error {
	return a.reminderController.DeleteReminder(id)
}
