package main

import (
	"flag"
	"fmt"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
)

func main() {
	var languages []string
	flag.Var(commonUtils.NewSliceValue([]string{}, &languages), "slice", "I like programming `languages`")
	flag.Parse()

	fmt.Println(languages)
}
