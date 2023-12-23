package engine

// Set of limits that can be applied on the BaseEngine before Execute
type Limits struct {
	Time     float32 // CPU time allowed in seconds
	WallTime float32 // Total program execution time in seconds
	Memory   int     // Max memory limit in KBs
	Stack    int     // Max stack limit in KBs, although part of the Memory
}

var (
	DEF_CMPL = &Limits{
		Time:     3,
		WallTime: 5,
		Memory:   20 * 1000,
		Stack:    -1,
	}
	DEF_EXEC = &Limits{
		Time:     10,
		WallTime: 15,
		Memory:   256 * 1000,
		Stack:    -1,
	}
)
