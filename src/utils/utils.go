package utils

import (
	"github.com/easysoft/zentaoatf/res"
	"github.com/easysoft/zentaoatf/src/misc"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
)

func Base(pathStr string) string {
	pathStr = filepath.ToSlash(pathStr)
	return path.Base(pathStr)
}

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
	nameWithSuffix := Base(file)
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
		return RunMode.String() + "-" + DateTimeStrFmt(time.Now(), "2006-01-02T150405") + string(os.PathSeparator)
	}

	name := Base(filePath)

	ext := path.Ext(filePath)
	name = strings.Replace(name, ext, "", -1)

	runName := RunMode.String() + "-" + name + "-" + DateTimeStrFmt(time.Now(), "2006-01-02T150405") + string(os.PathSeparator)

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

// base on "."
func ConvertWorkDir(p string) string {
	sepa := string(os.PathSeparator)
	var temp string
	if p == "." {
		temp, _ := filepath.Abs(`.`)
		temp = temp + sepa
		return temp
	}

	if IsRelativePath(p) {
		temp, _ = filepath.Abs(`.`)
		temp = temp + sepa + p
	}
	if !PathEndWithSeparator(p) {
		temp = p + sepa
	}

	return temp
}

// base on workdir
func ConvertRunDir(p string) string {
	sepa := string(os.PathSeparator)

	//if p == "." {
	//	return Prefer.WorkDir
	//}

	if !PathEndWithSeparator(p) {
		p = p + sepa
	}

	return p
}

func IsRelativePath(path string) bool {
	var idx int

	if IsWin() {
		idx = strings.Index(path, ":")
		return idx < 0
	} else {
		idx = strings.Index(path, string(os.PathSeparator))
		return idx != 0
	}
}

func PathEndWithSeparator(path string) bool {
	idx := strings.LastIndex(path, string(os.PathSeparator))
	ret := idx == len(path)-1
	return ret
}

func ReadResData(path string) string {
	isRelease := IsRelease()

	var jsonStr string
	if isRelease {
		data, _ := res.Asset(path)
		jsonStr = string(data)
	} else {
		jsonStr = ReadFile(path)
	}

	return jsonStr
}

func GetSortKey(mp map[int]interface{}) []int {
	var keys []int
	for key := range mp {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	return keys
}
