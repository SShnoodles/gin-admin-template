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

func InitDB() error {
	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault()

	dsn := AppConfig.Datasource.Username + ":" + AppConfig.Datasource.Password + "@" + AppConfig.Datasource.Url
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	return AutoMigrate()
}

func AutoMigrate() error {
	return DB.AutoMigrate(&domain.User{},
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
