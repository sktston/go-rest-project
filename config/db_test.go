package config

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"testing"
)

// Tests

func TestConnectingDatabase(t *testing.T) {
	// Load configuration
	assert.NoError(t, LoadConfig())

	// Open test DB with random prefix
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
	assert.NoError(t, err)

	// Migrate the schema
	assert.NoError(t, MigrateSchema(testDB))

	// Test
	sqlDB, err := testDB.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())

	// Free
	assert.NoError(t, DropSchema(testDB))
}