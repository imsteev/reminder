package app

import (
	"reminder-app/controllers"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type App struct {
	db                 *pgxpool.Pool
	river              *river.Client[pgx.Tx]
	reminderController *controllers.ReminderController
}

func New(db *pgxpool.Pool, riverClient *river.Client[pgx.Tx], reminderController *controllers.ReminderController) *App {
	return &App{
		db:                 db,
		river:              riverClient,
		reminderController: reminderController,
	}
}

func (a *App) GetReminders(userID string) ([]controllers.Reminder, error) {
	return a.reminderController.GetReminders(userID)
}

func (a *App) CreateReminder(reminder *controllers.Reminder) error {
	return a.reminderController.CreateReminder(reminder)
}

func (a *App) UpdateReminder(id int, reminder *controllers.Reminder) error {
	return a.reminderController.UpdateReminder(id, reminder)
}

func (a *App) DeleteReminder(id int) error {
	return a.reminderController.DeleteReminder(id)
}
