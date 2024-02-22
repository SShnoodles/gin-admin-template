package config

import (
	"gin-admin-template/internal/domain"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
	"time"
)

var DB *gorm.DB

func init() {
	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault()

	dsn := AppConfig.Datasource.Username + ":" + AppConfig.Datasource.Password + "@" + AppConfig.Datasource.Url
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		Log.Fatalf("Connection mysql:%s\n", err)
	}
	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	AutoMigrate()
}

func AutoMigrate() {
	DB.AutoMigrate(&domain.User{},
		&domain.Org{},
		&domain.Menu{},
		&domain.Role{},
		&domain.Resource{},
		&domain.RoleMenuRelation{},
		&domain.UserRoleRelation{},
		&domain.MenuResourceRelation{},
		&domain.OrgMenuRelation{},
	)
}
