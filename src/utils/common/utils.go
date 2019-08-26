package commonUtils

import (
	"github.com/easysoft/zentaoatf/src/utils/const"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
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

func BoolToPass(b bool) string {
	if b {
		return constant.PASS.String()
	} else {
		return constant.FAIL.String()
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
	if _, err := os.Stat("res"); os.IsNotExist(err) {
		return true
	}

	return false
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
	if strings.Index(p, ".") == 0 {
		temp, _ := filepath.Abs(".")
		p = temp + p[1:]
	} else if strings.Index(p, "~") == 0 && !IsWin() {
		user, err := user.Current()
		if nil == err {
			temp := user.HomeDir
			p = temp + p[1:]
		}
	} else if IsRelativePath(p) {
		temp, _ = filepath.Abs(`.`)
		p = temp + sepa + p
	}

	if !PathEndWithSeparator(p) {
		p = p + sepa
	}

	return p
}

// base on workdir
func ConvertRunDir(p string) string {
	sepa := string(os.PathSeparator)

	if p == "." || p == "./" {
		p = "scripts"
	}

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

func GetSortKey(mp map[int]interface{}) []int {
	var keys []int
	for key := range mp {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	return keys
}
