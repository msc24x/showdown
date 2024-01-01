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

	DefaultEnvs  []string
	SubCommands  []string
	DefaultFlags []string
}

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
		BuildRequired: false,
		RunnerPath:    "",
		CompilerPath:  "g++",
		Supported:     true,
	}
	C = Language{
		Format:        config.C,
		BuildRequired: true,
		RunnerPath:    "",
		CompilerPath:  "gcc",
		Supported:     true,
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
		BuildRequired: false,
		RunnerPath:    "go",
		CompilerPath:  "go",
		Supported:     true,
		SubCommands:   []string{"run"},
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

func ImportPaths() {
	var (
		path string
		key  string
	)

	paths, err := godotenv.Read(config.PATHS_FILE)

	if err != nil {
		log.Printf("Could not read file %s, continuing with defaults\n", config.PATHS_FILE)
		return
	}

	for _, lang := range SUPPORTED_LANGUAGES {
		key = strings.ToUpper(lang.Format)
		path = paths[key]

		if path != "" {
			if lang.RunnerPath != "" {
				lang.RunnerPath = path
				log.Println(lang.Format, "runner:\t", path)
			}

			if lang.CompilerPath != "" {
				lang.CompilerPath = path
				log.Println(lang.Format, "compiler:\t", path)
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
