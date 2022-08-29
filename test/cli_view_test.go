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
	"testing"
	"time"

	"github.com/bmizerany/assert"
	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/stretchr/testify/suite"
)

type ViewSuit struct {
	suite.Suite
	testCount uint32
}

func (s *ViewSuit) TestViewSuite() {
	assert.Equal(s.Suite.T(), "Success", testView("ztf view ./demo/1_string_match_fail.php", regexp.MustCompile("check string matches pattern")))
	assert.Equal(s.Suite.T(), "Success", testView("ztf -v ./demo -k 1", regexp.MustCompile("check string matches pattern")))
	assert.Equal(s.Suite.T(), "Success", testView("ztf view demo -k match", regexp.MustCompile("Found 3 test cases|发现3个用例")))
}

func testView(cmd string, successRe *regexp.Regexp) string {
	if runtime.GOOS == "windows" {
		cmd = strings.ReplaceAll(cmd, "/", "\\")
	}
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err = child.Expect(successRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successRe, err.Error())
	}

	return "Success"
}

func TestView(t *testing.T) {
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}
	suite.Run(t, new(ViewSuit))
}
