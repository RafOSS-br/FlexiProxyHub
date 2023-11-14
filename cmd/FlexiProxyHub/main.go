package main

import (
	utils "FlexiProxyHub/internal/utils/generic"
	"FlexiProxyHub/internal/utils/logs"
	"fmt"
	"log"
	"os"

	"FlexiProxyHub/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	environmentMode := configureEnvironment()
	config := utils.CreateConfigFromEnv(logs.GetLogger())
	logs.Init(environmentMode, config.VisibleHeaders)
	logs.Log.Info(config.VisibleHeaders)
	log := logs.GetLogger()
	if config.LogLevel == "DEVELOPMENT" {
		logs.SetDebugMode()
	} else {
		logs.DisableDebugMode()
	}
	router.Router(config, log)
}

func configureEnvironment() string {
	// Check if dotenv exists
	var environment string
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		environment = os.Getenv("ENVIRONMENT")
		if environment == "DEVELOPMENT" {
			fmt.Println("Environment is DEVELOPMENT and .env file does not exist.")
			return "DEVELOPMENT"
		}
		return "PRODUCTION"
	} else {
		// Load .env
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		environment = os.Getenv("ENVIRONMENT")
		if environment == "DEVELOPMENT" {
			return "DEVELOPMENT"
		} else {
			return "PRODUCTION"
		}
	}
}
