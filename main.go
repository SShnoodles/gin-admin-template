package main

import (
	"gin-admin-template/internal/api"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/router"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	router.SetApiRouter(r)
	router.SetOtherRouter(r)

	dir, _ := os.Getwd()
	r.Static("/assets", dir+"/web/dist/assets")
	r.GET("/", api.HtmlHandler)

	r.Run(":" + strconv.Itoa(config.AppConfig.Server.Port))
}
