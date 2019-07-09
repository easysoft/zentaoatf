package utils

import (
	"io/ioutil"
	"os"
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

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
