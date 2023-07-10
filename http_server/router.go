package http_server

import (
	"github.com/gin-gonic/gin"
)

//InitRouter InitRouter
func InitRouter(ginEngine *gin.Engine) (ginEngineDone *gin.Engine, err error) {

	ginEngine.GET("/ping", ping)
	ginEngine.POST("/write", write2)

	ginEngine.Use()
	ginEngineDone = ginEngine
	return
}
