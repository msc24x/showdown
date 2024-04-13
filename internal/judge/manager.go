package judge

import (
	"log"
	"msc24x/showdown/config"
	"msc24x/showdown/internal/utils"
	"strings"
)

type ShowdownWorker struct {
	InstanceId    int
	Address       string
	Authenticated bool
	Last          bool
}

var workers = []*ShowdownWorker{}

func PingWorkers(workers_string *string) {
	worker_urls := strings.Split(*workers_string, ",")
	workers_authenticated := 0

	for _, worker_url := range worker_urls {
		stats, err := authenticateWorker(worker_url, config.T_WORKER)

		if err != nil {
			utils.LogWarn("%s: %s", worker_url, err.Error())
			workers = append(workers, &ShowdownWorker{Address: worker_url})
			continue
		}

		workers = append(workers, &ShowdownWorker{Address: worker_url, InstanceId: stats.InstanceId, Authenticated: true})
		workers_authenticated++
	}

	log.Printf("Total workers authenticated: %d", workers_authenticated)
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
