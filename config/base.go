package config

const (
	LOG_FILE             = "tmp/app.log"
	PRODUCTION           = false
	ENGINE_WORKDIR       = "files"
	MAX_ACTIVE_PROCESSES = 5

	ISOLATE_BIN       = "tools/isolate/cli"
	ISOLATE_WORKDIR   = "/var/local/lib/isolate"
	MAX_ISOLATE_BOXES = MAX_ACTIVE_PROCESSES

	CPP        = "cpp"
	C          = "c"
	TYPESCRIPT = "ts"
	JAVASCRIPT = "js"
	PYTHON     = "py"
	GOLANG     = "go"
)
