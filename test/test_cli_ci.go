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

	expect "github.com/google/goexpect"
)

var (
	continueRE = regexp.MustCompile("Will commit cases below to Zentao|以下用例信息将被更新到禅道")
	successRE  = regexp.MustCompile("Totally commit 1 cases to Zentao|合计更新1个用例到禅道")
	newline    = "\n"
)

func testCi() {
	cmd := "ztf ci ./demo/1_string_match_pass.php"
	if runtime.GOOS == "windows" {
		cmd = "ztf ci .\\demo\\1_string_match_pass.php"
	}
	child, _, err := expect.Spawn(cmd, -1)
	if err != nil {
		fmt.Println(err)
	}
	defer child.Close()
	if _, _, err = child.Expect(continueRE, time.Second); err != nil {
		fmt.Printf("%s: %s%s", continueRE, err, newline)
		return
	}

	if err = child.Send("y" + newline); err != nil {
		fmt.Println(err)
		return
	}

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
