package http_server

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
)

// Config struct
var Config = struct {
	HOST    string `default:"127.0.0.1:8080"`
	CONTEXT string `default:"/"`
	//TODO: USE_TOKEN    bool   `default:"false"`
	//TODO: TOKEN_SECRET string `default:"12345678910"`

	COMMANDS []struct {
		Name    string
		Command string
		Arg     bool `default:"false"`
	}
}{}

func healthGet(c *gin.Context) {
	c.String(200, "pong")
}

func createCommands() {
	// initiate gin
	router := gin.New()
	// enable logging and recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	//TODO: Create func for checking if the api key is set
	//and if set: check the headers and the key

	// health test
	router.GET("/ping", healthGet)

	// group commands in the same context
	commands := router.Group(Config.CONTEXT)

	for _, command := range Config.COMMANDS {

		// loop through the command struct and configure the routes
		if command.Arg {
			path := command.Name + "/:arg"
			commands.GET(path, func(c *gin.Context) {
				arg := c.Params.ByName("arg")
				if arg != "" {
					cmd := command.Command
					println("executing command: ", cmd, arg)
					if cmdOut, err := exec.Command(cmd, arg).Output(); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"std.error": string(cmdOut), "error": err, "status": err})
					} else {
						c.JSON(http.StatusOK, gin.H{"status": "ok", "output": string(cmdOut)})
					}
				} else {
					c.String(500, "failed")
				}
			})
		} else {
			path := command.Name
			commands.PUT(path, func(c *gin.Context) {
				cmd := command.Command
				println("executing command: ", cmd)
				if cmdOut, err := exec.Command(cmd).Output(); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"std.error": string(cmdOut), "error": err})
				} else {
					c.JSON(http.StatusOK, gin.H{"status": "ok", "cmd": string(cmdOut)})
				}

			})
		}

	}
	router.Run(Config.HOST)

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
