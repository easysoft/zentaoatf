package main

import (
	"flag"
	"fmt"
	"testing"

	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
)

func init() {
	var version = flag.String("zentaoVersion", "", "")
	// for _, arg := range os.Args {
	// 	if strings.Contains(arg, "zentaoVersion") {
	// 		*version = arg[strings.Index(arg, "=")+1:]
	// 	}
	// }
	testing.Init()
	flag.Parse()
	fmt.Println(1111, commonTestHelper.NewLine, commonTestHelper.RootPath, *version)
	// err := commonTestHelper.InitZentao(*version)
	// if err != nil {
	// 	fmt.Println("Init zentao data fail ", err)
	// }
	// err = commonTestHelper.Pull()
	// if err != nil {
	// 	fmt.Println("Git pull code fail ", err)
	// }
	// err = commonTestHelper.BuildCli()
	// if err != nil {
	// 	fmt.Println("Build cli fail ", err)
	// }
	// err = commonTestHelper.BuildServer()
	// if err != nil {
	// 	fmt.Println("Build server fail ")
	// }
}

func main() {

}
