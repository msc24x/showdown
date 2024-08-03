// Package engine provides the methods to execute an execution request in an
// isolated and limited environment.
package engine

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"

	"github.com/msc24x/showdown/internal/config"
	"github.com/msc24x/showdown/internal/utils"

	"github.com/google/uuid"
)

// Constructs of fields required for a execution request to be valid.
type ExecutionRequest struct {
	// Code submitted by the user.
	Code string `json:"code"`
	// Programming language of the code.
	Language string `json:"language"`
	// What should be streamed into the program.
	Input string `json:"input"`
	// What is the expected output for the program.
	Output string `json:"output"`
}

func (exeReq *ExecutionRequest) Validate() error {

	flag, _ := IsSupportedLanguage(exeReq.Language)
	if !flag {
		return errors.New("language is not supported")
	}

	fields := reflect.ValueOf(exeReq).Elem()
	is_optional := func(name string) bool {
		optional_fields := []string{"Input", "Output"}

		for _, optional_field := range optional_fields {
			if optional_field == name {
				return true
			}
		}

		return false
	}

	for i := 0; i < fields.NumField(); i++ {

		field_val := fields.Field(i)
		field_name := fields.Type().Field(i)

		if field_val.String() == "" && !is_optional(field_name.Name) {
			return fmt.Errorf("%s is a required field", field_name.Tag)
		}
	}
	return nil

}

// BaseEngine is responsible for for isolated execution of the ExecutionRequest.
type BaseEngine struct {
	Request *ExecutionRequest

	// Is automatically assigned by the showdown when onboarding the execution request.
	PID uuid.UUID
	// Specifies which isolate box id/directory has been assigned for the request to reside in.
	isolateBoxID int
	// Language information fetched from ExecutionRequest.
	languageInfo *Language
	// Memory and compute limits during execution.
	limits       []*Limits
	activeLimits *Limits
	command      string
	commandArgs  []string
	envs         []string

	// Not all instances have same working directory, which is determined by the isolateBoxId.
	workDirectory  string
	sourceFilePath string
	inputFilePath  string
}

// Implements the public methods.
type BaseEnginePI interface {
	Init(exe_req *ExecutionRequest) error
	Execute() ([]byte, error)
	CollectMeta() ([]byte, error)
	Clean()
}

// Implements the Initialization helpers.
type BaseEngineInit interface {
	prepareFiles() error
	prepareCommand() error
}

// Implements the methods to work with Isolate binary.
type BaseEngineIsolate interface {
	prepareIsolatedBox(retry bool) error
	cleanIsolatedBox() error
	getIsolatedCommand() (*exec.Cmd, error)
}

// Parses the execution request and creates required files in the isolate box.
func (engine *BaseEngine) prepareFiles() error {
	file_name := fmt.Sprintf("%s.%s", engine.PID.String(), engine.Request.Language)
	engine.sourceFilePath = file_name
	engine.inputFilePath = engine.PID.String() + ".txt"

	absSourceFilePath := fmt.Sprintf("%s/%s", engine.workDirectory, engine.sourceFilePath)
	if err := os.WriteFile(absSourceFilePath, []byte(engine.Request.Code), 0644); err != nil {
		return utils.NewError(err, "Writing code file")
	}

	absInputFilePath := fmt.Sprintf("%s/%s", engine.workDirectory, engine.inputFilePath)
	if err := os.WriteFile(absInputFilePath, []byte(engine.Request.Input), 0644); err != nil {
		return utils.NewError(err, "Writing input file")
	}

	return nil
}

// Loads required envs and command args based on the language of the request.
func (engine *BaseEngine) prepareCommand() error {
	if engine.languageInfo.Compiled {
		engine.command = engine.languageInfo.CompilerPath
	} else {
		engine.command = engine.languageInfo.RunnerPath
	}

	engine.envs = append(engine.envs, engine.languageInfo.DefaultEnvs...)
	engine.commandArgs = append(engine.commandArgs, engine.languageInfo.SubCommands...)
	sourceFileSet := false

	for i, arg := range engine.commandArgs {
		if arg == CMD_FILE {
			engine.commandArgs[i] = engine.sourceFilePath
			sourceFileSet = true
		} else if arg == CMD_OUT {
			engine.commandArgs[i] = engine.PID.String()
		}
	}

	if !sourceFileSet {
		engine.commandArgs = append(engine.commandArgs, engine.sourceFilePath)
	}

	return nil
}

// Cleans up the isolate box directories.
func (engine *BaseEngine) cleanIsolatedBox() error {
	boxCleanupCmd := exec.Command(
		"sudo", config.ISOLATE_BIN, "--cg",
		"-b", fmt.Sprintf("%d", engine.isolateBoxID),
		"--cleanup",
	)

	_, err := boxCleanupCmd.CombinedOutput()

	return utils.NewError(err, "Cleaning isolate box")
}

