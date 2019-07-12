package main

import (
	"biz"
	"flag"
	"fmt"
	"misc"
	"model"
	"os"
	"utils"
)

func main() {
	p := misc.GetInstance()

	langType := flag.String("l", "", "Script Language like python, php etc.")
	workDir := flag.String("p", "", "Folder that contains the scripts")

	flag.Parse()

	if *langType == "" || *workDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	p.Printf("HELLO_1", "Peter")
	fmt.Println(p.Sprintf("HELLO_1", "Peter"))

	files, _ := utils.GetAllFiles(*workDir, *langType)

	var report = model.TestReport{Path: *workDir, Env: utils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	biz.RunScripts(files, *workDir, *langType, &report)

	biz.CheckResults(*workDir, *langType, &report)
	biz.Print(report, *workDir)
}
