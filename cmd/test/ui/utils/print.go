package utils

import (
	"log"

	"github.com/easysoft/zentaoatf/cmd/test/ui/conf"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func PrintErrOrNot(err error, t provider.T) {
	if err != nil && t != nil {
		if conf.ExitAllOnError {
			t.Error(err)
			t.FailNow()
		} else if conf.ShowErr {
			t.Error(err)
		}
	}
}

func PrintErrMsg(err string, t provider.T) {
	if err != "" {
		if conf.ExitAllOnError {
			log.Panicln(err)
		} else {
			log.Println(err)
		}
	}
}
