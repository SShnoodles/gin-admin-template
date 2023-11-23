package routers

import (
	"gin-admin-template/controllers"
	"gin-admin-template/middlewares"
	"github.com/gin-gonic/gin"
)

func SetOtherRouter(router *gin.Engine) {
	apiRouter := router.Group("/api")
	apiRouter.Use(middlewares.CORS())
	{
		apiRouter.GET("/version", controllers.GetVersion)
	}
}
