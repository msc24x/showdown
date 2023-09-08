package main

import (
	"msc24x/showdown/api"
	"msc24x/showdown/config"

	"github.com/gin-gonic/gin"
)

func main() {

	if config.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	api.Init(router)

	router.Run("127.0.0.1:8080")

}
