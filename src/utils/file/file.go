package fileUtils

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"io/ioutil"
	"os"
	"path"
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

func GetAllFilesInDir(dirPth string, ext string, files *[]string) error {
	sep := string(os.PathSeparator)

	dir, err := ioutil.ReadDir(vari.Prefer.WorkDir + dirPth)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if fi.IsDir() { // 目录, 递归遍历
			if name == "res" || name == "xdoc" {
				continue
			}

			GetAllFilesInDir(dirPth+name+sep, ext, files)
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(name, "."+ext)
			if ok {
				*files = append(*files, dirPth+name)
			}
		}
	}

	return nil
}

func GetSpecifiedFilesInWorkDir(fileNames []string) (files []string, err error) {
	ret := make([]string, 0)

	for _, file := range fileNames {
		if path.Ext(file) == "."+constant.ExtNameSuite {
			fileList := make([]string, 0)
			GetSuiteFiles(file, &fileList)

			for _, f := range fileList {
				ret = append(ret, f)
			}
		} else {
			ret = append(ret, file)
		}
	}

	return ret, nil
}

func GetFailedFilesFromTestResult(resultFile string) ([]string, string) {
	ret := make([]string, 0)
	dir := ""
	extName := path.Ext(resultFile)

	if extName == "."+constant.ExtNameTxt {
		resultFile = strings.Replace(resultFile, extName, "."+constant.ExtNameJson, -1)
	}

	content := ReadFile(vari.Prefer.WorkDir + resultFile)

	var report model.TestReport
	json.Unmarshal([]byte(content), &report)

	for _, cs := range report.Cases {
		if cs.Status != constant.PASS.String() {
			ret = append(ret, cs.Path)
		}
		if dir == "" {
			dir = path.Dir(cs.Path)
		}
	}

	return ret, dir
}

func GetSuiteFiles(name string, fileList *[]string) {
	content := ReadFile(vari.Prefer.WorkDir + name)

	for _, line := range strings.Split(content, "\n") {
		file := strings.TrimSpace(line)
		if file == "" {
			return
		}

		if path.Ext(file) == "."+constant.ExtNameSuite {
			GetSuiteFiles(file, fileList)
		} else {
			*fileList = append(*fileList, file)
		}
	}
}

func MkDirIfNeeded(dir string) {
	if !FileExist(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
}
