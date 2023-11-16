package engine

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
		Format:        "py",
		BuildRequired: false,
		RunnerPath:    "/usr/bin/python3",
		CompilerPath:  "",
		Supported:     true,
	}
	CPP = Language{
		Format:        "cpp",
		BuildRequired: false,
		RunnerPath:    "",
		CompilerPath:  "g++",
		Supported:     false,
	}
	C = Language{
		Format:        "c",
		BuildRequired: true,
		RunnerPath:    "",
		CompilerPath:  "gcc",
		Supported:     true,
	}
	JAVASCRIPT = Language{
		Format:        "js",
		BuildRequired: true,
		RunnerPath:    "node",
		CompilerPath:  "",
		Supported:     true,
	}
	TYPESCRIPT = Language{
		Format:        "ts",
		BuildRequired: true,
		RunnerPath:    "ts-node",
		CompilerPath:  "",
		Supported:     true,
		DefaultEnvs:   []string{"TS_NODE_FILES=true"},
	}
	GOLANG = Language{
		Format:        "go",
		BuildRequired: false,
		RunnerPath:    "/usr/local/go/bin/go",
		CompilerPath:  "/usr/local/go/bin/go",
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

func IsSupportedLanguage(lang string) (bool, *Language) {
	for _, value := range SUPPORTED_LANGUAGES {
		if value.Format == lang && value.Supported {
			return true, value
		}
	}
	return false, nil
}
