package config

const (
	DUMP_FILE            = "/var/lib/showdown/dump"
	ENGINE_WORKDIR       = "/var/lib/showdown/files"
	LOG_FILE             = "/var/lib/showdown/log.log"
	MAX_ACTIVE_PROCESSES = 5
	MAX_WORKER_RETRIES   = 6

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
	ACCESS_TOKEN                  = ""
	WEBHOOK_SECRET                = ""
	INSTANCE_ID              uint = 1
	MANAGER_INSTANCE_ID      uint = 0
	MANAGER_INSTANCE_ADDRESS      = ""
	INSTANCE_TYPE                 = T_STANDALONE

	ACTIVE_POLLING_RATE  = 10      // Seconds
	REVIVAL_POLLING_RATE = 30 * 60 // Seconds

	PROTOCOL    = "http"
	HOST        = "0.0.0.0"
	PORT        = 7070
	ENV         = ENV_DEV
	CONFIG_FILE = "env/.config"
	CREDS_FILE  = "env/.env.creds"

	RABBIT_MQ_HOST     = "localhost"
	RABBIT_MQ_PORT     = "5672"
	RABBIT_MQ_USER     = "guest"
	RABBIT_MQ_PASSWORD = "guest"
)
