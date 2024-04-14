package judge

import (
	"errors"
	"msc24x/showdown/config"
	"sync"
	"time"

	"github.com/google/uuid"
)

type InstanceState struct {
	InstanceId   int
	InstanceType string
	Private      bool

	JudgeStats *JudgeState
	Workers    []*ShowdownWorker
}

// An in memory runtime information and statistics of the Showdown
type JudgeState struct {
	Started time.Time

	TotalProcessed  int // Total number of requests processed since start
	ActiveProcesses int // Number of active requests being processed
	DeniedProcesses int // Number of denied requests to due limits or policies

	processes map[string]bool // Map of processes currently being processed
}

var (
	judge_state = JudgeState{
		Started:         time.Now(),
		TotalProcessed:  0,
		ActiveProcesses: 0,
		DeniedProcesses: 0,
		processes:       make(map[string]bool),
	}
	judge_state_mutex sync.Mutex
)

func GetState() *JudgeState {
	return &judge_state
}

// Record/verify an execution request
func OnboardProcess() (uuid.UUID, error) {
	judge_state_mutex.Lock()
	defer judge_state_mutex.Unlock()

	if judge_state.ActiveProcesses >= config.MAX_ACTIVE_PROCESSES {
		judge_state.DeniedProcesses++
		return uuid.Nil, errors.New("too many requests")
	}

	pid := uuid.New()
	judge_state.ActiveProcesses++
	judge_state.processes[pid.String()] = true

	return pid, nil
}

// Record the completion of an execution request
func OffboardProcess(pid uuid.UUID) {
	judge_state_mutex.Lock()
	defer judge_state_mutex.Unlock()

	if !judge_state.processes[pid.String()] {
		panic("recording a process completion that was not recorded")
	}
	delete(judge_state.processes, pid.String())
	judge_state.ActiveProcesses--
	judge_state.TotalProcessed++
}
