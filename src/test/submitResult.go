package main

import (
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
)

func main() {
	zentaoService.Login("http://ztpmp.ngtesting.org/", "admin", "P2ssw0rd")
	zentaoService.SubmitResult("scripts/all.suite", "2019-08-15T170425")
}
