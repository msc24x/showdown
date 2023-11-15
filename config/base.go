package config

const (
	PRODUCTION           = false
	ENGINE_WORKDIR       = "files"
	MAX_ACTIVE_PROCESSES = 5

	ISOLATE_BIN       = "internal/isolate/isolate"
	ISOLATE_WORKDIR   = "/var/local/lib/isolate"
	MAX_ISOLATE_BOXES = MAX_ACTIVE_PROCESSES
)

const (
	CPP        = "cpp"
	C          = "c"
	TYPESCRIPT = "ts"
	JAVASCRIPT = "js"
	PYTHON     = "py"
	GOLANG     = "go"
)
