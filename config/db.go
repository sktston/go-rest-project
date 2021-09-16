package config

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sktston/go-rest-project/model/entity"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB() error {
	// Connect to DB
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

	DB = db
	return nil
}

func MigrateSchema(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.Book{}); err != nil {
		return err
	}
	return nil
}

// Belows are for unit testing

func InitTestDB() (*gorm.DB, error) {
	testDBPrefix := uuid.New().String()+"_"
	testDsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Seoul",
		viper.GetString("test-database.host"),
		viper.GetString("test-database.user"),
		viper.GetString("test-database.password"),
		viper.GetString("test-database.dbname"),
		viper.GetInt("test-database.port"),
	)
	testDB, err := gorm.Open(postgres.Open(testDsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: testDBPrefix, // prefix is testId_
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	if err := MigrateSchema(testDB); err != nil {
		return nil, err
	}

	DB = testDB
	return testDB, nil
}

func FreeTestDB(testDB *gorm.DB) error {
	// Drop tables
	if err := testDB.Migrator().DropTable(&entity.Book{}); err != nil {
		return err
	}
	return nil
}
