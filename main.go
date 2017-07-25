package main

import (
	"log"
	"os"

	"github.com/jinzhu/configor"
)

func main() {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		log.Fatal("Environment variable CONFIG_FILE not set")
	}

	// load config file
	configor.Load(&Config, configFile)
	// create the commands
	initHTTPServer()

}
