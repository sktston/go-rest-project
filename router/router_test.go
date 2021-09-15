package router

import (
	"github.com/gin-gonic/gin"
	"github.com/heptiolabs/healthcheck"
	"github.com/stretchr/testify/assert"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthLive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	health := healthcheck.NewHandler()
	r := gin.New()
	r.GET("/health/live", gin.WrapF(health.LiveEndpoint))

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHealthReady(t *testing.T) {
	gin.SetMode(gin.TestMode)
	health := healthcheck.NewHandler()
	r := gin.New()
	r.GET("/health/ready", gin.WrapF(health.ReadyEndpoint))

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSwaggerDoc(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "SWAGGER_DISABLE"))

	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
