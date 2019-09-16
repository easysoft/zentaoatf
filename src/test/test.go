package main

import (
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
)

func main() {

	_ = configUtils.SaveConfig("a", "2", "3", "4")

}
