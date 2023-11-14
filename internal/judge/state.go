package judge

import (
	"errors"
	"msc24x/showdown/config"
	"sync"
	"time"
)

type JudgeState struct {
	Started time.Time

	TotalProcessed  int
	ActiveProcesses int
	DeniedProcesses int

	MaxUsers    int
	ActiveUsers int

	processes map[string]bool
	users     map[string]int
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

func OnboardProcess(pid string, uid string) error {
	judge_state_mutex.Lock()
	defer judge_state_mutex.Unlock()

	if judge_state.ActiveProcesses >= config.MAX_ACTIVE_PROCESSES {
		judge_state.DeniedProcesses++
		return errors.New("too many requests")
	}

	if judge_state.processes[pid] {
		judge_state.DeniedProcesses++
		return errors.New("duplicate request")
	}
	judge_state.ActiveProcesses++
	judge_state.processes[pid] = true

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

	return nil
}

func OffboardProcess(pid string, uid string) {
	judge_state_mutex.Lock()
	defer judge_state_mutex.Unlock()

	if !judge_state.processes[pid] {
		panic("recording a process completion that was not recorded")
	}
	delete(judge_state.processes, pid)
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
