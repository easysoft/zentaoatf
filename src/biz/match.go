package biz

import (
	"misc"
	"regexp"
	"strings"
)

func MatchString(expect string, actual string, langType string) bool {
	if langType == misc.PHP.String() {
		// $wantedReg = preg_quote($wanted, '/');
		expect = strings.Replace(expect, "%s", `.+?`, -1)
		expect = strings.Replace(expect, "%i", `[+\-]?[0-9]+`, -1)
		expect = strings.Replace(expect, "%d", `[0-9]+`, -1)
		expect = strings.Replace(expect, "%x", `[0-9a-fA-F]+`, -1)
		expect = strings.Replace(expect, "%f", `[+\-]?\.?[0-9]+\.?[0-9]*(E-?[0-9]+)?`, -1)
		expect = strings.Replace(expect, "%c", ".", -1)
	} else if langType == misc.GO.String() {
		expect = strings.Replace(expect, "%s", `.+?`, -1)
		expect = strings.Replace(expect, "%v", `.+?`, -1)
		expect = strings.Replace(expect, "%d", `[0-9]+`, -1)
		expect = strings.Replace(expect, "%t", `(true|false)+`, -1)
		expect = strings.Replace(expect, "%f", `[+\-]?\.?[0-9]+\.?[0-9]*(E-?[0-9]+)?`, -1)
	}

	pass, _ := regexp.MatchString(expect, actual)
	return pass
}
