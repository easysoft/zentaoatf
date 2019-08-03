package script

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

func LoadTestResults(assert string) []string {
	ret := make([]string, 0)

	dir := utils.Prefer.WorkDir + utils.LogDir

	mode, name := GetRunModeAndName(assert)

	reg := fmt.Sprintf("result-%s-%s-(.+)\\.txt", mode, name)
	myExp := regexp.MustCompile(reg)

	files, _ := ioutil.ReadDir(dir)
	for _, fi := range files {
		if !fi.IsDir() {

			arr := myExp.FindStringSubmatch(fi.Name())
			if len(arr) > 1 {
				ret = append(ret, arr[1])
			}
		}
	}

	return ret
}

func GetTestResult(nassert string, date string) string {
	mode, name := GetRunModeAndName(nassert)
	resultPath := utils.Prefer.WorkDir + utils.LogDir + fmt.Sprintf("result-%s-%s-%s.txt", mode, name, date)

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

	return strings.Join(arr, "\n")
}

func GetRunModeAndName(assert string) (string, string) {
	ext := path.Ext(assert)
	name := strings.Replace(path.Base(assert), ext, "", -1)

	var mode string
	if ext == ".suite" {
		mode = "suite"
	} else {
		mode = "script"
	}

	return mode, name
}

func GetLogFileByScript(file string) string {
	ext := path.Ext(file)
	logName := strings.Replace(path.Base(file), ext, ".log", -1)

	return utils.Prefer.WorkDir + utils.LogDir + logName
}
