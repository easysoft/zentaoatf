package main

import (
	"biz"
	"flag"
	"os"
	"utils"
)

func main() {
	langType := flag.String("l", "", "Script Language like python, php etc.")
	workDir := flag.String("p", "", "Folder that contains the scripts")

	flag.Parse()

	if *langType == "" || *workDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	files, _ := utils.GetAllFiles(*workDir, *langType)

	summaryMap := make(map[string]interface{})
	biz.RunScripts(files, *workDir, *langType, &summaryMap)

	resultMap := make(map[string]bool)
	checkpointMap := make(map[string][]string)
	biz.CheckResults(*workDir, *langType, &summaryMap, &resultMap, &checkpointMap)
	biz.Print(summaryMap, resultMap, checkpointMap, *workDir)
}
