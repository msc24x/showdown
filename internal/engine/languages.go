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
	BuildRequired bool
	RunnerPath    string
	CompilerPath  string

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
		BuildRequired: false,
		RunnerPath:    "node",
		CompilerPath:  "",
		Supported:     true,
	}
	TYPESCRIPT = Language{
		Format:        config.TYPESCRIPT,
		BuildRequired: false,
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
