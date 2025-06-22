package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	AppName          string
	DBUrl            string
	AppEnv           string
	AppDomain        string
	Port             string
	JWTSecret        string
	JWTAccessExpiry  int
	JWTRefreshExpiry int
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

	// get and validate required env values
	dbUser := mustGetString("DB_USER")
	dbPassword := mustGetString("DB_PASSWORD")
	dbName := mustGetString("DB_NAME")
	dbHost := mustGetString("DB_HOST")
	dbPort := mustGetString("DB_PORT")
	appEnv := mustGetString("APP_ENV")
	appPort := mustGetString("PORT")
	jwtSecret := mustGetString("JWT_SECRET")
	// in minutes
	jwtAccessExpiry := mustGetInt("JWT_ACCESS_EXPIRY")
	// in hours
	jwtRefreshExpiry := mustGetInt("JWT_REFRESH_EXPIRY")

	appName := mustGetString("APP_NAME")

	appDomain := mustGetString("APP_DOMAIN")

	// Set/Create the db connection url string
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	log.Info().Msg("> ✅ viper configuration loaded correctly")

	cfg := Config{
		AppName:          appName,
		DBUrl:            dbUrl,
		AppEnv:           appEnv,
		AppDomain:        appDomain,
		Port:             appPort,
		JWTSecret:        jwtSecret,
		JWTAccessExpiry:  jwtAccessExpiry,
		JWTRefreshExpiry: jwtRefreshExpiry,
	}
	fmt.Printf("config values: %v", cfg)

	return &cfg, nil
}

func mustGetString(key string) string {
	v := viper.GetString(key)

	if v == "" {
		log.Fatal().Msg(fmt.Sprintf("required config key '%s' is missing or empty", key))
	}

	return v
}

func mustGetInt(key string) int {
	if !viper.IsSet(key) {
		log.Fatal().Msg(fmt.Sprintf("required config key '%s' is missing or empty", key))
	}

	return viper.GetInt(key)
}
