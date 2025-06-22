package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func SetupConfig() error {
	// Configure viper
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("> Error reading .env file. If you are in development, please make sure that this exists")
		return err
	}

	log.Info().Msg("> ✅ viper configuration loaded correctly")

	return nil
}
