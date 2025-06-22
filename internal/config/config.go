package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	DBUrl  string
	AppEnv string
	Port   string
}

func SetupConfig() (*Config, error) {
	// Configure viper
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("> Error reading .env file. If you are in development, please make sure that this exists")
		return nil, err
	}

	// Set/Create the db connection url string
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_NAME"),
	)

	log.Info().Msg("> ✅ viper configuration loaded correctly")

	cfg := Config{
		DBUrl:  dbUrl,
		AppEnv: viper.GetString("APP_ENV"),
		Port:   viper.GetString("PORT"),
	}

	return &cfg, nil
}
