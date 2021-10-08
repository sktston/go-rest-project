package router

import (
	"github.com/gin-gonic/gin"
	"github.com/heptiolabs/healthcheck"
	"github.com/sktston/go-rest-project/config"
	"github.com/sktston/go-rest-project/test"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"testing"
)

// Tests

func TestHealthLive(t *testing.T) {
	health := healthcheck.NewHandler()
	_, code := test.SendRequest(
		http.MethodGet,
		"/health/live",
		nil,
		test.SetupRouter(http.MethodGet, "/health/live", gin.WrapF(health.LiveEndpoint)),
	)
	assert.Equal(t, http.StatusOK, code)
}

func TestHealthReady(t *testing.T) {
	health := healthcheck.NewHandler()
	_, code := test.SendRequest(
		http.MethodGet,
		"/health/ready",
		nil,
		test.SetupRouter(http.MethodGet, "/health/ready", gin.WrapF(health.ReadyEndpoint)),
	)
	assert.Equal(t, http.StatusOK, code)
}

func TestSwaggerDoc(t *testing.T) {
	_, code := test.SendRequest(
		http.MethodGet,
		"/swagger/doc.json",
		nil,
		test.SetupRouter(http.MethodGet, "/swagger/*any",ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "SWAGGER_DISABLE")),
	)
	assert.Equal(t, http.StatusOK, code)
}

// Helpers

// TestMain main function
func TestMain(m *testing.M) {
	viper.Set("log.level", "TEST")
	config.SetLogLevel()

	// run tests
	code := m.Run()
	os.Exit(code)
}
