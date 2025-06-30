package handler

import (
	"reminder-app/controller"

	"github.com/samber/do/v2"
)

// NewHandler creates HTTP handler (clean constructor)
func NewHandler(app *controller.App) *Handler {
	return New(app)
}

// newHandlerDI is a wrapper for DI that calls the clean constructor
func newHandlerDI(i do.Injector) (*Handler, error) {
	app, err := do.Invoke[*controller.App](i)
	if err != nil {
		return nil, err
	}

	return NewHandler(app), nil
}

// Package defines the handler dependency injection package
var Package = do.Package(
	do.Lazy(newHandlerDI),
)