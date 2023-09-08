package config

const (
	PRODUCTION           = false
	ENGINE_WORKDIR       = "files"
	MAX_ACTIVE_PROCESSES = 5
)

const (
	CPP        = "cpp"
	C          = "c"
	TYPESCRIPT = "ts"
	JAVASCRIPT = "js"
	PYTHON     = "py"
	GOLANG     = "go"
)

var SUPPORTED_LANGUAGES = []string{
	CPP,
	C,
	TYPESCRIPT,
	JAVASCRIPT,
	PYTHON,
	GOLANG,
}

var RUNNER_PATHS = map[string]string{
	CPP:        "g++",
	C:          "gcc",
	JAVASCRIPT: "node",
	TYPESCRIPT: "ts-node",
	PYTHON:     "python",
	GOLANG:     "go",
}

func GetRunnerPath(lang *string) string {
	return RUNNER_PATHS[*lang]
}
