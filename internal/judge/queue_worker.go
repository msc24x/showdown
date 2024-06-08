package judge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/msc24x/showdown/internal/config"
	"github.com/msc24x/showdown/internal/mq"
	"github.com/msc24x/showdown/internal/utils"

	"github.com/google/uuid"
)

type WorkerRegistration struct {
	Address string
}

type WorkerRegistrationResponse struct {
	AssignedInstanceId uint
}

func logProcess(pid string, fmsg string, a ...any) {
	utils.LogWorker("Process ID %s: %s", pid, fmt.Sprintf(fmsg, a...))
}

// Initiate the consumer to start executing processes
func InitQueueWorker() {
	exe_proc_msgs := mq.Consume("executables")

	go func() {
		for exe_proc_msg := range exe_proc_msgs {
			var exec_obj ExecutionProcess
			max_bytes := 52
			err := json.Unmarshal(exe_proc_msg.Body, &exec_obj)

			if err != nil {
				utils.LogWorker(
					"unable to parse message\n%s (showing only %d bytes)\ndue to %s\n",
					exe_proc_msg.Body[:max_bytes], max_bytes, err.Error(),
				)
				continue
			}

			pid, err := uuid.Parse(exec_obj.PID)
			if err != nil {
				utils.LogWorker("unable to parse PID %s", exec_obj.PID)
			}

			logProcess(exec_obj.PID, "received")
			capacity := OnboardProcess(pid)

			if capacity > 0 {
				go processWorker(&exec_obj, OffboardProcess)
			} else {
				processWorker(&exec_obj, OffboardProcess)
			}
		}
	}()

	utils.LogWorker("Worker initiated. Waiting for messages...")
}

// Executes one process
func processWorker(exe_proc *ExecutionProcess, offboard_c func(pid uuid.UUID)) {
	logProcess(exe_proc.PID, "started")
	start_time := time.Now().Unix()
	pid := uuid.MustParse(exe_proc.PID)
	defer func() {
		offboard_c(pid)
		logProcess(exe_proc.PID, "off boarded. Took %ds", time.Now().Unix()-start_time)
	}()

	var exe_response ExecutionResponse
	exe_response.PID = exe_proc.PID
	exe_response.Webhook = exe_proc.Params.Webhook

	err := processRequest(
		pid,
		&exe_proc.Request,
		&exe_proc.Params,
		&exe_response,
	)

	exe_response.ServerFault = err != nil
	res_bytes, _ := json.Marshal(exe_response)
	res_text := string(res_bytes)

	logProcess(exe_proc.PID, "processed. Took %ds", time.Now().Unix()-start_time)

	content := bytes.NewBufferString(res_text)
	webhook_req, err := http.NewRequest("POST", exe_response.Webhook, content)

	if err != nil {
		logProcess(exe_proc.PID, "unable create post request to the webhook '%s'", exe_response.Webhook)
		return
	}

	client := http.Client{}
	webhook_req.Header.Set("Webhook-Secret", config.WEBHOOK_SECRET)
	_, err = client.Do(webhook_req)

	if err != nil {
		logProcess(exe_proc.PID, "unable send post request to the webhook '%s'", exe_response.Webhook)
		return
	}
}
