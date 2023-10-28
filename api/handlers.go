package api

import (
	"encoding/json"
	"msc24x/showdown/internal/engine"
	"msc24x/showdown/internal/judge"

	"github.com/gin-gonic/gin"
)

func Judge(c *gin.Context) {

	var exe_req engine.ExecutionRequest

	if err := c.BindJSON((&exe_req)); err != nil {
		WriteServerError(c, err.Error())
		return
	}

	if err := exe_req.Validate(); err != nil {
		WriteBadRequest(c, err.Error())
		return
	}

	output, err := judge.JudgeExecutionRequest(&exe_req)

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
