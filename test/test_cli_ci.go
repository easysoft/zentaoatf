package main

/**

cid=0
pid=0

1.更新用例到禅道 >> Success

*/
import (
	"fmt"
	"regexp"
	"runtime"
	"time"

	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
)

var (
	continueRe = regexp.MustCompile("Will commit cases below to Zentao|以下用例信息将被更新到禅道")
	successRe  = regexp.MustCompile("Totally commit 1 cases to Zentao|合计更新1个用例到禅道")
	newline    = "\n"
)

func testCi() {
	cmd := "ztf ci ./demo/1_string_match_pass.php"
	if runtime.GOOS == "windows" {
		cmd = "ztf ci .\\demo\\1_string_match_pass.php"
	}
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		fmt.Println(err)
	}
	defer child.Close()
	if _, _, err = child.Expect(continueRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", continueRe, err, newline)
		return
	}

	if err = child.Send("y" + newline); err != nil {
		fmt.Println(err)
		return
	}

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
	testCi()
}
