package handler

import (
	"reminder-app/controller"

	"go.uber.org/fx"
)

type HandlerParams struct {
	fx.In

	App *controller.App
}

// NewHandler creates HTTP handler (clean constructor)
func NewHandler(p HandlerParams) *Handler {
	return New(p.App)
}

// Module defines the handler fx module
var Module = fx.Module("handler",
	fx.Provide(NewHandler),
)