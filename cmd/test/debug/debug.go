package main

import (
	"flag"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	execHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
)

var resultfile string

func init() {
	flag.StringVar(&resultfile, "resultfile", "", "result file")
}

func main() {
	flag.Parse()
	req := serverDomain.TestSet{
		TestTool:  commConsts.PyTest,
		ResultDir: resultfile,
	}
	testSuites := execHelper.RetrieveUnitResult(req, 0)
	spew.Dump(testSuites)
	cases, _, _ := execHelper.ParserUnitTestResult(testSuites)
	spew.Dump(cases)
	failedCaseLines := make([]string, 0)
	failedCaseLinesDesc := make([]string, 0)
	Fail, Pass, Total := 0, 0, 0
	for idx, cs := range cases {
		if cs.Failure != nil {
			Fail++

			className := cases[idx].TestSuite

			line := fmt.Sprintf("[%s] %s", className, cs.Title)
			failedCaseLines = append(failedCaseLines, line)
			failedCaseLinesDesc = append(failedCaseLinesDesc, line)
			failDesc := fmt.Sprintf("   %s - %s", cs.Failure.Type, cs.Failure.Desc)
			failedCaseLinesDesc = append(failedCaseLinesDesc, failDesc)
		} else {
			Pass++
		}

		Total++
	}
	fmt.Printf("Pass: %d, Fail: %d, Total: %d\n", Pass, Fail, Total)
}
