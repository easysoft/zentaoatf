package utils

import (
	"github.com/easysoft/zentaoatf/src/misc"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func RemoveBlankLine(str string) string {
	myExp := regexp.MustCompile(`\n{2,}`) // 连续换行
	ret := myExp.ReplaceAllString(str, "\n")

	ret = strings.Trim(ret, "\n")
	ret = strings.TrimSpace(ret)

	return ret
}

func ScriptToLogName(dir string, file string) string {
	logDir := dir + LogDir + RunDir
	MkDirIfNeeded(logDir)

	nameSuffix := path.Ext(file)
	nameWithSuffix := path.Base(file)
	name := strings.TrimSuffix(nameWithSuffix, nameSuffix)

	logFile := logDir + name + ".log"

	return logFile
}

func ScriptToExpectName(file string) string {
	fileSuffix := path.Ext(file)
	expectName := strings.TrimSuffix(file, fileSuffix) + ".ex"

	return expectName
}

func PathToRunName(filePath string) string {
	if filePath == "" {
		return RunMode.String() + "-" + DateTimeStrFmt(time.Now(), "2006-01-02 15:04:05") + "/"
	}

	name := path.Base(filePath)
	ext := path.Ext(filePath)
	name = strings.Replace(name, ext, "", -1)

	runName := RunMode.String() + "-" + name + "-" + DateTimeStrFmt(time.Now(), "2006-01-02 15:04:05") + "/"

	return runName
}

func BoolToPass(b bool) string {
	if b {
		return misc.PASS.String()
	} else {
		return misc.FAIL.String()
	}
}

func GetOs() string {
	osName := runtime.GOOS

	if osName == "darwin" {
		return "mac"
	} else {
		return osName
	}
}
func IsWin() bool {
	return GetOs() == "windows"
}
func IsLinux() bool {
	return GetOs() == "linux"
}
func IsMac() bool {
	return GetOs() == "mac"
}

func IsRelease() bool {
	return !FileExist("res")
}

func UpdateUrl(url string) string {
	if strings.LastIndex(url, "/") < len(url)-1 {
		url += "/"
	}
	return url
}
