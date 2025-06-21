package main

import (
	"blog-app/cmd/config"
	"fmt"
)

func main() {
	// Use viper to setup config
	config.SetupConfig()

	fmt.Println("Server is running!")
}
