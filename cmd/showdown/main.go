package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/msc24x/showdown/internal/api"
	"github.com/msc24x/showdown/internal/config"
	"github.com/msc24x/showdown/internal/judge"
	"github.com/msc24x/showdown/internal/mq"
	"github.com/msc24x/showdown/internal/utils"

	"github.com/gin-gonic/gin"
)

func initLogs() *os.File {
	log.Println("Attaching logger at", config.LOG_FILE)

	log_file, err := os.OpenFile(config.LOG_FILE, os.O_APPEND, 0666)

	if err != nil {
		log_file, err = os.Create(config.LOG_FILE)

		if err != nil {
			log.Fatal(err)
		}
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
		utils.LogWarn("Connection not restricted, consider using ACCESS_TOKEN")

		if config.INSTANCE_TYPE != config.T_STANDALONE {
			log.Fatalln("Role based instances must use ACCESS_TOKEN")
		}

	} else {
		router.Use(api.AccessToken())
	}

	api.AttachHandlers(router)

	closeConnection := mq.InitMessageQueue()
	defer closeConnection()

	if config.INSTANCE_TYPE != config.T_MANAGER {
		judge.InitQueueWorker()
	}

	if config.INSTANCE_TYPE == config.T_MANAGER {
		judge.InitWorkersTicker()
		judge.RestoreManagerState()
	}

	address := fmt.Sprintf("%s:%d", config.HOST, config.PORT)
	log.Printf("Started Showdown %s-%d on %s", config.INSTANCE_TYPE, config.INSTANCE_ID, address)
	go func() {
		if config.INSTANCE_TYPE == config.T_WORKER {
			judge.ConnectManager(config.MANAGER_INSTANCE_ADDRESS)
			log.SetPrefix(fmt.Sprintf("[SHOWDOWN-%s-%d] ", config.INSTANCE_TYPE, config.INSTANCE_ID))
		}
	}()

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
