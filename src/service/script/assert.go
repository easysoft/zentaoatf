package scriptService

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	langUtils "github.com/easysoft/zentaoatf/src/utils/lang"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func GetAllScriptsInDir(filePth string, files *[]string) error {
	filePth = fileUtils.AbosutePath(filePth)
	sep := string(os.PathSeparator)

	if !fileUtils.IsDir(filePth) { // first call, param is file
		regx := langUtils.GetSupportLangageRegx()
		pass, _ := regexp.MatchString("^*.\\."+regx+"$", filePth)

		if pass {
			id, _, _ := zentaoUtils.GetCaseInfo(filePth)

			if id > 0 {
				*files = append(*files, filePth)
			}
		}

		return nil
	}

	dir, err := ioutil.ReadDir(filePth)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if commonUtils.IngoreFile(name) {
			continue
		}

		if fi.IsDir() { // 目录, 递归遍历
			GetAllScriptsInDir(filePth+name+sep, files)
		} else {
			path := filePth + name
			regx := langUtils.GetSupportLangageRegx()
			pass, _ := regexp.MatchString("^*.\\."+regx+"$", path)

			if pass {
				id, _, _ := zentaoUtils.GetCaseInfo(path)

				if id > 0 {
					*files = append(*files, path)
				}
			}
		}
	}

	return nil
}

func GetScriptByIdsInDir(dirPth string, idMap map[int]string, files *[]string) error {
	dirPth = fileUtils.AbosutePath(dirPth)

	sep := string(os.PathSeparator)

	if commonUtils.IngoreFile(dirPth) {
		return nil
	}

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if fi.IsDir() { // 目录, 递归遍历
			GetScriptByIdsInDir(dirPth+name+sep, idMap, files)
		} else {
			regx := langUtils.GetSupportLangageRegx()
			pass, _ := regexp.MatchString("^*.\\."+regx+"$", name)

			if !pass {
				continue
			}

			path := dirPth + name
			id, _, _ := zentaoUtils.GetCaseInfo(path)

			if id > 0 {
				_, ok := idMap[id]

				if ok {
					*files = append(*files, path)
				}
			}
		}
	}

	return nil
}

func GetCaseIdsInSuiteFile(name string, fileIdMap *map[int]string) {
	content := fileUtils.ReadFile(name)

	for _, line := range strings.Split(content, "\n") {
		idStr := strings.TrimSpace(line)
		if idStr == "" {
			continue
		}

		id, err := strconv.Atoi(idStr)
		if err == nil {
			(*fileIdMap)[id] = ""
		}
	}
}

func GetFailedCasesDirectlyFromTestResult(resultFile string, cases *[]string) {
	extName := path.Ext(resultFile)

	if extName == "."+constant.ExtNameResult {
		resultFile = strings.Replace(resultFile, extName, "."+constant.ExtNameJson, -1)
	}

	content := fileUtils.ReadFile(resultFile)

	var report model.TestReport
	json.Unmarshal([]byte(content), &report)

	for _, cs := range report.Cases {
		*cases = append(*cases, cs.Path)
	}
}
