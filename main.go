package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sktston/go-rest-project/config"
	"github.com/sktston/go-rest-project/database"
	"github.com/sktston/go-rest-project/router"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Go Rest Project API
// @version 0.1.0
func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	//Load application config
	if err := config.LoadConfig(); err != nil {
		log.Fatal().Msgf("cannot load config: %v", err)
	}

	// Set log level for zerolog and gin
	config.SetLogLevel()

	//Connect to gormDB and Migrate Schema if not exist
	if err := database.InitDB(); err != nil {
		log.Fatal().Err(err).Caller().Msgf("cannot connect gormDB")
	}

	// Start the gin server
	// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-without-context/server.go
	log.Info().Msgf("Starting the server")
	r := router.SetupRouter()
	srv := &http.Server{
		Addr:    ":" + viper.GetString("server.port"),
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Caller().Msgf("cannot run server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msgf("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Caller().Msgf("Server forced to shutdown")
	}

	log.Info().Msgf("Server exiting")
}
