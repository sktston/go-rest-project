package config

import (
	"fmt"
	"github.com/airoasis/go-rest-project/model/entity"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Seoul",
		viper.GetString("database.host"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetInt("database.port"),
	)
}

func MigrateSchema() {
	DB.AutoMigrate(&entity.Book{})
}