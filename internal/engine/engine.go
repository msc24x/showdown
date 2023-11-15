package engine

import (
	"errors"
	"fmt"
	"msc24x/showdown/config"
	"os"
	"os/exec"
	"reflect"
	"strconv"
)

type ExecutionRequest struct {
	PID      string `json:"pid"`
	UID      string `json:"uid"`
	Code     string `json:"code"`
	Language string `json:"language"`
	Input    string `json:"input"`
	Output   string `json:"output"`
}

func (exeReq *ExecutionRequest) Validate() error {

	flag, _ := IsSupportedLanguage(exeReq.Language)
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

	isolateBoxID int
	languageInfo *Language
	command      string
	commandArgs  []string
	envs         []string

	workDirectory  string
	sourceFilePath string
	inputFilePath  string
}

type BaseEngineInit interface {
	prepareFiles() error
	prepareCommand() error
	prepareIsolatedBox() error
}

type BaseEnginePI interface {
	Init(exe_req *ExecutionRequest) error
	getIsolatedCommand() (*exec.Cmd, error)
	CollectMeta() ([]byte, error)
	Execute() ([]byte, error)
	Clean()
}

func (engine *BaseEngine) prepareFiles() error {

	file_name := fmt.Sprintf("%s.%s", engine.Request.PID, engine.Request.Language)
	engine.sourceFilePath = file_name
	engine.inputFilePath = engine.Request.PID + ".txt"

	absSourceFilePath := fmt.Sprintf("%s/%s", engine.workDirectory, engine.sourceFilePath)
	if err := os.WriteFile(absSourceFilePath, []byte(engine.Request.Code), 0644); err != nil {
		return err
	}

	absInputFilePath := fmt.Sprintf("%s/%s", engine.workDirectory, engine.inputFilePath)
	if err := os.WriteFile(absInputFilePath, []byte(engine.Request.Input), 0644); err != nil {
		return err
	}

	return nil
}

func (engine *BaseEngine) prepareCommand() error {
	engine.languageInfo = GetLanguageInfo(&engine.Request.Language)
	if engine.languageInfo.BuildRequired {
		engine.command = engine.languageInfo.CompilerPath
	} else {
		engine.command = engine.languageInfo.RunnerPath
	}

	engine.envs = append(engine.envs, engine.languageInfo.DefaultEnvs...)
	engine.commandArgs = append(engine.commandArgs, engine.languageInfo.SubCommands...)
	engine.commandArgs = append(engine.commandArgs, engine.sourceFilePath)

	return nil
}

func (engine *BaseEngine) prepareIsolatedBox() error {
	boxInitCmd := exec.Command(
		config.ISOLATE_BIN, "--cg",
		"-b", fmt.Sprintf("%d", engine.isolateBoxID),
		"--init",
	)
	_, err := boxInitCmd.CombinedOutput()
	if err != nil && boxInitCmd.ProcessState.ExitCode() != 2 {
		return err
	}
	return nil
}

func (engine *BaseEngine) CollectMeta() ([]byte, error) {
	meta_file := fmt.Sprintf("%s/%s.info", engine.workDirectory, engine.Request.PID)

	return os.ReadFile(meta_file)
}

func (engine *BaseEngine) getIsolatedCommand(name string, args ...string) (*exec.Cmd, error) {
	isolate_args := []string{
		"-b", fmt.Sprintf("%d", engine.isolateBoxID),
		"-p90",
		"--cg",
		"--stderr-to-stdout",
		"-s",
		"-M", fmt.Sprintf("%s/%s.info", engine.workDirectory, engine.Request.PID),
		"--open-files", "90",
		"-E", "HOME=/tmp",
		"-E", "PATH=$PATH",
		"--run", "--", name}
	isolate_args = append(isolate_args, args...)
	isolate_cmd := exec.Command(config.ISOLATE_BIN, isolate_args...)

	return isolate_cmd, nil
}

func (engine *BaseEngine) Init(exe_req *ExecutionRequest) error {
	engine.Request = exe_req

	pid, _ := strconv.Atoi(exe_req.PID)
	engine.isolateBoxID = pid % config.MAX_ISOLATE_BOXES
	engine.workDirectory = fmt.Sprintf("%s/%d/box", config.ISOLATE_WORKDIR, engine.isolateBoxID)

	if err := engine.prepareIsolatedBox(); err != nil {
		return err
	}

	if err := engine.prepareFiles(); err != nil {
		return err
	}

	if err := engine.prepareCommand(); err != nil {
		return err
	}

	return nil
}

func (engine *BaseEngine) Execute() ([]byte, error) {
	isolated_cmd, _ := engine.getIsolatedCommand(
		engine.command,
		engine.commandArgs...,
	)

	input_writer, _ := isolated_cmd.StdinPipe()
	input_writer.Write([]byte(engine.Request.Input))
	input_writer.Close()

	output, err := isolated_cmd.CombinedOutput()

	if err != nil {
		return output, err
	}

	if engine.languageInfo.BuildRequired {
		isolated_exec, _ := engine.getIsolatedCommand(
			fmt.Sprintf("./%s", engine.Request.PID),
		)

		return isolated_exec.CombinedOutput()
	}

	return output, err
}

func (engine *BaseEngine) Clean() {
	files_to_remove := []string{
		fmt.Sprintf("%s/%s", engine.workDirectory, engine.sourceFilePath),
		fmt.Sprintf("%s/%s", engine.workDirectory, engine.Request.PID),
		fmt.Sprintf("%s/%s.txt", engine.workDirectory, engine.Request.PID),
		fmt.Sprintf("%s/%s.info", engine.workDirectory, engine.Request.PID),
	}

	for _, file := range files_to_remove {
		os.Remove(file)
	}
}
