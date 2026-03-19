package app

import (
	"go_web/internal/db"
	"go_web/internal/handler"
	"go_web/internal/repository"
	"go_web/internal/router"
	"go_web/internal/service"
	"go_web/internal/txmanager"

	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
}

func New() (*App, error) {
	database, err := db.NewSQLite()
	if err != nil {
		return nil, err
	}

	txMgr := txmanager.New(database)
	userRepo := repository.NewUserRepository(txMgr)
	userService := service.NewUserService(userRepo, txMgr)
	userHandler := handler.NewUserHandler(userService)

	return &App{engine: router.New(userHandler)}, nil
}

func (a *App) Run() error {
	return a.engine.Run(":8080")
}
