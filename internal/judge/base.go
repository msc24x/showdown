package judge

import (
	"msc24x/showdown/internal/engine"
)

type ExecutionResponse struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Output   string `json:"output"`
	Expected string `json:"expected"`
}

func StringToLineBuffers(data string) [][]rune {
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

func JudgeLines(test [][]rune, truth [][]rune) bool {
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

func JudgeExecutionRequest(exe_req *engine.ExecutionRequest) (*ExecutionResponse, error) {

	err := OnboardProcess(exe_req.PID, exe_req.UID)
	if err != nil {
		return nil, err
	}
	defer OffboardProcess(exe_req.PID, exe_req.UID)

	engine := engine.BaseEngine{}

	engine.Init(exe_req)
	defer engine.Clean()

	output, err := engine.Execute()

	test_output_lines := StringToLineBuffers(string(output))
	truth_output_lines := StringToLineBuffers(exe_req.Output)

	response := ExecutionResponse{}

	response.Success = JudgeLines(test_output_lines, truth_output_lines)
	response.Output = string(output)
	response.Expected = string(exe_req.Output)

	if err != nil {
		response.Error = err.Error()

	}

	return &response, nil

}
