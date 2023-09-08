package engine

import (
	"errors"
	"fmt"
	"msc24x/showdown/config"
	"os"
	"os/exec"
	"reflect"
	"time"
)

type ExecutionRequest struct {
	Token    string `json:"token"`
	Code     string `json:"code"`
	Language string `json:"language"`
	Input    string `json:"input"`
	Output   string `json:"output"`
}

func IsSupportedLanguage(lang string) bool {
	for _, value := range config.SUPPORTED_LANGUAGES {
		if value == lang {
			return true
		}
	}
	return false
}

func (exeReq *ExecutionRequest) Validate() error {

	flag := IsSupportedLanguage(exeReq.Language)
	if !flag {
		return errors.New("language is not supported")
	}

	fields := reflect.ValueOf(exeReq).Elem()

	for i := 0; i < fields.NumField(); i++ {

		field_val := fields.Field(i)
		field_name := fields.Type().Field(i).Tag

		if field_val.String() == "" {
			return fmt.Errorf("%s is a required field", field_name)
		}
	}
	return nil

}

type BaseEngine struct {
	Request *ExecutionRequest

	runnerPath     string
	runnerCommands []string
	runnerEnviron  []string

	workDirectory  string
	sourceFilePath string
	inputFilePath  string
}

type BaseEnginePreparer interface {
	prepareFiles() error
	prepareCommand() error
}

type BaseEnginePI interface {
	Init(exe_req *ExecutionRequest) error
	Execute() ([]byte, error)
	Clean()
}

func (engine *BaseEngine) prepareFiles() error {
	engine.workDirectory = config.ENGINE_WORKDIR

	file_name := fmt.Sprintf("%s.%s", engine.Request.Token, engine.Request.Language)
	engine.sourceFilePath = fmt.Sprintf("%s/%s", engine.workDirectory, file_name)
	engine.inputFilePath = fmt.Sprintf("%s/%s.txt", engine.workDirectory, engine.Request.Token)

	err := os.WriteFile(engine.sourceFilePath, []byte(engine.Request.Code), 0644)

	if err != nil {
		return err
	}

	err = os.WriteFile(engine.inputFilePath, []byte(engine.Request.Input), 0644)

	if err != nil {
		return err
	}

	return nil
}

func (engine *BaseEngine) prepareCommand() error {
	engine.runnerPath = config.GetRunnerPath(&engine.Request.Language)

	switch engine.Request.Language {
	case config.TYPESCRIPT:
		engine.runnerEnviron = append(engine.runnerEnviron, "TS_NODE_FILES=true")
	case config.GOLANG:
		engine.runnerCommands = append(engine.runnerCommands, "run")
	}

	engine.runnerCommands = append(engine.runnerCommands, engine.sourceFilePath)

	return nil
}

func (engine *BaseEngine) Init(exe_req *ExecutionRequest) error {
	engine.Request = exe_req

	err := engine.prepareFiles()
	if err != nil {
		return err
	}

	err = engine.prepareCommand()
	if err != nil {
		return err
	}

	return nil
}

func (engine *BaseEngine) Execute() ([]byte, error) {

	exe := exec.Command(engine.runnerPath, engine.runnerCommands...)
	timeout := false

	exe.Env = append(os.Environ(), engine.runnerEnviron...)

	input_writer, _ := exe.StdinPipe()

	input_writer.Write([]byte(engine.Request.Input))
	input_writer.Close()

	timer := time.AfterFunc(time.Second*time.Duration(5), func() {
		exe.Process.Kill()
		timeout = true

	})
	defer timer.Stop()
	output, _ := exe.CombinedOutput()

	var err error
	if timeout {
		err = errors.New("timeout")
	}

	return output, err
}

func (engine *BaseEngine) Clean() {

	os.Remove(engine.inputFilePath)
	os.Remove(engine.sourceFilePath)
}
