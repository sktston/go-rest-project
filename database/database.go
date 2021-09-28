package database

import (
	"fmt"
	"github.com/sktston/go-rest-project/model/entity"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
