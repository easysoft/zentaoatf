package main

import (
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
)

func main() {
	config := configUtils.ReadCurrConfig()

	commonUtils.SetFieldVal(config, "php", "sdfdsf")
}
