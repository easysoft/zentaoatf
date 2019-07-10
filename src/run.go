package main

import (
	"biz"
	"fmt"
	"os"
	"utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: run <script-dir> <langType>")
	}

	workDir, langType := os.Args[1], os.Args[2]

	files, _ := utils.GetAllFiles(workDir, langType)
	biz.RunScripts(files, workDir, langType)
	biz.CheckResults(workDir, langType)
}
