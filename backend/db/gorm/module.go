package gorm

import (
	"go.uber.org/fx"
)

// Module defines the GORM fx module
var Module = fx.Module("gorm",
	fx.Provide(New),
)