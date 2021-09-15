package router

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/heptiolabs/healthcheck"
	"github.com/rs/zerolog"
	_ "github.com/sktston/go-rest-project/docs"
	"github.com/sktston/go-rest-project/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"time"
)

func SetupRouter() *gin.Engine {
	health := healthcheck.NewHandler()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.SetLogger(logger.WithWriter(zerolog.ConsoleWriter{Out: os.Stderr,TimeFormat: time.RFC3339})))

	// Books
	router.POST("/books", handler.CreateBook)
	router.GET("/books", handler.GetBookList)
	router.GET("/books/:id", handler.GetBookByID)
	router.PUT("/books/:id", handler.UpdateBook)
	router.DELETE("/books/:id", handler.DeleteBook)

	// Health check apis for k8s
	router.GET("/health/live", gin.WrapF(health.LiveEndpoint))
	router.GET("/health/ready", gin.WrapF(health.ReadyEndpoint))

	// Swagger
	router.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "SWAGGER_DISABLE"))

	return router
}
