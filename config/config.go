package config

import (
	"bytes"
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"strings"
)

//go:embed application-config.yaml
var config []byte

func LoadConfig() (err error) {
	viper.SetConfigType("yaml")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()

	err = viper.ReadConfig(bytes.NewBuffer(config))
	if err != nil {
		return err
	}

	return
}

func SetLogLevel() {
	// Set log level for zerolog and gin
	switch strings.ToUpper(viper.GetString("log.level")) {
	case "DEBUG" :
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		gin.SetMode(gin.DebugMode)
	case "TEST" :
		zerolog.SetGlobalLevel(zerolog.Disabled)
		gin.SetMode(gin.TestMode)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}
}