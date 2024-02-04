package judge

import (
	"errors"
	"msc24x/showdown/config"
	"sync"
	"time"

	"github.com/google/uuid"
)

// An in memory runtime information and statistics of the Showdown
type JudgeState struct {
	Started time.Time

	TotalProcessed  int // Total number of requests processed since start
	ActiveProcesses int // Number of active requests being processed
	DeniedProcesses int // Number of denied requests to due limits or policies

	MaxUsers    int
	ActiveUsers int

	processes map[string]bool // Map of processes currently being processed
	users     map[string]int  // Map of users having at least one request being processed by the application.
}

var (
	judge_state = JudgeState{
		Started:         time.Now(),
		TotalProcessed:  0,
		ActiveProcesses: 0,
		DeniedProcesses: 0,
		MaxUsers:        0,
		ActiveUsers:     0,
		processes:       make(map[string]bool),
		users:           make(map[string]int),
	}
	judge_state_mutex sync.Mutex
)

func GetState() *JudgeState {
	return &judge_state
}

// Record/verify an execution request
func OnboardProcess(uid string) (uuid.UUID, error) {
	judge_state_mutex.Lock()
	defer judge_state_mutex.Unlock()

	if judge_state.ActiveProcesses >= config.MAX_ACTIVE_PROCESSES {
		judge_state.DeniedProcesses++
		return uuid.Nil, errors.New("too many requests")
	}

	pid := uuid.New()
	judge_state.ActiveProcesses++
	judge_state.processes[pid.String()] = true

	_, found_user := judge_state.users[uid]

	if found_user {
		judge_state.users[uid]++
	} else {
		judge_state.users[uid] = 1
	}

	judge_state.ActiveUsers++

	if judge_state.MaxUsers < judge_state.ActiveUsers {
		judge_state.MaxUsers = judge_state.ActiveUsers
	}

	return pid, nil
}

// Record the completion of an execution request
func OffboardProcess(pid uuid.UUID, uid string) {
	judge_state_mutex.Lock()
	defer judge_state_mutex.Unlock()

	if !judge_state.processes[pid.String()] {
		panic("recording a process completion that was not recorded")
	}
	delete(judge_state.processes, pid.String())
	judge_state.ActiveProcesses--
	judge_state.TotalProcessed++

	user_processes, found_user := judge_state.users[uid]

	if !found_user {
		panic("off boarding a user that was not recorded")
	}

	if user_processes == 1 {
		delete(judge_state.users, uid)
	} else {
		judge_state.users[uid]--
	}
	judge_state.ActiveUsers--

}
