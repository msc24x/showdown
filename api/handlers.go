package api

import (
	"encoding/json"
	"fmt"
	"io"
	"msc24x/showdown/config"
	"msc24x/showdown/internal/engine"
	"msc24x/showdown/internal/judge"

	"github.com/gin-gonic/gin"
)

// Struct to define http request to showdown
type JudgeRequest struct {
	JudgeParams judge.Params            `json:"judge_params"`
	Exe         engine.ExecutionRequest `json:"exe"`
}

func Tmp(c *gin.Context) {
	if b, err := io.ReadAll(c.Request.Body); err == nil {
		fmt.Printf("Dummy webhook triggered with content of %d bytes\n", len(b))
	}
}

func Judge(c *gin.Context) {
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
		WriteServerError(c, err.Error())
		return
	}

	c.IndentedJSON(200, output)

}

func RegisterWorker(c *gin.Context) {
	judge.Workers_mutex.Lock()
	defer judge.Workers_mutex.Unlock()

	var req judge.WorkerRegisteration

	if err := c.BindJSON(&req); err != nil {
		WriteBadRequest(c, err.Error())
		return
	}

	_, err := judge.AuthenticateInstance(req.Address, config.T_WORKER)

	if err != nil {
		WriteBadRequest(c, err.Error())
		return
	}

	res := judge.WorkerRegisterationResponse{
		AssignedInstanceId: judge.GetMaxWorkerId() + 1,
	}

	response, _ := json.Marshal(res)

	judge.AddWorker(&judge.ShowdownWorker{
		InstanceId:    res.AssignedInstanceId,
		Address:       req.Address,
		Authenticated: true,
	})

	c.Writer.Write(response)
}

func GetStats(c *gin.Context) {
	stats := judge.GetState()

	res := judge.InstanceState{
		InstanceId:   config.INSTANCE_ID,
		InstanceType: config.INSTANCE_TYPE,
		Private:      config.ACCESS_TOKEN != "",
	}

	if config.INSTANCE_TYPE != config.T_MANAGER {
		res.JudgeStats = stats
	} else if config.INSTANCE_TYPE != config.T_WORKER {
		res.Workers = judge.GetWorkers(true)
	}

	response, _ := json.Marshal(res)

	c.Writer.Write(response)
}
