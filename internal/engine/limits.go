package engine

// Set of limits that can be applied on the BaseEngine before Execute
type Limits struct {
	Time     float32
	WallTime float32
	Memory   int
	Stack    int
}

var (
	DEF_CMPL = &Limits{
		Time:     3,
		WallTime: 5,
		Memory:   20 * 1000,
		Stack:    -1,
	}
	DEF_EXEC = &Limits{
		Time:     5,
		WallTime: 7,
		Memory:   128 * 1000,
		Stack:    -1,
	}
)
