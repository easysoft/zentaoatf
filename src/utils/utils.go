package utils

import (
	"os"
	"path"
	"regexp"
	"strings"
)

func RemoveBlankLine(str string) string {
	myExp := regexp.MustCompile(`\n{2,}`) // 连续换行
	ret := myExp.ReplaceAllString(str, "\n")

	myExp = regexp.MustCompile(`#[^\n]*\n`) // 空行
	ret = myExp.ReplaceAllString(ret, "")

	return ret
}

func ScriptToLogName(file string) string {
	pthSep := string(os.PathSeparator)

	dir := path.Dir(file)

	nameSuffix := path.Ext(file)
	nameWithSuffix := path.Base(file)
	name := strings.TrimSuffix(nameWithSuffix, nameSuffix)

	logFile := dir + pthSep + "logs" + pthSep + name + ".log"

	return logFile
}

func ScriptToExpectName(file string) string {
	fileSuffix := path.Ext(file)
	expectName := strings.TrimSuffix(file, fileSuffix) + ".ex"

	return expectName
}
