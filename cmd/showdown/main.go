package main

import (
	"flag"
	"fmt"
	"log"
	"msc24x/showdown/api"
	"msc24x/showdown/config"
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
	flag.Parse()

	if fHelp {
		return A_HELP
	} else if fStart {
		return A_START
	}
	return A_DEFAULT
}

func setLogfile() *os.File {
	log.Println("Log file", config.LOG_FILE, "initializing...")

	log_file, err := os.Create(config.LOG_FILE)

	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(log_file)

	return log_file
}

func startServer() {
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
	printHeadline()
	parseEnv()

	if config.ENV == config.ENV_PROD {
		gin.SetMode(gin.ReleaseMode)
	}

	log_file := setLogfile()
	defer log_file.Close()

	action := parseFlags()

	switch action {
	case A_HELP:
		printHelp()
	case A_START:
		startServer()
	case A_DEFAULT:
		printHelp()
		os.Exit(-1)
	}

}
