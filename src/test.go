package main

import (
	"fmt"
	"regexp"
)

func main() {
	text := ` df 
 FAIL scripts/tc-200.py
  Step1: FAIL   @step2010 第4次尝试登录
    Checkpoint1: FAIL
      Expect Result CODE: @step2010期望结果, 可以有多行
      Actual Result N/A

  Step2: FAIL   @step2104 再输入1次正确的密码
    Checkpoint1: FAIL
      Expect Result CODE: @step2104期望结果, 可以有多行
      Actual Result N/A
 
dd`

	str := "(?m:^\\s" + "FAIL\\sscripts/tc-200.py" + "\\n([\\s\\S]*?)((^\\s(PASS|FAIL))|\\z))"
	// myExp := regexp.MustCompile("(?m:^\\s(?:PASS|FAIL) scripts/tc-200.py\n([\\s\\S]*?)((^\\s(PASS|FAIL))|\\z))")

	myExp := regexp.MustCompile(str)
	arr := myExp.FindStringSubmatch(text)

	fmt.Println(arr[1])
}
