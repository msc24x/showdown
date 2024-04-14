package judge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"msc24x/showdown/config"
	"msc24x/showdown/internal/mq"
	"msc24x/showdown/internal/utils"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type WorkerRegisteration struct {
	Address string
}

type WorkerRegisterationResponse struct {
	AssignedInstanceId int
}

func RegisterWorker(url string) {
	failIf := func(err error, context string) {
		if err != nil {
			utils.LogWorker("%s\nFailed to %s the manager on %s\n", err.Error(), context, url)
			os.Exit(1)
		}
	}

	stats, err := AuthenticateInstance(url, config.T_MANAGER)
	failIf(err, "ping")

	config.MANAGER_INSTANCE_ID = stats.InstanceId
	utils.LogWorker("Ping successful to manager instance %d running on %s", stats.InstanceId, url)

	req_body_struct := WorkerRegisteration{
		Address: fmt.Sprintf("%s://%s:%d", config.PROTOCOL, config.HOST, config.PORT),
	}

	req_body_bytes, err := json.Marshal(req_body_struct)
	utils.PanicIf(err)

	register_url := fmt.Sprintf("%s/workers/register", url)
	client := &http.Client{}
	req, err := http.NewRequest("POST", register_url, bytes.NewBuffer(req_body_bytes))
	failIf(err, "connect")

	req.Header.Set("Access-Token", config.ACCESS_TOKEN)
	res, err := client.Do(req)
	failIf(err, "connect")

	if res.StatusCode != 200 {
		res_bytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		utils.LogWorker("%s : Failed to connect to the manager on %s\n", string(res_bytes), url)
		os.Exit(1)
	}

	res_obj := WorkerRegisterationResponse{}
	err = json.NewDecoder(res.Body).Decode(&res_obj)
	failIf(err, "connect")

	config.INSTANCE_ID = res_obj.AssignedInstanceId

	utils.LogWorker("Connection with manager instance %d running on %s successful", stats.InstanceId, url)
}

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
