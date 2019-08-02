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

	ext := path.Ext(assert)
	name := strings.Replace(path.Base(assert), ext, "", -1)

	var mode string
	if ext == ".suite" {
		mode = "suite"
	} else {
		mode = "script"
	}

	reg := fmt.Sprintf("result-%s-(%s-.+)\\.txt", mode, name)
	myExp := regexp.MustCompile(reg)
	utils.PrintToCmd(utils.Cui, reg)

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
