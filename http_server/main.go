package http_server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
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

func healthGet(c *gin.Context) {
	c.String(200, "pong")
}

func addRouteWithArg(name string, cmd string, router *gin.RouterGroup) {
	fmt.Printf("Configuring route %v with command %s and argument\n", name, cmd)
	path := name + "/:arg"
	router.GET(path, func(c *gin.Context) {
		arg := c.Params.ByName("arg")
		if cmdOut, err := exec.Command(cmd, arg).Output(); err == nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "cmd": string(cmd), "argument": string(arg), "stdout": string(cmdOut)})
		}
	})
}

func addRouteWithoutArg(name string, cmd string, router *gin.RouterGroup) {
	fmt.Printf("Configuring route %s with command: %s\n", name, cmd)
	router.GET(name, func(c *gin.Context) {
		if cmdOut, err := exec.Command(cmd).Output(); err == nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "cmd": string(cmd), "stdout": string(cmdOut)})
		}
	})
}

func createCommands() {
	// initiate gin
	router := gin.New()
	// enable logging and recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	if Config.USE_TOKENS {
		log.Println("Using tokens")
		router.Use(tokenMiddleware)
	}

	//TODO: Create func for checking if the api key is set
	//and if set: check the headers and the key

	// health test
	router.GET("/ping", healthGet)

	// group commands in the same context
	commands := router.Group(Config.CONTEXT)

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
		c.Next()
	}
}

// Main function
func Main() {
	// load config file
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		log.Fatal("Environment variable CONFIG_FILE not set")
	}

	// load config file
	configor.Load(&Config, configFile)
	// create the commands
	createCommands()

}
