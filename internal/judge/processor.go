package judge

import (
	"msc24x/showdown/internal/engine"

	"github.com/google/uuid"
)

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

// Executes and judge (optional) the execution request synchronously
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
