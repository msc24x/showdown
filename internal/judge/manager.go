package judge

import (
	"log"
	"msc24x/showdown/config"
	"msc24x/showdown/internal/utils"
	"sync"
	"time"
)

type ShowdownWorker struct {
	InstanceId    int
	Address       string
	Authenticated bool
	Retries       int
	LastState     *InstanceState
}

var workers = []*ShowdownWorker{}
var Workers_mutex sync.Mutex
var worker_ticker = time.NewTicker(10 * time.Second)

func InitWorkersTicker() {
	go func() {
		for range worker_ticker.C {
			// log.Printf("Ticker running %s", t.String())
			PingWorkers()
		}
	}()
}

func PingWorkers() {

	for _, worker := range GetWorkers(true) {
		worker_url := worker.Address
		state, err := AuthenticateInstance(worker_url, config.T_WORKER)

		if err != nil {
			utils.LogWarn("%s: %s", worker_url, err.Error())

			if worker.Retries == 0 {
				worker.Authenticated = false
				utils.LogWarn("worker-%d dropped", worker.InstanceId)
			} else {
				worker.Retries--
				utils.LogWarn("worker-%d not responded, %d retries pending", worker.InstanceId, worker.Retries)
			}
			continue
		}

		if worker.Retries != config.MAX_WORKER_RETRIES {
			log.Printf("worker-%d revived", worker.InstanceId)
			worker.Retries = config.MAX_WORKER_RETRIES
		}

		worker.LastState = state
	}

}

func GetWorkers(authenticated bool) []*ShowdownWorker {
	if !authenticated {
		return workers
	}

	var auth_workers []*ShowdownWorker

	for _, worker := range workers {
		if worker.Authenticated {
			auth_workers = append(auth_workers, worker)
		}
	}

	return auth_workers
}

func AddWorker(instance *ShowdownWorker) {
	instance.Retries = config.MAX_WORKER_RETRIES
	workers = append(workers, instance)
	log.Printf("worker-%d authenticated", instance.InstanceId)
}

func GetMaxWorkerId() int {
	var max_id int = config.INSTANCE_ID
	for _, worker := range workers {
		if worker.InstanceId > max_id {
			max_id = worker.InstanceId
		}
	}

	return max_id
}
