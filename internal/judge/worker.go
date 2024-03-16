package judge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"msc24x/showdown/internal/mq"
	"msc24x/showdown/internal/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func logProcess(pid string, msg string, a ...any) {
	utils.LogWorker("Process ID %s: %s", pid, fmt.Sprintf(msg, a...))
}

// Responsible for executing one process
func processWorker(exe_proc *ExecutionProcess) {
	logProcess(exe_proc.PID, "recieved")
	start_time := time.Now().Unix()
	pid, err := uuid.Parse(exe_proc.PID)

	if err != nil {
		logProcess(exe_proc.PID, "failed, unable to parse PID %s", exe_proc.PID)
	}

	var exe_response ExecutionResponse
	exe_response.PID = exe_proc.PID
	exe_response.Webhook = exe_proc.Params.Webhook

	err = processRequest(
		pid,
		&exe_proc.Request,
		&exe_proc.Params,
		&exe_response,
	)

	exe_response.ServerFault = err != nil
	res_bytes, _ := json.Marshal(exe_response)
	res_text := string(res_bytes)

	content := bytes.NewBufferString(res_text)
	http.Post(exe_response.Webhook, "application/json", content)

	logProcess(exe_proc.PID, "completed. Took %d ms", time.Now().Unix()-start_time)
}

// Initiate the process queue consumption
func InitQueueWorker() {
	exe_proc_msgs := mq.Consume("executables")

	go func() {
		for exe_proc_msg := range exe_proc_msgs {
			var exec_obj ExecutionProcess
			err := json.Unmarshal(exe_proc_msg.Body, &exec_obj)

			if err != nil {
				utils.LogWorker("Unable to parse message due to %s\n%s (showing only 52 bytes)", err.Error(), exe_proc_msg.Body[:52])
				continue
			}

			go processWorker(&exec_obj)
		}
	}()

	utils.LogWorker("Worker initiated. Waiting for messages...")
}
