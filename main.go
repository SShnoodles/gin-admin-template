package main

import (
	_ "gin-admin-template/docs"
	"gin-admin-template/internal/app"
	"log"
)

// @title           Admin API
// @version         1.0.0
// @BasePath        /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	if err := application.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
