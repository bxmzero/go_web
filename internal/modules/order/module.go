package order

import (
	"go.uber.org/fx"

	"example.com/gin-fx-sqlite-demo/internal/httpserver"
)

func newRoute(handler *Handler) httpserver.Route {
	return httpserver.Route{
		Register: handler.Register,
	}
}

var Module = fx.Module("order",
	fx.Provide(
		NewRepository,
		NewService,
		NewHandler,
		fx.Annotate(
			newRoute,
			fx.ResultTags(`group:"routes"`),
		),
	),
)
