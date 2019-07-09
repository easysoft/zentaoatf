package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Response struct {
	Code  int
	Cases []TestCase
}

type TestCase struct {
	Id    int
	Title string
	Steps []TestStep
}

type TestStep struct {
	Id    int
	Title string
	Steps []TestStep

	Expect       string
	IsGroup      bool
	IsCheckPoint bool
}

func ReadAll(filePth string) []byte {
	buf, err := ioutil.ReadFile(filePth)
	if err != nil {
		return nil
	}

	return buf
}

func DealwithTestCase(tc TestCase, langType string) {
	caseId := tc.Id
	caseTitle := tc.Title

	steps := []string{}
	expects := []string{}
	srcCode := []string{}

	level := 1

	for _, ts := range tc.Steps {
		DealwithTestStep(ts, langType, level, &steps, &expects, &srcCode)
	}

	template := ReadAll("xdoc/script-template.txt")
	content := string(fmt.Sprintf(string(template), langType, caseTitle, caseId,
		strings.Join(steps, "\n"), strings.Join(expects, "\n"), strings.Join(srcCode, "\n")))

	fmt.Println(content)
}

func DealwithTestStep(ts TestStep, langType string, level int,
	steps *[]string, expects *[]string, srcCode *[]string) {
	isGroup := ts.IsGroup
	isCheckPoint := ts.IsCheckPoint

	stepId := ts.Id
	stepTitle := ts.Title
	stepExpect := ts.Expect

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

		expectsLine = "# " + stepLineSimple + " " + stepExpect + "\n"
		expectsLine += "<期望值字符串, 可以有多行>\n"

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
	buf := ReadAll(caseFile)

	var resp Response
	json.Unmarshal(buf, &resp)

	if resp.Code != 1 {
		fmt.Println(string(buf))
		return
	}

	for _, testCase := range resp.Cases {
		DealwithTestCase(testCase, langType)
	}
}
