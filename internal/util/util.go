package util

import (
	"os"
	"strings"
)

func IsCommentLine(line string) bool {
	switch {
	case strings.HasPrefix(line, "#"):
	case strings.HasPrefix(line, ";"):
	case strings.HasPrefix(line, "//"):
	default:
		return false
	}

	return true
}

func StringTrim(str string) string {
	ret := str
	ret = strings.TrimSpace(str)
	ret = strings.Trim(ret, "\t")
	ret = strings.Trim(ret, "\r")
	ret = strings.Trim(ret, "\n")

	return ret
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
