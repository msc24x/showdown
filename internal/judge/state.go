package judge

import (
	"errors"
	"msc24x/showdown/config"
	"sync"
	"time"
)

type JudgeState struct {
	Started         time.Time
	TotalProcessed  int
	ActiveProcesses int
	DeniedProcesses int
	entries         map[string]bool
}

var judge_state = JudgeState{
	Started:         time.Now(),
	TotalProcessed:  0,
	DeniedProcesses: 0,
	ActiveProcesses: 0,
	entries:         make(map[string]bool),
}
var judge_state_mutex sync.Mutex

func GetState() *JudgeState {
	return &judge_state
}

func OnboardProcess(token string) error {
	judge_state_mutex.Lock()
	defer judge_state_mutex.Unlock()

	if judge_state.ActiveProcesses >= config.MAX_ACTIVE_PROCESSES {
		judge_state.DeniedProcesses++
		return errors.New("too many requests")
	}

	if judge_state.entries[token] {
		judge_state.DeniedProcesses++
		return errors.New("duplicate request")
	}
	judge_state.ActiveProcesses++
	judge_state.entries[token] = true
	return nil
}

func OffboardProcess(token string) {
	judge_state_mutex.Lock()
	defer judge_state_mutex.Unlock()

	if !judge_state.entries[token] {
		panic("recording a process completion that was not recorded")
	}
	judge_state.entries[token] = false
	judge_state.ActiveProcesses--
	judge_state.TotalProcessed++
}
