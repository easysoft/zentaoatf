package main

import (
	"encoding/json"
	"fmt"
	"model"
	"os"
	"strconv"
	"strings"
	"utils"
)

func DealwithTestCase(tc model.TestCase, langType string) {
	caseId := tc.Id
	caseTitle := tc.Title

	steps := make([]string, 0)
	expects := make([]string, 0)
	srcCode := make([]string, 0)

	level := 1
	checkPointIndex := 0
	for _, ts := range tc.Steps {
		DealwithTestStep(ts, langType, level, &checkPointIndex, &steps, &expects, &srcCode)
	}

	template := utils.ReadFile("xdoc/script-template.txt")
	content := string(fmt.Sprintf(string(template), langType, caseTitle, caseId,
		strings.Join(steps, "\n"), strings.Join(expects, "\n"), strings.Join(srcCode, "\n")))

	fmt.Println(content)

	utils.WriteFile("xdoc/tc-"+strconv.Itoa(caseId)+"."+langType, content)
}

func DealwithTestStep(ts model.TestStep, langType string, level int, checkPointIndex *int,
	steps *[]string, expects *[]string, srcCode *[]string) {
	isGroup := ts.IsGroup
	isCheckPoint := ts.IsCheckPoint

	stepId := ts.Id
	stepTitle := ts.Title
	stepExpect := ts.Expect

	// 处理steps
	var stepType string
	if isGroup {
		stepType = "group"
	} else {
		stepType = "step"
	}

	stepIdent := stepType + strconv.Itoa(stepId)
	if isCheckPoint {
		stepIdent = "@" + stepIdent
		(*checkPointIndex)++
	}
	i := level

	stepLine := stepIdent + " // " + stepTitle
	for {
		stepLine = "   " + stepLine
		i--
		if i == 0 {
			break
		}
	}
	*steps = append(*steps, stepLine)

	// 处理expects
	if isCheckPoint {
		expectsLine := ""

		expectsLine = "# \n"
		expectsLine += "<" + stepIdent + " 期望结果, 可以有多行>\n"

		*expects = append(*expects, expectsLine)
	}

	// 处理srcCode
	if isCheckPoint {
		codeLine := ""

		codeLine = "# " + stepIdent + " - " + stepExpect + "\n"
		codeLine += "// 此处编写上述验证点代码"
		if *checkPointIndex == 1 {
			codeLine += "，输出实际结果, 可以有多行 \n"
		} else {
			codeLine += "\n"
		}

		*srcCode = append(*srcCode, codeLine)
	}

	if isGroup {
		for _, tsChild := range ts.Steps {
			DealwithTestStep(tsChild, langType, level+1, checkPointIndex, steps, expects, srcCode)
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: gen-project.go <path> <langType>")
	}

	caseFile, langType := os.Args[1], os.Args[2]
	buf := utils.ReadFile(caseFile)

	var resp model.Response
	json.Unmarshal(buf, &resp)

	if resp.Code != 1 {
		fmt.Println(string(buf))
		return
	}

	for _, testCase := range resp.Cases {
		DealwithTestCase(testCase, langType)
	}
}
