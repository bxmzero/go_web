package httpserver

import "github.com/gin-gonic/gin"

type Route struct {
	Register func(group *gin.RouterGroup)
}
