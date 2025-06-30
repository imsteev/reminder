package gorm

import (
	"go.uber.org/fx"
)

var Module = fx.Module("gorm",
	fx.Provide(New),
)
