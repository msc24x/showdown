package judge

import (
	"encoding/json"
	"log"
	"msc24x/showdown/config"
	"msc24x/showdown/internal/app"
	"msc24x/showdown/internal/utils"
	"sync"
	"time"

	"github.com/google/uuid"
)

type InstanceState struct {
	InstanceId   uint
	ManagerId    uint
	InstanceType string
	Private      bool

	WorkerState *WorkerState
	Workers     []*ShowdownWorker
}

// An in memory runtime information and statistics of the Showdown
type WorkerState struct {
	StartedSince time.Time

	TotalProcessed  uint // Total number of requests processed since start
	ActiveProcesses uint // Number of active requests being processed

	Processes map[string]bool // Map of processes currently being processed
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

// Record/verify an execution request
func OnboardProcess(pid uuid.UUID) uint {
	worker_state_mutex.Lock()
	defer worker_state_mutex.Unlock()

	worker_state.ActiveProcesses++
	worker_state.Processes[pid.String()] = true

	return config.MAX_ACTIVE_PROCESSES - worker_state.ActiveProcesses
}

// Record the completion of an execution request
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
