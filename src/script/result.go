package script

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

func LoadTestResults(assert string) []string {
	ret := make([]string, 0)

	dir := utils.Prefer.WorkDir + utils.LogDir

	mode, name := GetRunModeAndName(assert)

	reg := fmt.Sprintf("%s-%s-(.+)", mode, name)
	myExp := regexp.MustCompile(reg)

	files, _ := ioutil.ReadDir(dir)
	for _, fi := range files {
		if fi.IsDir() {
			arr := myExp.FindStringSubmatch(fi.Name())
			if len(arr) > 1 {
				ret = append(ret, arr[1])
			}
		}
	}

	return ret
}

func GetTestResult(assert string, date string) []string {
	mode, name := GetRunModeAndName(assert)
	resultPath := utils.Prefer.WorkDir + utils.LogDir + logFolder(mode, name, date) + string(os.PathSeparator) + "result.txt"

	arr := make([]string, 0)
	content := utils.ReadFile(resultPath)
	for _, line := range strings.Split(content, "\n") {
		pass, _ := regexp.MatchString("^\\s(PASS|FAIL).*", line)
		if !pass {
			continue
		}

		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}

		arr = append(arr, line)
	}

	return arr
}

func GetCheckpointsResult(assert string, date string, caseLine string) string {
	mode, name := GetRunModeAndName(assert)
	resultPath := utils.Prefer.WorkDir + utils.LogDir + logFolder(mode, name, date) + string(os.PathSeparator) + "result.txt"

	content := utils.ReadFile(resultPath)

	caseLine = strings.Replace(caseLine, " ", "\\s", -1)

	myExp := regexp.MustCompile(`(?m:^\s` + caseLine + `\n([\s\S]*?)((^\s(PASS|FAIL))|\z))`)
	arr := myExp.FindStringSubmatch(content)
	str := ""
	if len(arr) > 1 {
		str = arr[1]
	}

	return str
}

func GetRunModeAndName(assert string) (string, string) {
	ext := path.Ext(assert)
	name := strings.Replace(utils.Base(assert), ext, "", -1)

	var mode string
	if ext == ".suite" {
		mode = "suite"
	} else {
		mode = "script"
	}

	return mode, name
}

func GetLogFileByCase(assert string, date string, file string) string {
	mode, name := GetRunModeAndName(assert)

	ext := path.Ext(file)
	logName := strings.Replace(utils.Base(file), ext, ".log", -1)

	return utils.Prefer.WorkDir + utils.LogDir + logFolder(mode, name, date) + string(os.PathSeparator) + logName
}

func logFolder(mode string, name string, date string) string {
	return fmt.Sprintf("%s-%s-%s", mode, name, date)
}
