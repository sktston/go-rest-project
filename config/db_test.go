package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnectingDatabase(t *testing.T) {
	assert.NoError(t, LoadConfig())

	testDB, err := InitTestDB()
	assert.NoError(t, err)

	sqlDB, err := testDB.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())

	assert.NoError(t, FreeTestDB(testDB))
}