package main

import (
	"gin-admin-template/internal/api"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/router"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(ginzap.Ginzap(config.Logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(config.Logger, true))

	router.SetApiRouter(r)
	router.SetOtherRouter(r)

	dir, _ := os.Getwd()
	r.Static("/assets", dir+"/web/dist/assets")
	r.GET("/", api.HtmlHandler)

	config.Log.Infof("Listening on %d", config.AppConfig.Server.Port)
	r.Run(":" + strconv.Itoa(config.AppConfig.Server.Port))
}
