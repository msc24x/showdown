// Package engine provides the methods to execute an execution request in an
// isolated and limited environment
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

// BaseEngine is responsible for for isolated execution of the ExecutionRequest
type BaseEngine struct {
	Request *ExecutionRequest

	isolateBoxID int
	languageInfo *Language // Language information fetched from ExecutionRequest
	limits       *Limits   // Memory and compute limits during execution
	command      string
	commandArgs  []string
	envs         []string

	workDirectory  string // Not all instances have same working directory, which is determined by the isolateBoxId
	sourceFilePath string
	inputFilePath  string
}

// Implements the public methods
type BaseEnginePI interface {
	Init(exe_req *ExecutionRequest) error
	Execute() ([]byte, error)
	CollectMeta() ([]byte, error)
	Clean()
}

// Implements the Initialization helpers
type BaseEngineInit interface {
	prepareFiles() error
	prepareCommand() error
}

// Implements the methods to work with Isolate binary
type BaseEngineIsolate interface {
	prepareIsolatedBox(retry bool) error
	cleanIsolatedBox() error
	getIsolatedCommand() (*exec.Cmd, error)
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

func (engine *BaseEngine) cleanIsolatedBox() error {
	boxCleanupCmd := exec.Command(
		config.ISOLATE_BIN, "--cg",
		"-b", fmt.Sprintf("%d", engine.isolateBoxID),
		"--cleanup",
	)
	_, err := boxCleanupCmd.CombinedOutput()
	return err
}

func (engine *BaseEngine) prepareIsolatedBox(retry bool) error {
	boxInitCmd := exec.Command(
		config.ISOLATE_BIN, "--cg",
		"-b", fmt.Sprintf("%d", engine.isolateBoxID),
		"--init",
	)
	_, err := boxInitCmd.CombinedOutput()
	if err != nil && retry {
		if err := engine.cleanIsolatedBox(); err != nil {
			return err
		}
		return engine.prepareIsolatedBox(false)
	}
	return err
}

// Reads the content of the meta file created during the execution of the request.
func (engine *BaseEngine) CollectMeta() ([]byte, error) {
	meta_file := fmt.Sprintf("%s/%s.info", engine.workDirectory, engine.Request.PID)

	return os.ReadFile(meta_file)
}

// Isolates the 'command' and 'commandArgs' in the engine using Isolate binary,
// applies the compute and memory limits and returns a ready to execute exec.Command.
func (engine *BaseEngine) getIsolatedCommand(name string, args ...string) (*exec.Cmd, error) {
	isolate_args := []string{
		"-b", fmt.Sprintf("%d", engine.isolateBoxID),
		"-p90",
		"--cg",
		"--cg-timing",
		"-x", "0",
		"--stderr-to-stdout",
		"-s",
		"-M", fmt.Sprintf("%s/%s.info", engine.workDirectory, engine.Request.PID),
		"--open-files", "90",
		"-E", "HOME=/tmp",
		"-E", "PATH=$PATH"}

	if engine.limits.Memory != -1 {
		isolate_args = append(isolate_args, "--cg-mem", fmt.Sprintf("%d", engine.limits.Memory))
	}
	if engine.limits.Stack != -1 {
		isolate_args = append(isolate_args, "-k", fmt.Sprintf("%d", engine.limits.Stack))
	}
	if engine.limits.Time != -1 {
		isolate_args = append(isolate_args, "-t", fmt.Sprintf("%f", engine.limits.Time))
	}
	if engine.limits.WallTime != -1 {
		isolate_args = append(isolate_args, "-w", fmt.Sprintf("%f", engine.limits.WallTime))
	}

	isolate_args = append(isolate_args, "--run", "--", name)
	isolate_args = append(isolate_args, args...)
	isolate_cmd := exec.Command(config.ISOLATE_BIN, isolate_args...)

	return isolate_cmd, nil
}

// Initializes an isolated box by Isolate binary, puts the required files in
// working directory created by the isolate and set the command and args
// to be executed. This must be executed in order to use BaseEngine.
func (engine *BaseEngine) Init(exe_req *ExecutionRequest) error {
	engine.Request = exe_req

	pid, _ := strconv.Atoi(exe_req.PID)
	engine.isolateBoxID = pid % config.MAX_ISOLATE_BOXES
	engine.workDirectory = fmt.Sprintf("%s/%d/box", config.ISOLATE_WORKDIR, engine.isolateBoxID)

	if err := engine.prepareIsolatedBox(true); err != nil {
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

// Works after the engine has been initialized. This method compiles the code
// if required and executes the submitted code. Returns the output bytes.
func (engine *BaseEngine) Execute() ([]byte, error) {
	if engine.languageInfo.BuildRequired {
		engine.limits = DEF_CMPL
	} else {
		engine.limits = DEF_EXEC
	}
	isolated_cmd, _ := engine.getIsolatedCommand(
		engine.command,
		engine.commandArgs...,
	)

	input_writer, _ := isolated_cmd.StdinPipe()
	input_writer.Write([]byte(engine.Request.Input))
	input_writer.Close()

	output, err := isolated_cmd.CombinedOutput()

	if err != nil {
		fmt.Println(string(output), err.Error())
		return output, err
	}

	if engine.languageInfo.BuildRequired {
		engine.limits = DEF_EXEC
		isolated_exec, _ := engine.getIsolatedCommand(
			fmt.Sprintf("./%s", engine.Request.PID),
		)

		return isolated_exec.CombinedOutput()
	}

	return output, err
}

// Clean redundant files created during the execution.
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
