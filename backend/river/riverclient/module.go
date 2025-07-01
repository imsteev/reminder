package riverclient

import (
	"go.uber.org/fx"
)

var Module = fx.Module("riverclient",
	fx.Provide(New),
)
