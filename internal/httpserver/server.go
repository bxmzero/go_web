package httpserver

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"example.com/gin-fx-sqlite-demo/internal/config"
)

type serverParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    config.Config
	Engine    *gin.Engine
	Routes    []Route `group:"routes"`
}

type ServerReady struct {
	Server *http.Server
}

func NewEngine() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	return engine
}

func NewHTTPServer(params serverParams) (*ServerReady, error) {
	apiGroup := params.Engine.Group("/api/v1")
	for _, route := range params.Routes {
		route.Register(apiGroup)
	}

	params.Engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	server := &http.Server{
		Addr:    params.Config.HTTPAddr,
		Handler: params.Engine,
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Printf("http server listening on %s", params.Config.HTTPAddr)
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("http server stopped with error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return server.Shutdown(shutdownCtx)
		},
	})

	return &ServerReady{Server: server}, nil
}
