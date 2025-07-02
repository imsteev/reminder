package controller

import (
	"reminder-app/controller/clerkcontroller"
	"reminder-app/controller/contactmethodcontroller"
	"reminder-app/controller/remindercontroller"

	"go.uber.org/fx"
)

var Module = fx.Module("controller",
	fx.Provide(
		remindercontroller.New,
		contactmethodcontroller.New,
		clerkcontroller.New,
	),
)
