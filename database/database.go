package database

import (
	"fmt"
	"github.com/spf13/viper"
	"go-rest-project/model/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
)

var gormDB *gorm.DB

func InitDB() error {
	// Connect to gormDB
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Seoul",
		viper.GetString("database.host"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetInt("database.port"),
	)

	// Set log level for gorm
	var level logger.LogLevel
	switch strings.ToUpper(viper.GetString("log.level")) {
	case "DEBUG":
		level = logger.Info
	case "TEST":
		level = logger.Silent
	default:
		level = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(level),
	})
	if err != nil {
		return err
	}

	// Migrate the schema (Create the book table)
	if err := MigrateSchema(db); err != nil {
		return err
	}

	SetDB(db)
	return nil
}

func GetDB() *gorm.DB {
	return gormDB
}

func SetDB(db *gorm.DB) {
	gormDB = db
}

func MigrateSchema(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.Book{}); err != nil {
		return err
	}
	return nil
}

func DropSchema(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&entity.Book{}); err != nil {
		return err
	}
	return nil
}
