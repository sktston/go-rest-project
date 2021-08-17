package router

import (
	_ "github.com/airoasis/go-rest-project/docs"
	"github.com/airoasis/go-rest-project/handler"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/heptiolabs/healthcheck"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"time"
)

func SetupRouter() *gin.Engine {
	health := healthcheck.NewHandler()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger.SetLogger(logger.WithWriter(zerolog.ConsoleWriter{Out:os.Stderr,TimeFormat: time.RFC3339})))

	ug := r.Group("/books")
	{
		ug.GET("", handler.GetBooks)
		ug.POST("", handler.CreateBook)
		ug.GET("/:id", handler.GetBookByID)
		ug.PUT("/:id", handler.UpdateBook)
		ug.DELETE("/:id", handler.DeleteBook)
	}

	// Health check apis for k8s
	hg := r.Group("/health")
	{
		hg.GET("/live", gin.WrapF(health.LiveEndpoint))
		hg.GET("/ready", gin.WrapF(health.ReadyEndpoint))
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
