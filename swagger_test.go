package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http/httptest"
	"testing"
)

func TestSwaggerDoc(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "SWAGGER_DISABLE"))

	r := httptest.NewRequest("GET", "/swagger/doc.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

