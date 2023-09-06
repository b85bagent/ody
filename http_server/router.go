package http_server

import (
	"newProject/handler"

	"github.com/gin-gonic/gin"
)

type HandlerWithServer func(c *gin.Context, server *handler.Server)

func WithServer(h HandlerWithServer, server *handler.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c, server)
	}
}

// InitRouter InitRouter
func InitRouter(ginEngine *gin.Engine, server *handler.Server) (ginEngineDone *gin.Engine, err error) {

	ginEngine.GET("/ping", ping)

	v1 := ginEngine.Group("/api/v1")
	{
		v1.GET("/ping", ping)
		// v1.POST("/write", WithServer(writeWithIndex, server))
		// v1.POST("/write/:index", WithServer(writeWithIndex, server))
	}

	ginEngine.Use()
	ginEngineDone = ginEngine
	return
}
