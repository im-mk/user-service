package main

import (
	"log"

	"github.com/spf13/viper"
)

type ApplicationConfig struct {
	Port string
}

func GetConfig() ApplicationConfig {

	appConfig := ApplicationConfig{
		Port: "8080",
	}

	viper.SetConfigType("json")
	viper.SetConfigFile(".env")

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
