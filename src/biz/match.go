package biz

import (
	"github.com/easysoft/zentaoatf/src/misc"
	"regexp"
	"strings"
)

func MatchString(expect string, actual string, langType string) bool {
	if langType == misc.PHP.String() {
		expect = strings.Replace(expect, "%s", `.+?`, -1)                                  // 字符串
		expect = strings.Replace(expect, "%i", `[+\-]?[0-9]+`, -1)                         // 十进制数字，可有符号
		expect = strings.Replace(expect, "%d", `[0-9]+`, -1)                               // 十进制数字，无符号
		expect = strings.Replace(expect, "%x", `[0-9a-fA-F]+`, -1)                         // 十六进制数字
		expect = strings.Replace(expect, "%f", `[+\-]?\.?[0-9]+\.?[0-9]*(E-?[0-9]+)?`, -1) // 十进制浮点数
		expect = strings.Replace(expect, "%c", ".", -1)                                    // 单个字符
	} else if langType == misc.GO.String() {
		expect = strings.Replace(expect, "%s", `.+?`, -1)                                  // 字符串
		expect = strings.Replace(expect, "%v", `.+?`, -1)                                  // 字符串
		expect = strings.Replace(expect, "%d", `[0-9]+`, -1)                               // 十进制数字，无符号
		expect = strings.Replace(expect, "%t", `(true|false)+`, -1)                        // true或false
		expect = strings.Replace(expect, "%f", `[+\-]?\.?[0-9]+\.?[0-9]*(E-?[0-9]+)?`, -1) // 浮点数
	}

	pass, _ := regexp.MatchString("^"+expect+"$", actual)
	return pass
}
