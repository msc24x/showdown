// Package judge utilizes the package engine and judges the execution output
package judge

import (
	"errors"
	"msc24x/showdown/config"
	"msc24x/showdown/internal/engine"
)

// Internal struct to conveniently pass execution response, params and pid
// from one process to another. Mainly for rabbit mq
type ExecutionProcess struct {
	PID     string                  `json:"pid"`
	Request engine.ExecutionRequest `json:"request"`
	Params  Params                  `json:"params"`
}

// Struc to define the end results the users will recieve
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
	// set the webhook, and showdown will send the response to that webhook
	// instead of an immediate response
	Webhook string `json:"webhook"`

	// set this to true and Showdown will only execute the code, not judge
	DoNotJudge bool `json:"donotjudge"`
}

func (params *Params) Validate() error {
	if params.Webhook == "" && config.INSTANCE_TYPE == config.T_MANAGER {
		return errors.New("webhook is mandatory for showdown manager instance")
	}

	return nil
}

// Entrypoint for an HTTP execution request
func JudgeExecutionRequest(exe_req *engine.ExecutionRequest, params *Params) (*ExecutionResponse, error) {

	pid, err := OnboardProcess()
	if err != nil {
		return nil, err
	}
	defer OffboardProcess(pid)

	response := ExecutionResponse{}
	response.PID = pid.String()
	response.Webhook = params.Webhook

	if params.Webhook != "" {
		queueRequest(pid, exe_req, params)
		return &response, nil
	}

	err = processRequest(pid, exe_req, params, &response)

	if err != nil {
		return &response, err
	}

	return &response, nil
}
