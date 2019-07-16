package main

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/misc"
)

func init() {
	p := misc.GetI118() // for test
	p.Printf("HELLO_1", "Peter")
	fmt.Println(p.Sprintf("HELLO_1", "Peter"))
}
