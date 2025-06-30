package handler

import (
	"reminder-app/controller"

	"github.com/samber/do/v2"
)

// NewHandler creates HTTP handler
func NewHandler(i do.Injector) (*Handler, error) {
	app, err := do.Invoke[*controller.App](i)
	if err != nil {
		return nil, err
	}

	return New(app), nil
}

// Package defines the handler dependency injection package
var Package = do.Package(
	do.Lazy(NewHandler),
)