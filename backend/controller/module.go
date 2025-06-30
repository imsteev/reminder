package controller

import (
	"reminder-app/controller/remindercontroller"

	"go.uber.org/fx"
)

// Module defines the app controller fx module
var Module = fx.Module("controller",
	fx.Provide(
		remindercontroller.New,
		New,
	),
)
