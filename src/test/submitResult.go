package main

import (
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
)

func main() {
	zentaoService.Login("http://ztpmp.ngtesting.org/", "admin", "P2ssw0rd")
	testingService.SubmitResult("scripts/all.suite", "2019-08-15T082332")
}
