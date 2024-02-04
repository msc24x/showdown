// Package judge utilizes the package engine and judges the execution output
package judge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"msc24x/showdown/internal/engine"
	"net/http"

	"github.com/google/uuid"
)

type ExecutionResponse struct {
	PID      string `json:"pid"`
	Webhook  string `json:"webhook"`
	Success  bool   `json:"success"`
	Judged   bool   `json:"judged"`
	Error    string `json:"error"`
	Output   string `json:"output"`
	Meta     string `json:"meta"`
	Expected string `json:"expected"`
}

type Params struct {
	Webhook    string `json:"webhook"`
	DoNotJudge bool   `json:"donotjudge`
}

// Creates a matrix of runes having each line separated. While parsing '\r'
// are ignored to support CRLF encoding and extra white spaces along with
// blank lines are removed.
func stringToLineBuffers(data string) [][]rune {
	var (
		output_lines [][]rune
		line_buf     []rune
		buf_len      int = 0
	)

	for _, val := range data {

		// ignore '\r'
		if val == '\r' {
			continue
		}

		// capture a line when '\n'
		if val == '\n' {
			if buf_len > 0 && line_buf[buf_len-1] == ' ' {
				line_buf = line_buf[:buf_len-1]
				buf_len--
			}
			if buf_len == 0 {
				continue
			}
			output_lines = append(output_lines, line_buf)
			line_buf = nil
			buf_len = 0
			continue
		}

		// white space skip to clean line
		if buf_len > 0 && line_buf[buf_len-1] == ' ' && val == ' ' || buf_len == 0 && val == ' ' {
			continue
		}

		// buffer the line
		line_buf = append(line_buf, val)
		buf_len++

	}

	if buf_len != 0 {
		if line_buf[buf_len-1] == ' ' {
			line_buf = line_buf[:buf_len-1]
			buf_len--
		}
		if buf_len > 0 {
			output_lines = append(output_lines, line_buf)
		}
	}

	return output_lines

}

// Core function to determine the correctness of the output with the expected
// output.
func judgeLines(test [][]rune, truth [][]rune) bool {
	if len(test) != len(truth) {
		return false
	}

	var current_pos = 0

	for current_line := 0; current_line < len(test); current_line++ {
		if len(test[current_line]) != len(truth[current_line]) {
			return false
		}
		current_pos = 0
		for ; current_pos < len(test[current_line]); current_pos++ {
			if test[current_line][current_pos] != truth[current_line][current_pos] {
				return false
			}
		}
	}

	return true
}

func processRequest(pid uuid.UUID, exe_req *engine.ExecutionRequest, params *Params, response *ExecutionResponse) error {
	engine := engine.BaseEngine{}

	if err := engine.Init(pid, exe_req); err != nil {
		return err
	}
	defer engine.Clean()

	output, err := engine.Execute()
	response.Output = string(output)

	if !params.DoNotJudge {
		response.Judged = true
		response.Expected = string(exe_req.Output)
		meta_info, _ := engine.CollectMeta()
		response.Meta = string(meta_info)

		if err != nil {
			response.Error = err.Error()
			return nil
		}

		test_output_lines := stringToLineBuffers(string(output))
		truth_output_lines := stringToLineBuffers(exe_req.Output)

		response.Success = judgeLines(test_output_lines, truth_output_lines)
	}

	return nil
}

func processAndSendRequest(url string, callback *func() (*ExecutionResponse, error)) {
	res, err := (*callback)()
	var res_text string

	if err != nil {
		log.Fatalln(err.Error())
		res_text = fmt.Sprintf("Your Showdown request with PID %s could not be completed. Report if the issue persists.", res.PID)
	} else {
		res_bytes, _ := json.Marshal(res)
		res_text = string(res_bytes)
	}

	content := bytes.NewBufferString(res_text)
	http.Post(url, "application/json", content)
}

// Entrypoint for an HTTP execution request
func JudgeExecutionRequest(exe_req *engine.ExecutionRequest, params *Params) (*ExecutionResponse, error) {

	pid, err := OnboardProcess(exe_req.UID)
	if err != nil {
		return nil, err
	}

	response := ExecutionResponse{}
	response.PID = pid.String()
	response.Webhook = params.Webhook

	if params.Webhook != "" {
		callback := func() (*ExecutionResponse, error) {
			defer OffboardProcess(pid, exe_req.UID)
			return &response, processRequest(pid, exe_req, params, &response)
		}

		go processAndSendRequest(params.Webhook, &callback)
		return &response, nil
	} else {
		defer OffboardProcess(pid, exe_req.UID)
	}

	err = processRequest(pid, exe_req, params, &response)

	if err != nil {
		return &response, err
	}

	return &response, nil

}
