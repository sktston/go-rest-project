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

func SetLogLevel(level string) {
	if strings.ToUpper(level) == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else if strings.ToUpper(level) == "INFO" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	} else {
		// default is INFO
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}
}