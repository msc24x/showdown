package config

import "runtime"

const (
	DUMP_FILE      = "/var/lib/showdown/dump"
	ENGINE_WORKDIR = "/var/lib/showdown/files"
	LOG_FILE       = "/var/lib/showdown/log.log"

	ISOLATE_BIN     = "tools/isolate/bin/isolate"
	ISOLATE_WORKDIR = "/var/local/lib/isolate"

	ENV_DEV  = "dev"
	ENV_PROD = "prod"

	CPP        = "cpp"
	C          = "c"
	TYPESCRIPT = "ts"
	JAVASCRIPT = "js"
	PYTHON     = "py"
	GOLANG     = "go"

	T_STANDALONE = "standalone"
	T_MANAGER    = "manager"
	T_WORKER     = "worker"
)

var (
	// If set, requests to this instances would require the client to send it.
	ACCESS_TOKEN = ""
	// It is set automatically, however prone to changes if a worker re-authenticates.
	INSTANCE_ID uint = 1
	// It is set for workers, as its mandatory for them to connect to a manager.
	MANAGER_INSTANCE_ID      uint = 0
	MANAGER_INSTANCE_ADDRESS      = ""
	INSTANCE_TYPE                 = T_STANDALONE

	// After how many seconds should a manager ping workers to synchronize.
	ACTIVE_POLLING_RATE = 10 // Seconds
	// After how many seconds should a manager retry to ping dropped workers.
	REVIVAL_POLLING_RATE = 30 * 60 // Seconds
	// How many processes should as worker process at a given time at max.
	MAX_ACTIVE_PROCESSES uint = uint(runtime.NumCPU())
	// After how many tries to reach a stalled worker should it be marked as
	// dropped.
	MAX_WORKER_RETRIES uint8 = 6

	PROTOCOL    = "http"
	HOST        = "0.0.0.0"
	PORT        = 7070
	ENV         = ENV_PROD
	CONFIG_FILE = "env/.config"
	CREDS_FILE  = "env/.env.creds"

	RABBIT_MQ_HOST     = "localhost"
	RABBIT_MQ_PORT     = "5672"
	RABBIT_MQ_USER     = "guest"
	RABBIT_MQ_PASSWORD = "guest"
)
