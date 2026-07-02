package app

import (
	"gin-admin-template/internal/api"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/middleware"
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

type App struct {
	Engine *gin.Engine
}

func New() (*App, error) {
	if err := config.Init(); err != nil {
		return nil, err
	}
	if err := middleware.InitValidator(); err != nil {
		return nil, err
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(ginzap.Ginzap(config.Logger, time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(config.Logger, true))
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.SetApiRouter(engine)
	router.SetOtherRouter(engine)
	registerStaticRoutes(engine)

	if err := service.SaveResourceFromSwagger("docs/swagger.json"); err != nil {
		config.Log.Error(err.Error())
	}
	if err := service.InitAdminUser(); err != nil {
		return nil, err
	}
	config.Log.Info("superadmin initialized successfully")

	return &App{
		Engine: engine,
	}, nil
}

func (a *App) Run() error {
	addr := ":" + strconv.Itoa(config.AppConfig.Server.Port)
	config.Log.Infof("Listening on %d", config.AppConfig.Server.Port)
	return a.Engine.Run(addr)
}

func registerStaticRoutes(engine *gin.Engine) {
	dir, err := os.Getwd()
	if err != nil {
		config.Log.Warnf("failed to get working directory: %v", err)
		return
	}
	engine.Static("/assets", dir+"/web/dist/assets")
	engine.GET("/", api.HtmlHandler)
}
