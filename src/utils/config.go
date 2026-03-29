package utils

import (
	"log"
	"strings"

	"github.com/im-mk/user-service/src/models"
	"github.com/spf13/viper"
)

func GetConfig() models.ApplicationConfig {

	appConfig := models.ApplicationConfig{}
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

	if appConfig.Auth.TokenExpirySeconds == 0 {
		log.Printf("using default time of 600 seconds")
		appConfig.Auth.TokenExpirySeconds = 600
	}

	if appConfig.Auth.RefreshTokenExpirySeconds == 0 {
		log.Printf("using default time of 604800 seconds")
		appConfig.Auth.RefreshTokenExpirySeconds = 604800
	}

	return appConfig
}
