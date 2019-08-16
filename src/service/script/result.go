package scriptService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

func LoadTestResults(assert string) []string {
	ret := make([]string, 0)

	dir := vari.Prefer.WorkDir + constant.LogDir

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

func GetRunModeAndName(assert string) (string, string) {
	ext := path.Ext(assert)
	name := strings.Replace(commonUtils.Base(assert), ext, "", -1)

	var mode string
	if ext == ".suite" {
		mode = "suite"
	} else {
		mode = "script"
	}

	return mode, name
}

func GetLogFolder(mode string, name string, date string) string {
	return fmt.Sprintf("%s-%s-%s", mode, name, date)
}
