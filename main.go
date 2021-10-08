package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sktston/go-rest-project/config"
	"github.com/sktston/go-rest-project/database"
	"github.com/sktston/go-rest-project/router"
	"github.com/spf13/viper"
	"os"
	"time"
)

// @title Go Rest Project API
// @version 0.1.0
func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out:os.Stderr,TimeFormat: time.RFC3339})

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

	//Start the gin server
	log.Info().Msgf("Starting the server")
	r := router.SetupRouter()
	if err := r.Run(":" + viper.GetString("server.port")); err != nil {
		log.Fatal().Err(err).Caller().Msgf("cannot run server")
	}
}
