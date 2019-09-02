package main

import (
	"fmt"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
)

func main() {
	osName := commonUtils.GetOs()
	fmt.Println(osName)
}
