package config

const (
	LOG_FILE             = "tmp/app.log"
	ENGINE_WORKDIR       = "files"
	MAX_ACTIVE_PROCESSES = 5

	ISOLATE_BIN       = "tools/isolate/bin/isolate"
	ISOLATE_WORKDIR   = "/var/local/lib/isolate"
	MAX_ISOLATE_BOXES = MAX_ACTIVE_PROCESSES

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

// Subject to change in response to process flags or envs
var (
	ACCESS_TOKEN  = ""
	INSTANCE_ID   = 0
	INSTANCE_TYPE = T_STANDALONE

	HOST       = "0.0.0.0"
	PORT       = 8080
	ENV        = ENV_DEV
	PATHS_FILE = ".env.paths"
	CREDS_FILE = ".env.creds"

	RABBIT_MQ_HOST     = "localhost"
	RABBIT_MQ_PORT     = "5672"
	RABBIT_MQ_USER     = "guest"
	RABBIT_MQ_PASSWORD = "guest"
)
