package judge

import (
	"errors"
	"msc24x/showdown/config"
	"sync"
)

type JudgeState struct {
	TotalProcessed  int
	ActiveProcesses int
	Entries         map[string]bool
}

var judge_state = JudgeState{
	TotalProcessed:  0,
	ActiveProcesses: 0,
	Entries:         make(map[string]bool),
}
var judge_state_mutex sync.Mutex

func GetState() *JudgeState {
	return &judge_state
}

func RecordEntry(token string) error {
	judge_state_mutex.Lock()
	defer judge_state_mutex.Unlock()

	if judge_state.ActiveProcesses >= config.MAX_ACTIVE_PROCESSES {
		return errors.New("too many requests")
	}

	if judge_state.Entries[token] {
		return errors.New("duplicate request")
	}
	judge_state.ActiveProcesses++
	judge_state.Entries[token] = true
	return nil
}

func RecordProcessCompletion(token string) {
	judge_state_mutex.Lock()
	defer judge_state_mutex.Unlock()

	if !judge_state.Entries[token] {
		panic("recording a process completion that was not recorded")
	}
	judge_state.Entries[token] = false
	judge_state.ActiveProcesses--
	judge_state.TotalProcessed++
}
