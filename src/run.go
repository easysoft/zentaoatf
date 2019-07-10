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

	summaryMap := make(map[string]interface{})
	biz.RunScripts(files, workDir, langType, &summaryMap)

	resultMap := make(map[string]bool)
	checkpointMap := make(map[string][]string)
	biz.CheckResults(workDir, langType, &summaryMap, &resultMap, &checkpointMap)
	biz.Print(summaryMap, resultMap, checkpointMap)
}
