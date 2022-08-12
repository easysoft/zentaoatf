package main

/**

cid=0
pid=0

1.列出目录下的所有脚本 >> Success
2.列出目录下cid=1的脚本 >> Success
3.列出目录下标题包含match的脚本 >> Success

*/
import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"time"

	expect "github.com/google/goexpect"
)

var (
	newline = "\n"
)

func testLs(cmd string, successRe *regexp.Regexp) {
	if runtime.GOOS == "windows" {
		cmd = strings.ReplaceAll(cmd, "/", "\\")
	}
	child, _, err := expect.Spawn(cmd, -1)
	if err != nil {
		fmt.Println(err)
	}
	defer child.Close()

	if _, _, err = child.Expect(successRe, 10*time.Second); err != nil {
		fmt.Printf("%s: %s%s", successRe, err, newline)
		return
	}

	fmt.Println("Success")
}

func main() {
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}
	testLs("ztf list ./demo", regexp.MustCompile("Found 3 test cases|发现3个用例"))
	testLs("ztf ls ./demo -k 1", regexp.MustCompile("Found 1 test cases|发现1个用例"))
	testLs("ztf ls demo -k match", regexp.MustCompile("Found 2 test cases|发现2个用例"))
}
