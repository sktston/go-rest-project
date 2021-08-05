package config

import (
	"bytes"
	_ "embed"
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