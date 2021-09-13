package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sktston/go-rest-project/config"
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
	err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("cannot load config: %v", err)
	}

	//Connect to DB
	err = config.InitDB()
	if err != nil {
		log.Fatal().Err(err).Caller().Msgf("cannot connect DB")
	}

	//Start the gin server
	log.Info().Msgf("Starting the server")
	r := router.SetupRouter()
	err = r.Run(":" + viper.GetString("server.port"))
	if err != nil {
		log.Fatal().Err(err).Caller().Msgf("cannot run server")
	}
}
