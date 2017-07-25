package main

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

// HealthGet check if the service is alive specific check
func getPing(c *gin.Context) {
	c.String(200, "pong")
}

func addRouteWithArg(name string, cmd string, router *gin.RouterGroup) {
	log.Printf("Configuring route %v with command %s and argument\n", name, cmd)
	path := name + "/:arg"
	router.GET(path, func(c *gin.Context) {
		arg := c.Params.ByName("arg")
		if cmdOut, err := exec.Command(cmd, arg).Output(); err == nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "cmd": string(cmd), "argument": string(arg), "stdout": string(cmdOut)})
		}
	})
}

func addRouteWithoutArg(name string, cmd string, router *gin.RouterGroup) {
	log.Printf("Configuring route %s with command: %s\n", name, cmd)
	router.GET(name, func(c *gin.Context) {
		if cmdOut, err := exec.Command(cmd).Output(); err == nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "cmd": string(cmd), "stdout": string(cmdOut)})
		}
	})
}
