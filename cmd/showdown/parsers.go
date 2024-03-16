package main

import (
	"flag"
	"msc24x/showdown/config"
	"msc24x/showdown/internal/engine"
	"msc24x/showdown/internal/utils"
	"os"

	"github.com/joho/godotenv"
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
	fCreds string
)

func parseEnv() {
	env := os.Getenv("ENV")

	if env != "" {
		config.ENV = env
	}
}

func parseCreds() {
	envMap, err := godotenv.Read(config.CREDS_FILE)
	utils.PanicIf(err)

	config.ACCESS_TOKEN = envMap["ACCESS_TOKEN"]
	config.RABBIT_MQ_HOST = envMap["RABBIT_MQ_HOST"]
	config.RABBIT_MQ_PORT = envMap["RABBIT_MQ_PORT"]
	config.RABBIT_MQ_USER = envMap["RABBIT_MQ_USER"]
	config.RABBIT_MQ_PASSWORD = envMap["RABBIT_MQ_PASSWORD"]
}

func parseFlags() int {
	flag.BoolVar(&fHelp, "help", false, "See usage")
	flag.BoolVar(&fStart, "start", false, "Start server to listen to requests")
	flag.IntVar(&fPort, "p", config.PORT, "Specify port to listen on")
	flag.StringVar(&fHost, "h", config.HOST, "Specify address to listen on")
	flag.StringVar(&fPaths, "paths", config.PATHS_FILE, "Specify .env file path to override defaults")
	flag.StringVar(&fCreds, "creds", config.PATHS_FILE, "Specify .env file path to override defaults")
	flag.Parse()

	config.CREDS_FILE = fCreds
	parseCreds()

	config.PATHS_FILE = fPaths
	engine.ImportPaths()

	if fHelp {
		return A_HELP
	} else if fStart {
		return A_START
	}
	return A_DEFAULT
}
