package engine

import (
	"log"
	"msc24x/showdown/config"
	"strings"

	"github.com/joho/godotenv"
)

type Language struct {
	Supported bool // To temporarily disable a language's support
	Format    string

	BuildRequired bool // Tells the BaseEngine whether to use RunnerPath or CompilerPath
	RunnerPath    string
	CompilerPath  string

	DefaultEnvs []string
	SubCommands []string
}

const (
	CMD_FILE = "__SHDN_FILE"
	CMD_OUT  = "__SHDN_OUT"
)

var (
	PYTHON = Language{
		Format:        config.PYTHON,
		BuildRequired: false,
		RunnerPath:    "python",
		CompilerPath:  "",
		Supported:     true,
	}
	CPP = Language{
		Format:        config.CPP,
		BuildRequired: true,
		RunnerPath:    "",
		CompilerPath:  "g++",
		Supported:     true,
		SubCommands:   []string{CMD_FILE, "-o", CMD_OUT},
	}
	C = Language{
		Format:        config.C,
		BuildRequired: true,
		RunnerPath:    "",
		CompilerPath:  "gcc",
		Supported:     true,
		SubCommands:   []string{CMD_FILE, "-o", CMD_OUT},
	}
	JAVASCRIPT = Language{
		Format:        config.JAVASCRIPT,
		BuildRequired: true,
		RunnerPath:    "node",
		CompilerPath:  "",
		Supported:     true,
	}
	TYPESCRIPT = Language{
		Format:        config.TYPESCRIPT,
		BuildRequired: true,
		RunnerPath:    "ts-node",
		CompilerPath:  "",
		Supported:     true,
		DefaultEnvs:   []string{"TS_NODE_FILES=true"},
	}
	GOLANG = Language{
		Format:        config.GOLANG,
		BuildRequired: true,
		RunnerPath:    "",
		CompilerPath:  "go",
		Supported:     true,
		SubCommands:   []string{"build"},
	}
)

var SUPPORTED_LANGUAGES = []*Language{
	&PYTHON,
	&CPP,
	&C,
	&JAVASCRIPT,
	&TYPESCRIPT,
	&GOLANG,
}

// Import paths from file config.CONFIG_FILE
func ImportPaths() {
	var (
		path string
		key  string
	)

	paths, err := godotenv.Read(config.CONFIG_FILE)

	if err != nil {
		log.Printf("Could not read file %s, continuing with defaults\n", config.CONFIG_FILE)
		return
	}

	for _, lang := range SUPPORTED_LANGUAGES {
		key = strings.ToUpper(lang.Format)
		path = paths[key]

		if path != "" {
			if lang.RunnerPath != "" {
				lang.RunnerPath = path
			}

			if lang.CompilerPath != "" {
				lang.CompilerPath = path
			}
		}
	}
}

func IsSupportedLanguage(lang string) (bool, *Language) {
	for _, value := range SUPPORTED_LANGUAGES {
		if value.Format == lang && value.Supported {
			return true, value
		}
	}
	return false, nil
}
