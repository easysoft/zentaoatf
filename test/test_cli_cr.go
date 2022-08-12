package main

/**

cid=0
pid=0

1.提交结果到禅道 >> Success

*/
import (
	"fmt"
	"regexp"
	"runtime"
	"time"

	expect "github.com/google/goexpect"
)

var (
	successRE = regexp.MustCompile("Submitted test results to ZenTao|提交测试结果到禅道成功")
	newline   = "\n"
)

func testCi() {
	cmd := "ztf cr demo/001 -p 1 -y -t testcr"
	if runtime.GOOS == "windows" {
		cmd = "ztf cr demo\\001 -p 1 -y -t testcr"
	}
	child, _, err := expect.Spawn(cmd, -1)
	if err != nil {
		fmt.Println(err)
	}
	defer child.Close()

	if _, _, err = child.Expect(successRE, 10*time.Second); err != nil {
		fmt.Printf("%s: %s%s", successRE, err, newline)
		return
	}

	fmt.Println("Success")
}

func main() {
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}
	testCi()
}
