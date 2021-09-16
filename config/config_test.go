package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	err := LoadConfig()
	assert.NoError(t, err)

	assert.NotEmpty(t, viper.GetString("server.port"))

	assert.NotEmpty(t, viper.GetString("database.host"))
	assert.NotEmpty(t, viper.GetString("database.port"))
	assert.NotEmpty(t, viper.GetString("database.user"))
	assert.NotEmpty(t, viper.GetString("database.password"))
	assert.NotEmpty(t, viper.GetString("database.dbname"))

	assert.NotEmpty(t, viper.GetString("test-database.host"))
	assert.NotEmpty(t, viper.GetString("test-database.port"))
	assert.NotEmpty(t, viper.GetString("test-database.user"))
	assert.NotEmpty(t, viper.GetString("test-database.password"))
	assert.NotEmpty(t, viper.GetString("test-database.dbname"))
}
