package contactmethodcontroller

import "go.uber.org/fx"

var Module = fx.Module("contactmethodcontroller",
	fx.Provide(New),
)