// Creates or recreates, if already exists, isolate box.
func (engine *BaseEngine) prepareIsolatedBox(retry bool) error {
	boxInitCmd := exec.Command(
		"sudo", config.ISOLATE_BIN, "--cg",
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
	return utils.NewError(err, "Initializing isolate box")
}

// Reads the content of the meta file created during the execution of the request.
func (engine *BaseEngine) CollectMeta() ([]byte, error) {
	meta_file := fmt.Sprintf("%s/%s.info", engine.workDirectory, engine.PID.String())

	return os.ReadFile(meta_file)
}

// Isolates the 'command' and 'commandArgs' in the engine using Isolate binary,
// applies the compute and memory limits and returns a ready to execute exec.Command.
func (engine *BaseEngine) getIsolatedCommand(name string, args ...string) (*exec.Cmd, error) {
	isolate_args := []string{
		config.ISOLATE_BIN,
		"-b", fmt.Sprintf("%d", engine.isolateBoxID),
		"-p90",
		"--cg",
		"--cg-timing",
		"-x", "0",
		"--stderr-to-stdout",
		"-s",
		"-d", "/opt",
		"-M", fmt.Sprintf("%s/%s.info", engine.workDirectory, engine.PID.String()),
		"--open-files", "90",
		"-E", "HOME=/tmp",
		"-E", "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin"}

	if engine.activeLimits.Memory != -1 {
		isolate_args = append(isolate_args, "--cg-mem", fmt.Sprintf("%d", engine.activeLimits.Memory))
	}
	if engine.activeLimits.Stack != -1 {
		isolate_args = append(isolate_args, "-k", fmt.Sprintf("%d", engine.activeLimits.Stack))
	}
	if engine.activeLimits.Time != -1 {
		isolate_args = append(isolate_args, "-t", fmt.Sprintf("%f", engine.activeLimits.Time))
	}
	if engine.activeLimits.WallTime != -1 {
		isolate_args = append(isolate_args, "-w", fmt.Sprintf("%f", engine.activeLimits.WallTime))
	}

	isolate_args = append(isolate_args, "--run", "--", name)
	isolate_args = append(isolate_args, args...)
	isolate_cmd := exec.Command("sudo", isolate_args...)

	return isolate_cmd, nil
}

// Initializes an isolated box by Isolate binary, puts the required files in
// working directory created by the isolate and set the command and args
// to be executed. This must be executed in order to use BaseEngine.
func (engine *BaseEngine) Init(pid uuid.UUID, exe_req *ExecutionRequest, limits []*Limits) error {
	engine.Request = exe_req
	engine.PID = pid
	engine.limits = limits
	engine.languageInfo = GetLanguageInfo(&exe_req.Language)

	if engine.languageInfo == nil {
		return errors.New("cannot process given language")
	}

	box_id, ok := AssignBoxId()
	utils.BPanicIf(!ok, "Unable to to acquire an isolate box")

	engine.isolateBoxID = box_id
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

func (engine *BaseEngine) selectCompileLimits() *Limits {
	return engine.limits[0]
}

func (engine *BaseEngine) selectExecuteLimits() *Limits {
	if len(engine.limits) > 1 {
		return engine.limits[1]
	} else {
		return engine.selectCompileLimits()
	}
}

func (engine *BaseEngine) setLimits() {
	utils.BPanicIf(len(engine.limits) == 0, "Invalid limits")

	if engine.languageInfo.Compiled {
		engine.activeLimits = engine.selectCompileLimits()
	} else {
		engine.activeLimits = engine.selectExecuteLimits()
	}
}

// Works after the engine has been initialized. This method compiles the code
// if required and executes the submitted code. Returns the output bytes.
func (engine *BaseEngine) Execute() ([]byte, error) {
	engine.setLimits()

	isolated_cmd, _ := engine.getIsolatedCommand(
		engine.command,
		engine.commandArgs...,
	)

	input_file, err := os.Open(fmt.Sprintf("%s/%s", engine.workDirectory, engine.inputFilePath))

	if err != nil {
		return nil, err
	}

	defer input_file.Close()

	if !engine.languageInfo.Compiled {
		isolated_cmd.Stdin = input_file
	}

	output, err := isolated_cmd.CombinedOutput()

	if err != nil {
		return output, err
	}

	if engine.languageInfo.Compiled {
		engine.activeLimits = engine.selectExecuteLimits()
		isolated_exec, _ := engine.getIsolatedCommand(
			fmt.Sprintf("./%s", engine.PID.String()),
		)
		isolated_exec.Stdin = input_file

		return isolated_exec.CombinedOutput()
	}

	return output, err
}

// Clean redundant files created during the execution, and frees the allocated.
// isolate box.
func (engine *BaseEngine) Clean() {
	FreeBoxId(engine.isolateBoxID)

	files_to_remove := []string{
		fmt.Sprintf("%s/%s", engine.workDirectory, engine.sourceFilePath),
		fmt.Sprintf("%s/%s", engine.workDirectory, engine.PID.String()),
		fmt.Sprintf("%s/%s.txt", engine.workDirectory, engine.PID.String()),
		fmt.Sprintf("%s/%s.info", engine.workDirectory, engine.PID.String()),
	}

	for _, file := range files_to_remove {
		os.Remove(file)
	}
}
