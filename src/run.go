package main

import (
	"flag"
	"fmt"
	"github.com/easysoft/zentaoatf/src/biz"
	. "github.com/easysoft/zentaoatf/src/misc"
	"github.com/easysoft/zentaoatf/src/model"
	. "github.com/easysoft/zentaoatf/src/utils"
	"os"
)

func main() {
	p := GetInstance()

	langType := flag.String("l", "", "Script Language like python, php etc.")
	workDir := flag.String("p", "", "Folder that contains the scripts")

	flag.Parse()

	if *langType == "" || *workDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	p.Printf("HELLO_1", "Peter")
	fmt.Println(p.Sprintf("HELLO_1", "Peter"))

	files, _ := GetAllFiles(*workDir, *langType)

	var report = model.TestReport{Path: *workDir, Env: GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	biz.RunScripts(files, *workDir, *langType, &report)

	biz.CheckResults(*workDir, *langType, &report)
	biz.Print(report, *workDir)
}
