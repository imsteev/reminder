package pgxpool

import (
	"go.uber.org/fx"
)

// Module defines the pgxpool fx module
var Module = fx.Module("pgxpool",
	fx.Provide(New),
)