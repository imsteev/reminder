package app

import (
	"reminder-app/controller/remindercontroller"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

type App struct {
	db                 *gorm.DB
	river              *river.Client[pgx.Tx]
	reminderController *remindercontroller.Controller
}

func New(db *gorm.DB, riverClient *river.Client[pgx.Tx], reminderController *remindercontroller.Controller) *App {
	return &App{
		db:                 db,
		river:              riverClient,
		reminderController: reminderController,
	}
}

func (a *App) GetReminders(userID int64) ([]remindercontroller.Reminder, error) {
	return a.reminderController.GetReminders(userID)
}

func (a *App) CreateReminder(reminder *remindercontroller.Reminder) error {
	return a.reminderController.CreateReminder(reminder)
}

func (a *App) UpdateReminder(id int64, reminder *remindercontroller.Reminder) error {
	return a.reminderController.UpdateReminder(id, reminder)
}

func (a *App) DeleteReminder(id int64) error {
	return a.reminderController.DeleteReminder(id)
}
