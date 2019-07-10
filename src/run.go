package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: run <script-dir> <langType>")
	}

	osName := runtime.GOOS
	pthSep := string(os.PathSeparator)

	dir, langType := os.Args[1], os.Args[2]
	logDir := dir + pthSep + "logs"
	mkLogDir(logDir)

	files, _ := utils.GetAllFiles(dir, langType)

	for _, file := range files {
		var command string
		if osName == "darwin" {
			logFile := strings.Replace(file, pthSep, pthSep+"logs"+pthSep,
				strings.LastIndex(file, pthSep)) + ".log"
			command = file + " > " + logFile

			if langType == "php" {
				command = langType + " " + command
			}
		}

		out, _ := utils.ExeShell(command)
		fmt.Printf(out)
	}
}

func mkLogDir(dir string) {
	if !utils.CheckFileIsExist(dir) {
		os.Mkdir(dir, os.ModePerm)
	}
}
