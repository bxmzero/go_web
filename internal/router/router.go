package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go_web/internal/handler"
)

func New(userHandler *handler.UserHandler) *gin.Engine {
	engine := gin.Default()

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	api := engine.Group("/api")
	userHandler.RegisterRoutes(api.Group("/users"))

	return engine
}
