package main

import (
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/vari"
)

func main() {
	zentaoService.Login("http://ztpmp.ngtesting.org/", "admin", "P2ssw0rd")

	vari.CurrScriptFile = "scripts/all.suite"
	vari.CurrResultDate = "2019-08-19T144559"
	zentaoService.SubmitResult()
}
