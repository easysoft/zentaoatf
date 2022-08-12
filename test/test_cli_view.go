package main

/**

cid=0
pid=0

1.查看脚本 >> Success
2.查看目录下cid=1的脚本 >> Success
3.查看目录下标题包含match的脚本 >> Success

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
	testLs("ztf view ./demo/1_string_match_fail.php", regexp.MustCompile("check string matches pattern"))
	testLs("ztf -v ./demo -k 5", regexp.MustCompile("extract content from webpage"))
	testLs("ztf view demo -k match", regexp.MustCompile("Found 2 test cases|发现2个用例"))
}
