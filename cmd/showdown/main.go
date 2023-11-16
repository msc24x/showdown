package main

import (
	"log"
	"msc24x/showdown/api"
	"msc24x/showdown/config"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	log_file, err := os.Create(config.LOG_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer log_file.Close()
	log.SetOutput(log_file)

	if config.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	api.Init(router)

	router.Run("127.0.0.1:8080")

}
