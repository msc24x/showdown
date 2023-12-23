package config

const (
	LOG_FILE             = "tmp/app.log"
	ENGINE_WORKDIR       = "files"
	MAX_ACTIVE_PROCESSES = 5

	ISOLATE_BIN       = "tools/isolate/bin/isolate"
	ISOLATE_WORKDIR   = "/var/local/lib/isolate"
	MAX_ISOLATE_BOXES = MAX_ACTIVE_PROCESSES

	CPP        = "cpp"
	C          = "c"
	TYPESCRIPT = "ts"
	JAVASCRIPT = "js"
	PYTHON     = "py"
	GOLANG     = "go"

	ENV_DEV  = "dev"
	ENV_PROD = "prod"
)

// Subject to change in response to process flags or envs
var (
	HOST = "0.0.0.0"
	PORT = 8080
	ENV  = ENV_DEV
)
