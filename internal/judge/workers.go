package judge

import (
	"encoding/json"
	"log"
	"msc24x/showdown/config"
	"msc24x/showdown/internal/app"
	"msc24x/showdown/internal/utils"
	"sync"
	"time"
)

type WorkerStatus int8

const (
	SW_ACTIVE  WorkerStatus = 1 << 0
	SW_DROPPED WorkerStatus = 1 << 1
	SW_STALLED WorkerStatus = 1 << 2
)

type ShowdownWorker struct {
	InstanceId       uint
	Address          string
	Status           WorkerStatus
	Retries          uint8
	LastFetchedState *InstanceState

	CreatedSince  time.Time
	InactiveSince time.Time
}

var workers = []*ShowdownWorker{}
var Workers_mutex sync.Mutex

func InitWorkersTicker() {
	init := func(ticker *time.Ticker, status WorkerStatus) {
		go func() {
			for range ticker.C {
				PingWorkers(status)
			}
		}()
	}

	init(time.NewTicker(
		time.Duration(config.ACTIVE_POLLING_RATE)*time.Second),
		SW_ACTIVE|SW_STALLED,
	)

	init(time.NewTicker(
		time.Duration(config.REVIVAL_POLLING_RATE)*time.Second),
		SW_DROPPED,
	)
}

func PingWorkers(status WorkerStatus) {
	for _, worker := range GetWorkers(status) {
		worker_url := worker.Address
		state, err := AuthenticateInstance(worker_url, config.T_WORKER)

		if err != nil {
			utils.LogWarn("%s: %s", worker_url, err.Error())

			if worker.Retries == 0 {
				worker.Status = SW_DROPPED

				utils.LogWarn("worker-%d dropped", worker.InstanceId)
			} else {
				worker.Retries--
				worker.Status = SW_STALLED

				utils.LogWarn("worker-%d not responded, %d retries pending", worker.InstanceId, worker.Retries)
			}

			worker.InactiveSince = time.Now()
			continue
		}

		if state.InstanceId != worker.InstanceId && state.ManagerId == config.INSTANCE_ID {
			index, _ := GetWorkerByInstanceId(worker.InstanceId)
			utils.BPanicIf(index == -1, "unable to find worker entry just after ping")
			workers[index] = workers[len(workers)-1]
			workers[len(workers)-1] = nil
			workers = workers[:len(workers)-1]
			utils.LogWarn("worker-%d dropped permanently, due to invalid entry", worker.InstanceId)
			continue
		}

		if worker.Retries != config.MAX_WORKER_RETRIES {
			worker.Retries = config.MAX_WORKER_RETRIES
			worker.Status = SW_ACTIVE
			worker.InactiveSince = time.Time{}

			log.Printf("worker-%d revived", worker.InstanceId)
		}

		worker.LastFetchedState = state
	}

	res := GetInstanceState()
	state, _ := json.Marshal(res)
	app.DumpInstanceState(state)
}

func GetWorkers(status WorkerStatus) []*ShowdownWorker {
	var auth_workers []*ShowdownWorker

	for _, worker := range workers {
		if worker.Status&status == worker.Status {
			auth_workers = append(auth_workers, worker)
		}
	}

	return auth_workers
}

func AddWorker(instance *ShowdownWorker) {
	worker := GetWorkerByAddress(instance.Address)

	if worker == nil {
		instance.CreatedSince = time.Now()
		workers = append(workers, instance)
	} else {
		instance = worker
	}

	instance.Retries = config.MAX_WORKER_RETRIES
	instance.InactiveSince = time.Time{}
	instance.Status = SW_ACTIVE

	log.Printf("worker-%d authenticated", instance.InstanceId)
}

func GetMaxWorkerId() uint {
	var max_id uint = config.INSTANCE_ID

	for _, worker := range workers {
		if worker.InstanceId > max_id {
			max_id = worker.InstanceId
		}
	}

	return max_id
}

func GetWorkerByInstanceId(instanceId uint) (int, *ShowdownWorker) {
	for i, worker := range GetWorkers(SW_ACTIVE | SW_DROPPED | SW_STALLED) {
		if worker.InstanceId == instanceId {
			return i, worker
		}
	}

	return -1, nil
}

func GetWorkerByAddress(address string) *ShowdownWorker {
	for _, worker := range GetWorkers(SW_ACTIVE | SW_DROPPED | SW_STALLED) {
		if worker.Address == address {
			return worker
		}
	}

	return nil
}
