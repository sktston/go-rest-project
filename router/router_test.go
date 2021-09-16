package router

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/heptiolabs/healthcheck"
	"github.com/stretchr/testify/assert"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthLive(t *testing.T) {
	health := healthcheck.NewHandler()
	_, code := sendRequest(
		http.MethodGet,
		"/health/live",
		nil,
		setupRouter(http.MethodGet, "/health/live", gin.WrapF(health.LiveEndpoint)),
	)
	assert.Equal(t, http.StatusOK, code)
}

func TestHealthReady(t *testing.T) {
	health := healthcheck.NewHandler()
	_, code := sendRequest(
		http.MethodGet,
		"/health/ready",
		nil,
		setupRouter(http.MethodGet, "/health/ready", gin.WrapF(health.ReadyEndpoint)),
	)
	assert.Equal(t, http.StatusOK, code)
}

func TestSwaggerDoc(t *testing.T) {
	_, code := sendRequest(
		http.MethodGet,
		"/swagger/doc.json",
		nil,
		setupRouter(http.MethodGet, "/swagger/*any",ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "SWAGGER_DISABLE")),
	)
	assert.Equal(t, http.StatusOK, code)
}


// setupRouter get router on given handler
func setupRouter(httpMethod string, path string, handler gin.HandlerFunc) *gin.Engine {
	// prepare router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Handle(httpMethod, path, handler)
	return router
}

// sendRequest reads response from given http request.
func sendRequest(httpMethod string, target string, requestBody io.Reader, router *gin.Engine) (*bytes.Buffer, int) {
	// serve http on given response and request
	req := httptest.NewRequest(httpMethod, target, requestBody)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr.Body, rr.Code
}
