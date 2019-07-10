package utils

import (
	"io/ioutil"
	"os"
	"strings"
)

func ReadFile(filePath string) []byte {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil
	}

	return buf
}

func WriteFile(filePath string, content string) {
	var d1 = []byte(content)
	err2 := ioutil.WriteFile(filePath, d1, 0666) //写入文件(字节数组)
	check(err2)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func GetAllFiles(dirPth string, ext string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth+PthSep+fi.Name(), ext)
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), "."+ext)
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table, ext)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}
