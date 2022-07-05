package database

import (
	"fmt"

	"github.com/jgcaceres97/go-auth-jwt/src/models"
	"github.com/jgcaceres97/go-auth-jwt/src/settings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		*settings.DB.User,
		*settings.DB.Password,
		*settings.DB.Host,
		*settings.DB.Port,
		*settings.DB.Name,
	)

	conn, err := gorm.Open(
		mysql.Open(connectionString),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)},
	)

	if err != nil {
		panic("could not connect to the database")
	}

	DB = conn

	conn.AutoMigrate(&models.User{})
}
