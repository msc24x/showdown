package engine

type Language struct {
	Supported bool
	Format    string

	HasRunner    bool
	RunnerPath   string
	CompilerPath string

	DefaultEnvs  []string
	SubCommands  []string
	DefaultFlags []string
}

var (
	PYTHON = Language{
		Format:       "py",
		HasRunner:    false,
		RunnerPath:   "python",
		CompilerPath: "",
		Supported:    true,
	}
	CPP = Language{
		Format:       "cpp",
		HasRunner:    true,
		RunnerPath:   "",
		CompilerPath: "g++",
		Supported:    true,
	}
	C = Language{
		Format:       "c",
		HasRunner:    true,
		RunnerPath:   "",
		CompilerPath: "gcc",
		Supported:    true,
	}
	JAVASCRIPT = Language{
		Format:       "js",
		HasRunner:    false,
		RunnerPath:   "node",
		CompilerPath: "",
		Supported:    true,
	}
	TYPESCRIPT = Language{
		Format:       "ts",
		HasRunner:    false,
		RunnerPath:   "ts-node",
		CompilerPath: "",
		Supported:    true,
		DefaultEnvs:  []string{"TS_NODE_FILES=true"},
	}
	GOLANG = Language{
		Format:       "go",
		HasRunner:    true,
		RunnerPath:   "go",
		CompilerPath: "",
		Supported:    true,
		SubCommands:  []string{"run"},
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
		if value.Format == lang {
			return true, value
		}
	}
	return false, nil
}
