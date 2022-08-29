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
	"testing"
	"time"

	"github.com/bmizerany/assert"
	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/stretchr/testify/suite"
)

type ListSuit struct {
	suite.Suite
	testCount uint32
}

func (s *ListSuit) TestListSuite() {
	assert.Equal(s.Suite.T(), "Success", testLs("ztf list ./demo", regexp.MustCompile("Found 5 test cases|发现5个用例")))
	assert.Equal(s.Suite.T(), "Success", testLs("ztf ls ./demo -k 1", regexp.MustCompile("Found 2 test cases|发现2个用例")))
	assert.Equal(s.Suite.T(), "Success", testLs("ztf ls demo -k match", regexp.MustCompile("Found 3 test cases|发现3个用例")))
}

func testLs(cmd string, successRe *regexp.Regexp) string {
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

func TestList(t *testing.T) {
	suite.Run(t, new(ListSuit))
}
