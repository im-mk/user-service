package config

import (
	"log"

	"github.com/spf13/viper"
)

type ApplicationConfig struct {
	Port string
}

func GetConfig() ApplicationConfig {

	viper.SetConfigType("json")
	viper.SetConfigFile(".env")

	appConfig := ApplicationConfig{
		Port: "8080",
	}

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value, ok := viper.Get("port").(string)
	if !ok {
		appConfig.Port = value
	}

	return appConfig
}
