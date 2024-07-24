package api

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/msc24x/showdown/internal/config"
	"github.com/msc24x/showdown/internal/engine"
	"github.com/msc24x/showdown/internal/judge"

	"github.com/gin-gonic/gin"
)

// Struct to define http request to showdown
type JudgeRequest struct {
	JudgeParams judge.Params            `json:"judge_params"`
	Exe         engine.ExecutionRequest `json:"exe"`
}

func DebugWebhook(c *gin.Context) {
	if b, err := io.ReadAll(c.Request.Body); err == nil {
		fmt.Printf("Debug webhook triggered with content of %d bytes\n", len(b))
		// log.Println(string(b))
	}
}

func Judge(c *gin.Context) {
	if config.INSTANCE_TYPE == config.T_WORKER {
		WriteBadRequest(c, "Not allowed on worker instance")
		return
	}

	var req JudgeRequest

	if err := c.BindJSON((&req)); err != nil {
		WriteBadRequest(c, err.Error())
		return
	}

	if err := req.Exe.Validate(); err != nil {
		WriteBadRequest(c, err.Error())
		return
	}

	if err := req.JudgeParams.Validate(); err != nil {
		WriteBadRequest(c, err.Error())
		return
	}

	output, err := judge.JudgeExecutionRequest(&req.Exe, &req.JudgeParams)

	if err != nil {
		if output != nil {
			WriteServerError(c, err.Error())
			return
		} else {
			WriteError(c, HTTP_TOO_MANY_REQ, err.Error())
			return
		}
	}

	c.IndentedJSON(200, output)

}

func RegisterWorker(c *gin.Context) {
	if config.INSTANCE_TYPE != config.T_MANAGER {
		WriteBadRequest(c, "Not allowed on non manager instances")
		return
	}

	judge.Workers_mutex.Lock()
	defer judge.Workers_mutex.Unlock()

	var req judge.WorkerRegistration

	if err := c.BindJSON(&req); err != nil {
		WriteBadRequest(c, err.Error())
		return
	}

	_, err := judge.AuthenticateInstance(req.Address, config.T_WORKER)

	if err != nil {
		WriteBadRequest(c, err.Error())
		return
	}

	res := judge.WorkerRegistrationResponse{
		AssignedInstanceId: judge.GetMaxWorkerId() + 1,
	}

	judge.AddWorker(&judge.ShowdownWorker{
		InstanceId: res.AssignedInstanceId,
		Address:    req.Address,
	})

	response, _ := json.Marshal(res)

	c.Writer.Write(response)
}

func Status(c *gin.Context) {
	f_config := c.Query("config") == "true"

	res := judge.GetInstanceState(f_config)
	response, _ := json.Marshal(res)

	c.Writer.Write(response)
}
