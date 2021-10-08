package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Tests

func TestLoadConfig(t *testing.T) {
	err := LoadConfig()
	assert.NoError(t, err)

	assert.NotEmpty(t, viper.GetString("server.port"))

	assert.NotEmpty(t, viper.GetString("database.host"))
	assert.NotEmpty(t, viper.GetString("database.port"))
	assert.NotEmpty(t, viper.GetString("database.user"))
	assert.NotEmpty(t, viper.GetString("database.password"))
	assert.NotEmpty(t, viper.GetString("database.dbname"))

	assert.NotEmpty(t, viper.GetString("log.level"))
}

// Helpers

// TestMain main function
func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	// run tests
	code := m.Run()
	os.Exit(code)
}
