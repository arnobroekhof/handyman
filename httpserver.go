package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Config struct
var Config = struct {
	HOST       string `default:"127.0.0.1:8080"`
	CONTEXT    string `default:"/"`
	USE_TOKENS bool   `default:"false"`

	COMMANDS []struct {
		Name    string
		Command string
		Arg     bool `default:"false"`
	}

	TOKENS []struct {
		Name  string
		Token string
	}
}{}

func initHTTPServer() {
	// initiate gin
	router := gin.New()
	// enable logging and recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// default unauthenticated ping check
	router.GET("/ping", getPing)

	// group commands in the same context
	commands := router.Group(Config.CONTEXT)

	if Config.USE_TOKENS {
		log.Println("Using tokens")
		commands.Use(tokenMiddleware)
	}

	// loop through the commands and configure the routes
	for _, command := range Config.COMMANDS {
		if command.Arg == true {
			addRouteWithArg(command.Name, command.Command, commands)
		} else if command.Arg == false {
			addRouteWithoutArg(command.Name, command.Command, commands)
		}
	}
	router.Run(Config.HOST)

}

// tokenMiddleware
func tokenMiddleware(c *gin.Context) {
	token := c.Request.Header.Get("X-Auth-Token")
	if token == "" {
		log.Println("No token provided")
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden", "reason": "X-Auth-Token not provided or empty"})
		c.Abort()
	} else {
		for _, t := range Config.TOKENS {
			if token == t.Token {
				log.Printf("%s authorized\n", t.Name)
				c.Next()
				return
			}
		}
		c.JSON(403, gin.H{"error": "unauthorized"})
		c.Abort()
	}
}
