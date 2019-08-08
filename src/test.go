package main

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils"
	"regexp"
)

func main() {
	text := utils.ReadFile("result.txt")

	str := "(?m:^\\s" + "FAIL\\sscripts\\\\tc-200.py" + "\\n([\\s\\S]*?)((^\\s(PASS|FAIL))|\\z))"
	// myExp := regexp.MustCompile("(?m:^\\s(?:PASS|FAIL) scripts\\tc-200.py\n([\\s\\S]*?)((^\\s(PASS|FAIL))|\\z))")

	myExp := regexp.MustCompile(str)
	arr := myExp.FindStringSubmatch(text)

	fmt.Println(arr[1])
}
