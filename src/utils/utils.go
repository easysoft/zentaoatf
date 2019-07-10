package utils

import (
	"os"
	"strings"
)

func ScriptToLog(file string) string {
	pthSep := string(os.PathSeparator)

	logFile := strings.Replace(file, pthSep, pthSep+"logs"+pthSep,
		strings.LastIndex(file, pthSep)) + ".log"

	return logFile
}
