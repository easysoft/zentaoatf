package scriptService

import (
	"fmt"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	config2 "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

/**
Get all test script and suite files in current work dir
*/
func LoadAssetFiles() ([]string, []string) {
	config := config2.ReadCurrConfig()
	ext := GetSupportedScriptLang()[config.LangType]["extName"]

	caseFiles := make([]string, 0)
	suitesFiles := make([]string, 0)

	fileUtils.GetAllFiles(vari.Prefer.WorkDir+constant.ScriptDir, ext, &caseFiles)
	fileUtils.GetAllFiles(vari.Prefer.WorkDir+constant.ScriptDir, "suite", &suitesFiles)

	return caseFiles, suitesFiles
}

/**
Get all test result histories for specific test script/suite
*/
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

/**
Run mode: refer to utils/const/enum
*/
func GetRunModeAndName(assert string) (string, string) {
	ext := path.Ext(assert)
	name := strings.Replace(commonUtils.Base(assert), ext, "", -1)

	var mode string
	if ext == ".suite" {
		mode = constant.RunModeSuite.String()
	} else {
		mode = constant.RunModeScript.String()
	}

	return mode, name
}

func GetLogFolder(mode string, name string, date string) string {
	return fmt.Sprintf("%s-%s-%s", mode, name, date)
}
