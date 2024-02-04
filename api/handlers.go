package api

import (
	"encoding/json"
	"fmt"
	"io"
	"msc24x/showdown/internal/engine"
	"msc24x/showdown/internal/judge"

	"github.com/gin-gonic/gin"
)

type JudgeRequest struct {
	JudgeParams judge.Params            `json:"judge_params"`
	Exe         engine.ExecutionRequest `json:"exe"`
}

func Tmp(c *gin.Context) {
	if b, err := io.ReadAll(c.Request.Body); err == nil {
		fmt.Println(string(b))
	}
}

func Judge(c *gin.Context) {
	var req JudgeRequest

	if err := c.BindJSON((&req)); err != nil {
		WriteServerError(c, err.Error())
		return
	}

	if err := req.Exe.Validate(); err != nil {
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

func GetStats(c *gin.Context) {
	stats := judge.GetState()

	response, _ := json.Marshal(stats)

	c.Writer.Write(response)

}
