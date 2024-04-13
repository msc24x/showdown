package main

import (
	"fmt"
	"io"
	"log"
	"msc24x/showdown/api"
	"msc24x/showdown/config"
	"msc24x/showdown/internal/judge"
	"msc24x/showdown/internal/mq"
	"msc24x/showdown/internal/utils"
	"os"

	"github.com/gin-gonic/gin"
)

func initLogs() *os.File {
	log.Println("Log file", config.LOG_FILE, "initializing...")

	log_file, err := os.Create(config.LOG_FILE)

	if err != nil {
		log.Fatal(err)
	}

	if config.ENV == config.ENV_PROD {
		log.SetOutput(log_file)
	} else {
		log_stream := io.MultiWriter(os.Stdout, log_file)
		log.SetOutput(log_stream)
	}

	log.SetPrefix(fmt.Sprintf("[SHOWDOWN-%s-%d] ", config.INSTANCE_TYPE, config.INSTANCE_ID))

	return log_file
}

func initServer() {
	log_file := initLogs()
	defer log_file.Close()

	router := gin.Default()

	if config.ACCESS_TOKEN == "" {
		utils.LogWarn("Connection not restricted, consider using ACCESS_TOKEN.")
	} else {
		router.Use(api.AccessToken())
	}

	api.AttachHandlers(router)

	if config.INSTANCE_TYPE != config.T_MANAGER {
		judge.PingManager(config.MANAGER_INSTANCE_ADDRESS)
		closeConnection := mq.InitMessageQueue()
		defer closeConnection()
		judge.InitQueueWorker()
	}

	address := fmt.Sprintf("%s:%d", fHost, fPort)
	log.Printf("Started Showdown %s-%d on %s", config.INSTANCE_TYPE, config.INSTANCE_ID, address)
	router.Run(address)

}

func main() {
	parseEnv()

	if config.ENV == config.ENV_PROD {
		gin.SetMode(gin.ReleaseMode)
	}

	action := parseFlags()

	switch action {
	case A_HELP:
		printHelp()
	case A_START:
		initServer()
	case A_DEFAULT:
		printHelp()
		os.Exit(1)
	}

}
