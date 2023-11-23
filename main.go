package main

import (
	"gin-admin-template/controllers"
	"gin-admin-template/initializations"
	"gin-admin-template/routers"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	routers.SetApiRouter(router)
	routers.SetOtherRouter(router)

	dir, _ := os.Getwd()
	router.Static("/assets", dir+"/web/dist/assets")
	router.GET("/", controllers.HtmlHandler)

	router.Run(":" + strconv.Itoa(initializations.AppConfig.Server.Port))
}
