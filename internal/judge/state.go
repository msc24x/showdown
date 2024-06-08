package judge

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/msc24x/showdown/internal/app"
	"github.com/msc24x/showdown/internal/config"
	"github.com/msc24x/showdown/internal/utils"

	"github.com/google/uuid"
)

// A general state struct for any showdown instance.
type InstanceState struct {
	InstanceId   uint
	InstanceType string
	// Provides manager id if instance type is worker.
	ManagerId uint
	// Specifies if instance uses Access-Token.
	Private bool

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
	worker_state_mutex sync.Mutex
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

func GetInstanceState() *InstanceState {
	res := InstanceState{
		InstanceId:   config.INSTANCE_ID,
		ManagerId:    config.MANAGER_INSTANCE_ID,
		InstanceType: config.INSTANCE_TYPE,
		Private:      config.ACCESS_TOKEN != "",
	}

	if config.INSTANCE_TYPE != config.T_MANAGER {
		res.WorkerState = &worker_state
	} else if config.INSTANCE_TYPE != config.T_WORKER {
		res.Workers = GetWorkers(SW_ACTIVE | SW_DROPPED | SW_STALLED)
	}

	return &res
}

// Record/verify an execution request.
func OnboardProcess(pid uuid.UUID) uint {
	worker_state_mutex.Lock()
	defer worker_state_mutex.Unlock()

	worker_state.ActiveProcesses++
	worker_state.Processes[pid.String()] = true

	return config.MAX_ACTIVE_PROCESSES - worker_state.ActiveProcesses
}

// Record the completion of an execution request.
func OffboardProcess(pid uuid.UUID) {
	worker_state_mutex.Lock()
	defer worker_state_mutex.Unlock()

	if !worker_state.Processes[pid.String()] {
		panic("recording a process completion that was not recorded")
	}

	delete(worker_state.Processes, pid.String())
	worker_state.ActiveProcesses--
	worker_state.TotalProcessed++
}
