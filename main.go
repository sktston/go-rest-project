package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sktston/go-rest-project/config"
	"github.com/sktston/go-rest-project/db"
	"github.com/sktston/go-rest-project/router"
	"github.com/spf13/viper"
	"os"
	"time"
)

// @title Go Rest Project API
// @version 0.1.0
func main() {
	//Set a logger with zerolog
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out:os.Stderr,TimeFormat: time.RFC3339})

	//Load application config
	if err := config.LoadConfig(); err != nil {
		log.Fatal().Msgf("cannot load config: %v", err)
	}

	//Connect to gormDB and Migrate Schema if not exist
	if err := db.InitDB(); err != nil {
		log.Fatal().Err(err).Caller().Msgf("cannot connect gormDB")
	}

	//Start the gin server
	log.Info().Msgf("Starting the server")
	r := router.SetupRouter()
	if err := r.Run(":" + viper.GetString("server.port")); err != nil {
		log.Fatal().Err(err).Caller().Msgf("cannot run server")
	}
}
