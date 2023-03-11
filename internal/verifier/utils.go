package verifier

import "fmt"

func restrictParamLen(rows [][]string, l int) string {
	if len(rows) != l {
		return fmt.Sprintf("parmas of task wanted %d row(s) , but got %d row(s)", l, len(rows))
	}
	return ""
}
