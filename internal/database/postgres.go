package database

import (
	"fmt"

	"github.com/ai-readone/go-url-shortner/configs"
	"github.com/ai-readone/go-url-shortner/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var session *gorm.DB

func GetSession() *gorm.DB {
	return session
}

func InitPostgres(config *configs.Config) error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Database.Host,
		config.Database.Port, config.Database.User, config.Database.Password, config.Database.DbName)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	session = db.Session(&gorm.Session{})

	logger.Info("Successfully initialized Postgres")
	return nil
}
