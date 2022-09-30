package utils

import (
	"github.com/easysoft/zentaoatf/test/ui/conf"
	"log"
)

func PrintErr(err error) {
	if err != nil && conf.ExitOnError {
		log.Panic(err)
	} else {
		log.Println(err)
	}
}
