package engine

import (
	"sync"

	"github.com/msc24x/showdown/internal/config"
	"github.com/msc24x/showdown/internal/utils"
)

// Maintains an array of available isolate boxes to assign
// to execution requests.
type BoxesPool struct {
	// Array of size MAX_ACTIVE_PROCESSES, each item represents whether the
	// isolate box id is acquired or not.
	acquired []bool
	mutex    sync.Mutex
}

var boxes_pool *BoxesPool

// Allocates memory to initialize isolate boxes pool.
func AllocateBoxesPool() {
	utils.BPanicIf(boxes_pool != nil, "boxes pool is already allocated")

	boxes_pool = &BoxesPool{
		acquired: make([]bool, config.MAX_ACTIVE_PROCESSES),
		mutex:    sync.Mutex{},
	}
}

// Assigns one box id if available.
func AssignBoxId() (int, bool) {
	boxes_pool.mutex.Lock()
	defer boxes_pool.mutex.Unlock()

	for i := 0; i < int(config.MAX_ACTIVE_PROCESSES); i++ {
		if !boxes_pool.acquired[i] {
			boxes_pool.acquired[i] = true
			return i, true
		}
	}

	return -1, false
}

// Frees given box id, panics if already freed.
func FreeBoxId(box_id int) {
	boxes_pool.mutex.Lock()
	defer boxes_pool.mutex.Unlock()

	utils.BPanicIf(!boxes_pool.acquired[box_id],
		"Trying to free a free box: freeing %d", box_id,
	)

	boxes_pool.acquired[box_id] = false
}
