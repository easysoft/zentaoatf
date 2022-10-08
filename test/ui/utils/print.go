package utils

import (
	"github.com/easysoft/zentaoatf/test/ui/conf"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"log"
)

func PrintErrOrNot(err error, t provider.T) {
	if err != nil {
		if conf.ExitAllOnError {
			log.Panicln(err)
		} else {
			log.Println(err)
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
