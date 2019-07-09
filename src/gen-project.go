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

	for _, ts := range tc.Steps {
		DealwithTestStep(ts, langType, level, &steps, &expects, &srcCode)
	}

	template := utils.ReadFile("xdoc/script-template.txt")
	content := string(fmt.Sprintf(string(template), langType, caseTitle, caseId,
		strings.Join(steps, "\n"), strings.Join(expects, "\n"), strings.Join(srcCode, "\n")))

	fmt.Println(content)

	utils.WriteFile("xdoc/tc-"+strconv.Itoa(caseId)+"."+langType, content)
}

func DealwithTestStep(ts model.TestStep, langType string, level int,
	steps *[]string, expects *[]string, srcCode *[]string) {
	isGroup := ts.IsGroup
	isCheckPoint := ts.IsCheckPoint

	stepId := ts.Id
	stepTitle := ts.Title
	// stepExpect := ts.Expect

	// 处理steps
	stepLine := ""

	var stepType string
	if isGroup {
		stepType = "group"
	} else {
		stepType = "step"
	}
	stepLine += stepType + strconv.Itoa(stepId) + " // " + stepTitle
	if isCheckPoint {
		stepLine = "@" + stepLine
	}
	i := level
	for {
		stepLine = "   " + stepLine
		i--
		if i == 0 {
			break
		}
	}
	stepLineSimple := strings.Replace(strings.TrimLeft(stepLine, " "), "//", "-", -1)
	*steps = append(*steps, stepLine)

	// 处理expects
	if isCheckPoint {
		expectsLine := ""

		expectsLine = "# \n" //  + stepLineSimple + " " + stepExpect + "\n"
		expectsLine += "<" + stepType + strconv.Itoa(stepId) + " 期望结果, 可以有多行>\n"

		*expects = append(*expects, expectsLine)
	}

	// 处理srcCode

	if isCheckPoint {
		codeLine := ""

		codeLine = "# " + stepLineSimple + "\n"
		codeLine += "// 此处编写上述验证点代码 \n"

		*srcCode = append(*srcCode, codeLine)
	}

	if isGroup {
		for _, tsChild := range ts.Steps {
			DealwithTestStep(tsChild, langType, level+1, steps, expects, srcCode)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
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
