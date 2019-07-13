package action

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/biz"
	. "github.com/easysoft/zentaoatf/src/misc"
	"github.com/easysoft/zentaoatf/src/model"
	. "github.com/easysoft/zentaoatf/src/utils"
)

func Run(scriptDir string, langType string) {
	p := GetInstance() // for test
	p.Printf("HELLO_1", "Peter")
	fmt.Println(p.Sprintf("HELLO_1", "Peter"))

	files, _ := GetAllFiles(scriptDir, langType)

	var report = model.TestReport{Path: scriptDir, Env: GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	biz.ExeScripts(files, scriptDir, langType, &report)

	biz.CheckResults(scriptDir, langType, &report)
	biz.Print(report, scriptDir)
}
