package controller

import (
	"reminder-app/controller/remindercontroller"
	"reminder-app/models"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Top-level controller
type App struct {
	db                 *gorm.DB
	river              *river.Client[pgx.Tx]
	reminderController *remindercontroller.Controller
}

type Params struct {
	fx.In

	DB                 *gorm.DB
	River              *river.Client[pgx.Tx]
	ReminderController *remindercontroller.Controller
}

func New(p Params) *App {
	return &App{
		db:                 p.DB,
		river:              p.River,
		reminderController: p.ReminderController,
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
