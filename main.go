package main

import (
	_ "gin-admin-template/docs"
	"gin-admin-template/internal/api"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/router"
	"gin-admin-template/internal/service"
	"os"
	"strconv"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Admin API
// @version         1.0.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(ginzap.Ginzap(config.Logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(config.Logger, true))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.SetApiRouter(r)
	router.SetOtherRouter(r)
	service.SaveResourceFromSwagger("docs/swagger.json")

	dir, _ := os.Getwd()
	r.Static("/assets", dir+"/web/dist/assets")
	r.GET("/", api.HtmlHandler)

	err := service.InitAdminUser()
	if err != nil {
		config.Log.Errorf("Failed to initialize superadmin: %v", err)
		return
	} else {
		config.Log.Info("superadmin initialized successfully")
	}

	config.Log.Infof("Listening on %d", config.AppConfig.Server.Port)
	r.Run(":" + strconv.Itoa(config.AppConfig.Server.Port))
}
