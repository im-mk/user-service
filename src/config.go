package main

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

type AppConfig struct {
	Host string
	Port string
}

type ApplicationConfig struct {
	App    AppConfig
	DB     DBConfig
	JWTKey string
}

func GetConfig() ApplicationConfig {

	appConfig := ApplicationConfig{}
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()

	if err != nil {
		log.Printf("Error while reading config file %s", err)
	}

	configErr := viper.Unmarshal(&appConfig)
	if configErr != nil {
		log.Printf("Invalid configuration %s", configErr)
	}

	return appConfig
}
