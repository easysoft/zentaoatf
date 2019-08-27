package fileUtils

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func ReadFile(filePath string) string {
	buf := ReadFileBuf(filePath)
	str := string(buf)
	str = commonUtils.RemoveBlankLine(str)
	return str
}

func ReadFileBuf(filePath string) []byte {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte(err.Error())
	}

	return buf
}

func WriteFile(filePath string, content string) {
	dir := path.Dir(filePath)
	MkDirIfNeeded(dir)

	var d1 = []byte(content)
	err2 := ioutil.WriteFile(filePath, d1, 0666) //写入文件(字节数组)
	check(err2)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func FileExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func GetAllScriptsInDir(dirPth string, files *[]string) error {
	sep := string(os.PathSeparator)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if fi.IsDir() { // 目录, 递归遍历
			GetAllScriptsInDir(dirPth+name+sep, files)
		} else {
			path := dirPth + name
			if CheckFileIsScript(path) {
				*files = append(*files, path)
			}
		}
	}

	return nil
}

func GetScriptByIdsInDir(dir string, idMap map[int]string, files *[]string) error {
	sep := string(os.PathSeparator)

	dir, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if fi.IsDir() { // 目录, 递归遍历
			GetAllScriptsInDir(dir+name+sep, files)
		} else {
			path := dir + name
			if CheckFileIsScript(path) {
				*files = append(*files, path)
			}
		}
	}

	return nil
}

//func GetSpecifiedFilesInWorkDir(fileNames []string) (files []string, err error) {
//	ret := make([]string, 0)
//
//	for _, file := range fileNames {
//		if !FileExist(file) {
//			continue
//		}
//
//		if path.Ext(file) == "."+constant.ExtNameSuite {
//			fileList := make([]string, 0)
//			GetCaseIdsInSuiteFile(file, &fileList)
//
//			for _, f := range fileList {
//				ret = append(ret, f)
//			}
//		} else {
//			ret = append(ret, file)
//		}
//	}
//
//	return ret, nil
//}

func GetCaseIdsInSuiteFile(name string, fileIdMap *map[int]string) {
	content := ReadFile(name)

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

func GetFailedCasesFromTestResult(resultFile string, fileIdMap *map[int]string) {
	extName := path.Ext(resultFile)

	if extName == "."+constant.ExtNameTxt { // txt format
		resultFile = strings.Replace(resultFile, extName, "."+constant.ExtNameJson, -1)
	}

	content := ReadFile(resultFile)

	var report model.TestReport
	json.Unmarshal([]byte(content), &report)

	for _, cs := range report.Cases {
		if cs.Status != constant.PASS.String() {
			(*fileIdMap)[cs.Id] = ""
		}
	}
}

func CheckFileIsScript(path string) bool {
	content := ReadFile(path)

	pass, _ := regexp.MatchString("<<<TC", content)
	return pass
}

func MkDirIfNeeded(dir string) {
	if !FileExist(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
}

func IsDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}
