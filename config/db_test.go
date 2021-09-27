package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Tests

func TestConnectingDatabase(t *testing.T) {
	testDB := InitTestDB(t)

	sqlDB, err := testDB.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())

	FreeTestDB(t, testDB)
}