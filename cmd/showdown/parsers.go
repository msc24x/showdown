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
	A_MANAGER = 3
	A_WORKER  = 4
)
var (
	fHelp           bool
	fStart          bool
	fManager        bool
	fWorker         bool
	fManagerAddress string
	fPort           int
	fHost           string
	fPaths          string
	fCreds          string
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
	config.WEBHOOK_SECRET = envMap["WEBHOOK_SECRET"]
	config.RABBIT_MQ_HOST = envMap["RABBIT_MQ_HOST"]
	config.RABBIT_MQ_PORT = envMap["RABBIT_MQ_PORT"]
	config.RABBIT_MQ_USER = envMap["RABBIT_MQ_USER"]
	config.RABBIT_MQ_PASSWORD = envMap["RABBIT_MQ_PASSWORD"]
}

func parseFlags() int {
	flag.BoolVar(&fHelp, "help", false, "See usage")
	flag.BoolVar(&fStart, "start", false, "Start showdown")
	flag.BoolVar(&fManager, "m", false, "Start showdown as a manager instance")
	flag.BoolVar(&fWorker, "w", false, "Start showdown as a worker instance")
	flag.StringVar(&fManagerAddress, "c", "http://localhost:8080", "Provide manager instance address (-w required)")
	flag.IntVar(&fPort, "p", config.PORT, "Specify port to listen on")
	flag.StringVar(&fHost, "h", config.HOST, "Specify address to listen on")
	flag.StringVar(&fPaths, "paths", config.PATHS_FILE, "Specify .env file path to override defaults")
	flag.StringVar(&fCreds, "creds", config.PATHS_FILE, "Specify .env file path to override defaults for standalone or worker instance")
	flag.Parse()

	if fWorker {
		config.INSTANCE_TYPE = config.T_WORKER
		config.MANAGER_INSTANCE_ADDRESS = fManagerAddress
	} else if fManager {
		config.INSTANCE_TYPE = config.T_MANAGER
	}

	config.PORT = fPort
	config.HOST = fHost

	config.CREDS_FILE = fCreds
	parseCreds()

	if config.INSTANCE_TYPE != config.T_MANAGER {
		config.PATHS_FILE = fPaths
		engine.ImportPaths()
	}

	if fHelp {
		return A_HELP
	} else if fStart {
		return A_START
	}
	return A_DEFAULT
}
