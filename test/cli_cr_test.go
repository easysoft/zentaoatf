package main

/**

cid=0
pid=0
timeout=10

1.提交结果到禅道 >> Success

*/
import (
	"fmt"
	"regexp"
	"runtime"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/stretchr/testify/suite"
)

var (
	successCrRe = regexp.MustCompile("Submitted test results to ZenTao|提交测试结果到禅道成功")
)

type CrSuit struct {
	suite.Suite
	testCount uint32
}

func (s *CrSuit) TestCrSuite() {
	assert.Equal(s.Suite.T(), "Success", testCr())
}

func testCr() string {
	cmd := "ztf cr demo/001 -p 1 -y -t testcr"
	if runtime.GOOS == "windows" {
		cmd = "ztf cr demo\\001 -p 1 -y -t testcr"
	}
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err = child.Expect(successCrRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successCrRe, err.Error())
	}
	return "Success"
}

func TestCr(t *testing.T) {
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}
	suite.Run(t, new(CrSuit))
}
