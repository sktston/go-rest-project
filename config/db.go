package config

import (
	"errors"
	"fmt"
	"github.com/sktston/go-rest-project/model/entity"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB
var TestPrefix = "test_"

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
	err = MigrateSchema(db)
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func InitTestDB() error {
	testDsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Seoul",
		viper.GetString("test-database.host"),
		viper.GetString("test-database.user"),
		viper.GetString("test-database.password"),
		viper.GetString("test-database.dbname"),
		viper.GetInt("test-database.port"),
	)
	testDb, err := gorm.Open(postgres.Open(testDsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: TestPrefix, // table name prefix
		},
	})
	if err != nil {
		return err
	}

	// Migrate the schema (Create the book table)
	err = MigrateSchema(testDb)
	if err != nil {
		return err
	}

	DB = testDb
	return nil
}

func MigrateSchema(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.Book{}); err != nil {
		return err
	}
	return nil
}

func FreeTestDB() error {
	// Check test DB
	if DB.NamingStrategy.TableName("") != TestPrefix {
		return errors.New("invalid db: current DB is not test DB")
	}

	// Drop tables
	if err := DB.Migrator().DropTable(&entity.Book{}); err != nil {
		return err
	}
	return nil
}
