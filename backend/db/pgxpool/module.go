package pgxpool

import (
	"go.uber.org/fx"
)

var Module = fx.Module("pgxpool",
	fx.Provide(New),
)
