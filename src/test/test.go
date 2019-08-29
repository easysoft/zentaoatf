package main

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils/lang"
	"strings"
)

func main() {
	fmtParam := strings.Join(langUtils.GetSupportLangageArr(), " / ")
	scriptLang := fmt.Sprintf("enter_co_language%s", fmtParam)

	fmt.Println(scriptLang)
}
