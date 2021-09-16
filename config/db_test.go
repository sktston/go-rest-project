package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnectingDatabase(t *testing.T) {
	err := LoadConfig()
	assert.NoError(t, err)

	testDB, err := InitTestDB()
	assert.NoError(t, err)

	sqlDB, err := GetDB().DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())

	err = FreeTestDB(testDB)
	assert.NoError(t, err)
}