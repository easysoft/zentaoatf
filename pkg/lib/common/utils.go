package commonUtils

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"github.com/emirpasic/gods/maps"
)

func Base(pathStr string) string {
	pathStr = filepath.ToSlash(pathStr)
	return path.Base(pathStr)
}

func RemoveBlankLine(str string) string {
	myExp := regexp.MustCompile(`\n{3,}`) // 连续换行
	ret := myExp.ReplaceAllString(str, "\n\n")

	ret = strings.Trim(ret, "\n")
	ret = strings.TrimSpace(ret)

	return ret
}

func GetOs() commConsts.OsType {
	osName := runtime.GOOS

	if osName == "darwin" {
		return commConsts.OsMac
	} else {
		return commConsts.OsType(osName)
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

func IntToStrArr(intArr []int) (strArr []string) {
	for _, i := range intArr {
		strArr = append(strArr, strconv.Itoa(i))
	}

	return
}
func UintToStrArr(intArr []uint) (strArr []string) {
	for _, i := range intArr {
		strArr = append(strArr, fmt.Sprintf("%d", i))
	}

	return
}

func LinkedMapToMap(mp maps.Map) map[string]string {
	ret := make(map[string]string, 0)

	for _, keyIfs := range mp.Keys() {
		valueIfs, _ := mp.Get(keyIfs)

		key := strings.TrimSpace(keyIfs.(string))
		value := strings.TrimSpace(valueIfs.(string))

		ret[key] = value
	}

	return ret
}

func IsRelease() bool {
	arg1 := strings.ToLower(os.Args[0])
	name := filepath.Base(arg1)

	ret := strings.Index(arg1, "go-build") < 0 &&
		strings.Index(name, "___") != 0 && strings.Index(name, "go-build") != 0

	return ret
}

func GetUserHome() string {
	userProfile, _ := user.Current()
	home := userProfile.HomeDir
	return home
}

func IsPortInUse(port int) bool {
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprintf("%d", port)), 3*time.Second); err == nil {
		conn.Close()
		return true
	}
	return false
}

func IsDisable(enable string) bool {
	if enable == "1" {
		return false
	} else {
		return true
	}
}

func IgnoreZtfFile(path string) bool {
	path = filepath.Base(path)

	arr := []string{"bin", "release", "logs", "xdoc",
		"log", "log-bak", "conf"}
	if strings.Index(path, ".") == 0 ||
		stringUtils.FindInArr(path, arr) ||
		(len(path) >= 4 && path[len(path)-4:] == ".exp") {
		return true
	} else {
		return false
	}
}
func IgnoreCodeFile(path string) bool {
	path = filepath.Base(path)

	arr := []string{"bin", "node_modules", ".git"}
	if strings.Index(path, ".") == 0 || stringUtils.FindInArr(path, arr) {
		return true
	} else {
		return false
	}
}

func AddSlashForUrl(url string) string {
	if strings.LastIndex(url, "/") < len(url)-1 {
		url += "/"
	}

	return url
}

//func ChangeScriptForDebug(dir *string) {
//	if !IsRelease() { // debug in ide
//		*dir = filepath.Join(*dir, "demo", "sample")
//	}
//}

func GetDebugParamForRun(args []string) (debug string, ret []string) {
	index := -1
	for i, item := range args {
		if item == "-debug" || item == "--debug" {
			index = i
			break
		}
	}

	if index > -1 && len(args) > index+1 {
		debug = args[index+1]
		ret = append(args[0:index], args[index+2:]...)
	} else {
		ret = args
	}

	return
}
