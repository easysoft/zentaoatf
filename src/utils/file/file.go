package fileUtils

import (
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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

func AbosutePath(pth string) string {
	if !path.IsAbs(pth) {
		pth, _ = filepath.Abs(pth)
	}

	pth = UpdateDir(pth)

	return pth
}

func UpdateDir(path string) string {
	sepa := string(os.PathSeparator)

	if strings.LastIndex(path, sepa) < len(path)-1 {
		path += sepa
	}
	return path
}
