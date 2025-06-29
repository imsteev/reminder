package app

import (
	"reminder-app/controller/remindercontroller"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type App struct {
	db                 *pgxpool.Pool
	river              *river.Client[pgx.Tx]
	reminderController *remindercontroller.Controller
}

func New(db *pgxpool.Pool, riverClient *river.Client[pgx.Tx], reminderController *remindercontroller.Controller) *App {
	return &App{
		db:                 db,
		river:              riverClient,
		reminderController: reminderController,
	}
}

func (a *App) GetReminders(userID string) ([]remindercontroller.Reminder, error) {
	return a.reminderController.GetReminders(userID)
}

func (a *App) CreateReminder(reminder *remindercontroller.Reminder) error {
	return a.reminderController.CreateReminder(reminder)
}

func (a *App) UpdateReminder(id int, reminder *remindercontroller.Reminder) error {
	return a.reminderController.UpdateReminder(id, reminder)
}

func (a *App) DeleteReminder(id int) error {
	return a.reminderController.DeleteReminder(id)
}
