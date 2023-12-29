package router

import (
	"gin-admin-template/internal/api"
	"gin-admin-template/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetOtherRouter(router *gin.Engine) {
	apiRouter := router.Group("/api")
	apiRouter.Use(middleware.CORS())
	{
		apiRouter.GET("/version", api.GetVersion)
	}
}
