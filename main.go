package main

import (
	"github.com/airoasis/go-rest-project/config"
	"github.com/airoasis/go-rest-project/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

func main() {
	//Set a logger with zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out:os.Stderr,TimeFormat: time.RFC3339})

	//Load application config
	err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("cannot load config: %v", err)
	}

	// Set log level
	config.SetLogLevel(viper.GetString("log.level"))

	//Connect to DB
	config.DB, err = gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal().Msgf("Status: %v", err)
	}

	//Migrate the schema (Create the user table)
	config.MigrateSchema()

	//Start the gin server
	log.Info().Msgf("Starting the server")
	r := router.SetupRouter()
	r.Run(":" + viper.GetString("server.port"))
}