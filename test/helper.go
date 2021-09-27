package test

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sktston/go-rest-project/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"net/http/httptest"
	"testing"
)

// InitTestDB init test database
func InitTestDB(t *testing.T) *gorm.DB {
	// Load configuration
	assert.NoError(t, config.LoadConfig())

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
	assert.NoError(t, config.MigrateSchema(testDB))

	config.SetDB(testDB)
	return testDB
}

// FreeTestDB free test database
func FreeTestDB(t *testing.T, testDB *gorm.DB) {
	assert.NoError(t, config.DropSchema(testDB))
}

// SetupRouter get router on given handler
func SetupRouter(httpMethod string, path string, handler gin.HandlerFunc) *gin.Engine {
	// prepare router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Handle(httpMethod, path, handler)
	return router
}

// SendRequest reads response from given http request.
func SendRequest(httpMethod string, target string, requestBody io.Reader, router *gin.Engine) (*bytes.Buffer, int) {
	// serve http on given response and request
	req := httptest.NewRequest(httpMethod, target, requestBody)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr.Body, rr.Code
}
