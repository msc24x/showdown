// Package judge utilizes the package engine and judges the execution output.
package judge

import (
	"encoding/json"
	"errors"

	"github.com/msc24x/showdown/internal/config"
	"github.com/msc24x/showdown/internal/engine"
	"github.com/msc24x/showdown/internal/mq"

	"github.com/google/uuid"
)

// Internal struct to conveniently pass execution response, params and pid
// from one process to another. Mainly for rabbit mq.
type ExecutionProcess struct {
	PID     string                  `json:"pid"`
	Request engine.ExecutionRequest `json:"request"`
	Params  Params                  `json:"params"`
}

// Struct to define the end results the users will receive.
type ExecutionResponse struct {
	PID         string `json:"pid"`
	Webhook     string `json:"webhook"`
	Success     bool   `json:"success"`
	Judged      bool   `json:"judged"`
	Error       string `json:"error"`
	Output      string `json:"output"`
	Meta        string `json:"meta"`
	Expected    string `json:"expected"`
	ServerFault bool   `json:"server_fault"`
}

type Params struct {
	// Set the webhook, and showdown will send the response to that webhook
	// instead of an immediate response.
	Webhook string `json:"webhook"`

	// Set this to true and Showdown will only execute the code, not judge.
	DoNotJudge bool `json:"donotjudge"`
}

func (params *Params) Validate() error {
	if params.Webhook == "" && config.INSTANCE_TYPE == config.T_MANAGER {
		return errors.New("webhook is mandatory for showdown manager instance")
	}

	return nil
}

// Entrypoint for an HTTP execution request.
func JudgeExecutionRequest(exe_req *engine.ExecutionRequest, params *Params) (*ExecutionResponse, error) {
	pid := uuid.New()
	response := ExecutionResponse{}
	response.PID = pid.String()
	response.Webhook = params.Webhook

	if params.Webhook != "" {
		queueRequest(pid, exe_req, params)
		return &response, nil
	}

	err := processRequest(pid, exe_req, params, &response)

	if err != nil {
		return &response, err
	}

	return &response, nil
}

// Takes the execution request and queues it into rabbit mq queue.
func queueRequest(pid uuid.UUID, exe_req *engine.ExecutionRequest, params *Params) {
	process_obj := ExecutionProcess{
		PID:     pid.String(),
		Request: *exe_req,
		Params:  *params,
	}
	process_body, _ := json.Marshal(process_obj)
	go mq.Queue("executables", 3, process_body)
}
