package controller

import (
	"reminder-app/controller/remindercontroller"

	"go.uber.org/fx"
)

var Module = fx.Module("controller",
	fx.Provide(
		remindercontroller.New,
		New,
	),
)
