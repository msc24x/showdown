package engine

import (
	"log"
	"strings"

	"github.com/msc24x/showdown/internal/config"

	"github.com/joho/godotenv"
)

type Language struct {
	// To temporarily disable a language's support.
	Supported bool
	// File extension of the programming language.
	Format string

	// Specifies if the language is compiled or not, whether to use a compiler
	// path of the runner path.
	Compiled     bool
	RunnerPath   string
	CompilerPath string

	DefaultEnvs []string
	SubCommands []string
}

const (
	// Placeholder identifier for the code file in the compilation command.
	CMD_FILE = "__SHDN_FILE"
	// Placeholder identifier for the output file in the compilation command.
	CMD_OUT = "__SHDN_OUT"
)

var (
	PYTHON = Language{
		Format:       config.PYTHON,
		Compiled:     false,
		RunnerPath:   "/opt/python/3.12.0/bin/python3",
		CompilerPath: "",
		Supported:    true,
	}
	CPP = Language{
		Format:       config.CPP,
		Compiled:     true,
		RunnerPath:   "",
		CompilerPath: "/usr/bin/g++",
		Supported:    true,
		SubCommands:  []string{CMD_FILE, "-o", CMD_OUT},
	}
	C = Language{
		Format:       config.C,
		Compiled:     true,
		RunnerPath:   "",
		CompilerPath: "/usr/bin/gcc",
		Supported:    true,
		SubCommands:  []string{CMD_FILE, "-o", CMD_OUT},
	}
	JAVASCRIPT = Language{
		Format:       config.JAVASCRIPT,
		Compiled:     false,
		RunnerPath:   "/usr/bin/node",
		CompilerPath: "",
		Supported:    true,
	}
	TYPESCRIPT = Language{
		Format:       config.TYPESCRIPT,
		Compiled:     false,
		RunnerPath:   "/usr/bin/ts-node",
		CompilerPath: "",
		Supported:    true,
		DefaultEnvs:  []string{"TS_NODE_FILES=true"},
	}
	GOLANG = Language{
		Format:       config.GOLANG,
		Compiled:     true,
		RunnerPath:   "",
		CompilerPath: "/usr/local/go/bin/go",
		Supported:    true,
		SubCommands:  []string{"build"},
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

// Import paths from file config.CONFIG_FILE.
func ImportPaths() {
	var (
		path string
		key  string
	)

	paths, err := godotenv.Read(config.CONFIG_FILE)

	if err != nil {
		log.Printf("Could not read file %s, continuing with default paths\n", config.CONFIG_FILE)
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
