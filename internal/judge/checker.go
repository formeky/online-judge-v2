package judge

import (
	"bytes"
)

func compareOutput(got, expected []byte) bool {
	gotLines := normalizeLines(got)
	expLines := normalizeLines(expected)
	if len(gotLines) != len(expLines) {
		return false
	}
	for i := range gotLines {
		if !bytes.Equal(gotLines[i], expLines[i]) {
			return false
		}
	}
	return true
}

func normalizeLines(data []byte) [][]byte {
	lines := bytes.Split(data, []byte("\n"))
	var result [][]byte
	for _, line := range lines {
		trimmed := bytes.TrimRight(line, " \t\r")
		result = append(result, trimmed)
	}
	for len(result) > 0 && len(result[len(result)-1]) == 0 {
		result = result[:len(result)-1]
	}
	return result
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
