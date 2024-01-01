package main

import (
	"flag"
	"fmt"
	"log"
	"msc24x/showdown/api"
	"msc24x/showdown/config"
	"msc24x/showdown/internal/engine"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	A_DEFAULT = 0
	A_HELP    = 1
	A_START   = 2
)
var (
	fHelp  bool
	fStart bool
	fPort  int
	fHost  string
	fPaths string
)

func parseEnv() {
	env := os.Getenv("ENV")

	if env != "" {
		config.ENV = env
	}

}

func parseFlags() int {
	flag.BoolVar(&fHelp, "help", false, "See usage")
	flag.BoolVar(&fStart, "start", false, "Start server to listen to requests")
	flag.IntVar(&fPort, "p", config.PORT, "Specify port to listen on")
	flag.StringVar(&fHost, "h", config.HOST, "Specify address to listen on")
	flag.StringVar(&fPaths, "paths", config.PATHS_FILE, "Specify .env file path to override defaults")
	flag.Parse()

	config.PATHS_FILE = fPaths
	engine.ImportPaths()

	if fHelp {
		return A_HELP
	} else if fStart {
		return A_START
	}
	return A_DEFAULT
}

func initLogs() *os.File {
	log.Println("Log file", config.LOG_FILE, "initializing...")

	log_file, err := os.Create(config.LOG_FILE)

	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(log_file)

	return log_file
}

func initServer() {
	router := gin.Default()
	api.AttachHandlers(router)
	address := fmt.Sprintf("%s:%d", fHost, fPort)

	log.Println("Listening on", address)
	router.Run(address)
}

func printHeadline() {
	fmt.Printf("Showdown! https://github.com/msc24x/showdown\nPortable server to execute and judge code.\n")
}

func printHelp() {
	printHeadline()
	fmt.Println("\nUsage:\n\t./showdown [...options] -start")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}

func main() {
	log_file := initLogs()
	parseEnv()

	if config.ENV == config.ENV_PROD {
		gin.SetMode(gin.ReleaseMode)
	}

	defer log_file.Close()

	action := parseFlags()

	switch action {
	case A_HELP:
		printHelp()
	case A_START:
		initServer()
	case A_DEFAULT:
		printHelp()
		os.Exit(-1)
	}

}
