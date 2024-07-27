package engine

// Set of limits that can be applied on the BaseEngine before Execute.
type Limits struct {
	// CPU time allowed in seconds.
	Time float32 `json:"time"`
	// Total program execution time in seconds.
	WallTime float32 `json:"wall_time"`
	// Max memory limit in KBs.
	Memory int `json:"memory"`
	// Max stack limit in KBs, although part of the Memory.
	Stack int `json:"stack"`
}

var (
	DEF_CMPL = &Limits{
		Time:     10,
		WallTime: 15,
		Memory:   256 * 1024,
		Stack:    -1,
	}
	DEF_EXEC = &Limits{
		Time:     3,
		WallTime: 5,
		Memory:   256 * 1024,
		Stack:    -1,
	}
)
