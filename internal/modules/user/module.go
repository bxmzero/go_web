package user

import (
	"go.uber.org/fx"

	"example.com/gin-fx-sqlite-demo/internal/httpserver"
)

func newRoute(handler *Handler) httpserver.Route {
	return httpserver.Route{
		Register: handler.Register,
	}
}

var Module = fx.Module("user",
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
