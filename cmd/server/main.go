package main

import (
	"blog-app/internal/config"
	"blog-app/internal/logger"
)

func main() {
	// Viper
	config.SetupConfig()
	// Zerolog
	logger.SetupLogger()
}
