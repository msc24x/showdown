package judge

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/msc24x/showdown/internal/app"
	"github.com/msc24x/showdown/internal/config"
	"github.com/msc24x/showdown/internal/engine"
	"github.com/msc24x/showdown/internal/utils"

	"github.com/google/uuid"
)

type InstanceConfig struct {
	ActivePollingRate  int
	RevivalPollingRate int
	MaxActiveProcesses uint
	MaxWorkerRetries   uint8

	Env          string
	ConfigFile   string
	CredsFile    string
	MessageQueue string

	SupportedLanguages []*engine.Language
}

// A general state struct for any showdown instance.
type InstanceState struct {
	InstanceId   uint
	InstanceType string
	// Provides manager id if instance type is worker.
	ManagerId              uint
	ManagerInstanceAddress string
	// Specifies if instance uses Access-Token.
	Private bool

	// Available instance configuration variables.
	Config *InstanceConfig

	// Not nil, if instance type is standalone/worker.
	WorkerState *WorkerState
	// Specifies connected workers if instance type is manager.
	Workers []*ShowdownWorker
}

// An in memory runtime information and statistics of the Showdown.
type WorkerState struct {
	StartedSince time.Time
	// Total number of requests processed since start.
	TotalProcessed uint
	// Number of active requests being processed.
	ActiveProcesses uint
	// Map of processes currently being processed.
	Processes map[string]bool
}

var (
	worker_state = WorkerState{
		StartedSince:    time.Now(),
		TotalProcessed:  0,
		ActiveProcesses: 0,
		Processes:       make(map[string]bool),
	}
	// Used to read/write protect worker_state.
	worker_state_mutex sync.Mutex
	// Used to protect worker_state.ActiveProcesses from exceeding beyond
	// MAX_ACTIVE_PROCESSES.
	onboarding_mutex sync.Mutex
)

// Reads the instance dump file and restores the workers data.
func RestoreManagerState() {
	worker_state_mutex.Lock()
	defer worker_state_mutex.Unlock()
	state_bytes, err := app.ReadInstanceState()

	if err != nil {
		utils.LogWarn("unable to read dump: %s", err.Error())
		return
	}

	instance_state := InstanceState{}
	err = json.Unmarshal(state_bytes, &instance_state)

	if err != nil {
		utils.LogWarn("unable to parse dump: %s", err.Error())
		return
	}

	workers = instance_state.Workers

	for _, worker := range workers {
		if worker.Status == SW_STALLED {
			worker.Retries = config.MAX_WORKER_RETRIES
			worker.Status = SW_ACTIVE
			worker.InactiveSince = time.Time{}
		}
	}

	log.Println("restored previous manager state")

	PingWorkers(SW_ACTIVE | SW_DROPPED | SW_STALLED)
}

func CollectInstanceConfig() *InstanceConfig {
	res := InstanceConfig{
		ActivePollingRate:  config.ACTIVE_POLLING_RATE,
		RevivalPollingRate: config.REVIVAL_POLLING_RATE,
		MaxActiveProcesses: config.MAX_ACTIVE_PROCESSES,
		MaxWorkerRetries:   config.MAX_WORKER_RETRIES,
		Env:                config.ENV,
		ConfigFile:         config.CONFIG_FILE,
		CredsFile:          config.CREDS_FILE,
		MessageQueue:       fmt.Sprintf("%s:%s", config.RABBIT_MQ_HOST, config.RABBIT_MQ_PORT),
		SupportedLanguages: engine.SUPPORTED_LANGUAGES,
	}

	return &res
}

func GetInstanceState(f_config bool) *InstanceState {
	res := InstanceState{
		InstanceId:             config.INSTANCE_ID,
		ManagerId:              config.MANAGER_INSTANCE_ID,
		InstanceType:           config.INSTANCE_TYPE,
		Private:                config.ACCESS_TOKEN != "",
		ManagerInstanceAddress: config.MANAGER_INSTANCE_ADDRESS,
	}

	if f_config {
		res.Config = CollectInstanceConfig()
	}

	if config.INSTANCE_TYPE != config.T_MANAGER {
		res.WorkerState = &worker_state
	} else if config.INSTANCE_TYPE != config.T_WORKER {
		res.Workers = GetWorkers(SW_ACTIVE | SW_DROPPED | SW_STALLED)
	}

	return &res
}

// Record/verify an execution request.
func OnboardProcess(pid uuid.UUID) (uint, error) {
	onboarding_mutex.Lock()

	worker_state_mutex.Lock()
	defer worker_state_mutex.Unlock()

	worker_state.ActiveProcesses++
	worker_state.Processes[pid.String()] = true

	if worker_state.ActiveProcesses < config.MAX_ACTIVE_PROCESSES {
		onboarding_mutex.Unlock()
	}

	return config.MAX_ACTIVE_PROCESSES - worker_state.ActiveProcesses, nil
}

// Record the completion of an execution request.
func OffboardProcess(pid uuid.UUID) {
	worker_state_mutex.Lock()
	defer worker_state_mutex.Unlock()

	if !worker_state.Processes[pid.String()] {
		panic("recording a process completion that was not recorded")
	}

	if worker_state.ActiveProcesses == config.MAX_ACTIVE_PROCESSES {
		onboarding_mutex.Unlock()
	}

	delete(worker_state.Processes, pid.String())
	worker_state.ActiveProcesses--
	worker_state.TotalProcessed++
}
