package riverclient

import (
	"go.uber.org/fx"
)

// Module defines the river client fx module
var Module = fx.Module("riverclient",
	fx.Provide(New),
)