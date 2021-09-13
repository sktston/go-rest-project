package router

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHealthLive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	r := httptest.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestHealthReady(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	r := httptest.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestSwaggerDoc(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	r := httptest.NewRequest("GET", "/swagger/doc.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}
