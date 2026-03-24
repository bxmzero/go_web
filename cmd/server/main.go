package main

import (
	"go.uber.org/fx"

	"example.com/gin-fx-sqlite-demo/internal/config"
	"example.com/gin-fx-sqlite-demo/internal/database"
	"example.com/gin-fx-sqlite-demo/internal/httpserver"
	"example.com/gin-fx-sqlite-demo/internal/modules/order"
	"example.com/gin-fx-sqlite-demo/internal/modules/user"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.New,
			database.NewSQLite,
			httpserver.NewEngine,
			httpserver.NewHTTPServer,
		),
		user.Module,
		order.Module,
		fx.Invoke(func(*httpserver.ServerReady) {}),
	)

	app.Run()
}
