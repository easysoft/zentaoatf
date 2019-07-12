package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/easysoft/zentaoatf/src/misc"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
	"log"
	"os"
	"strconv"
	"strings"
)

func DealwithTestCase(tc model.TestCase, langType string, independentExpect bool) {
	caseId := tc.Id
	caseTitle := tc.Title
	scriptFile := "xdoc/scripts/tc-" + strconv.Itoa(caseId) + "." + langType

	steps := make([]string, 0)
	expects := make([]string, 0)
	srcCode := make([]string, 0)

	level := 1
	checkPointIndex := 0
	for _, ts := range tc.Steps {
		DealwithTestStep(ts, langType, level, &checkPointIndex, &steps, &expects, &srcCode)
	}

	var expectsTxt string
	if independentExpect {
		expectFile := utils.ScriptToExpectName(scriptFile)

		expectsTxt = "@file"
		utils.WriteFile(expectFile, strings.Join(expects, "\n"))
	} else {
		expectsTxt = strings.Join(expects, "\n")
	}

	template := utils.ReadFile("xdoc/script-template.txt")
	content := string(fmt.Sprintf(string(template), langType, caseTitle, caseId,
		strings.Join(steps, "\n"), expectsTxt, strings.Join(srcCode, "\n")))

	fmt.Println(content)

	utils.WriteFile(scriptFile, content)
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
		expectsLine += "<" + stepIdent + " 期望结果, 可以有多行>"

		*expects = append(*expects, expectsLine)
	}

	// 处理srcCode
	if isCheckPoint {
		codeLine := ""

		if langType == misc.PHP.String() {
			codeLine += `echo "#\n";`
		} else if langType == misc.GO.String() {
			codeLine += `println("#")\n`
		}

		codeLine += "   // 验证点" + stepIdent + "的标记位，请勿删除\n"
		codeLine += "// 期待结果：" + stepExpect + "\n"
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

func usage() {
	log.Fatalf("usage: gen-project.go -p path -l lang [-e] \n")
}

func main() {
	independentExpect := flag.Bool("e", false, "Save ExpectResult in an independent file or not")
	langType := flag.String("l", "", "Script Language like python, php etc.")
	caseFile := flag.String("p", "", "Folder that contains the scripts")

	flag.Parse()

	if *caseFile == "" || *langType == "" || independentExpect == nil {
		flag.Usage()
		os.Exit(1)
	}

	buf := utils.ReadFileBuf(*caseFile)

	var resp model.Response
	json.Unmarshal(buf, &resp)

	if resp.Code != 1 {
		fmt.Println(string(buf))
		return
	}

	for _, testCase := range resp.Cases {
		DealwithTestCase(testCase, *langType, *independentExpect)
	}
}
