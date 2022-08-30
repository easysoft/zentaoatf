package main

/**

cid=0
pid=0

1.提交bug到禅道 >> Success

*/
import (
	"fmt"
	"regexp"
	"runtime"
	"testing"
	"time"

	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	continueRe  = regexp.MustCompile("Which case do you want to report bug for|请输入您想提交缺陷的用例ID")
	successCbRe = regexp.MustCompile("Success to report bug for case \\d+|成功为用例\\d+提交缺陷")
)

type CbSuit struct {
	suite.Suite
	testCount uint32
}

func (s *CbSuit) TestCbSuite() {
	assert.Equal(s.Suite.T(), "Success", testCb())
}

func testCb() string {
	cmd := "ztf cb demo/001 -p 1"
	if runtime.GOOS == "windows" {
		cmd = "ztf cb demo\\001 -p 1"
	}
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()
	if _, err := child.Expect(continueRe, 2*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", continueRe, err.Error())
	}

	if err = child.Send("1" + newline); err != nil {
		return err.Error()
	}

	if _, err = child.Expect(successCbRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successCbRe, err.Error())
	}
	child.Close()
	return "Success"
}

func TestCb(t *testing.T) {
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}
	suite.Run(t, new(CbSuit))
}
