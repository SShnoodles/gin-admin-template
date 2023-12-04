package initializations

import (
	"gin-admin-template/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

var DB *gorm.DB

func init() {
	dsn := AppConfig.Datasource.Username + ":" + AppConfig.Datasource.Password + "@" + AppConfig.Datasource.Url
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Connection mysql:%s\n", err)
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
	DB.AutoMigrate(&models.User{},
		&models.Org{},
		&models.Menu{},
		&models.Role{},
		&models.Resource{},
		&models.RoleMenuRelation{},
		&models.UserRoleRelation{},
		&models.MenuResourceRelation{},
		&models.OrgMenuRelation{},
	)
}
