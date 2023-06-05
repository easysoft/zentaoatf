package main

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	commonTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	continueRe  = regexp.MustCompile("Which case do you want to report bug for|请输入您想提交缺陷的用例ID")
	successCbRe = regexp.MustCompile("Success to report bug for case \\d+|成功为用例\\d+提交缺陷")
)

type CbSuite struct {
	suite.Suite
}

func (s *CbSuite) BeforeEach(t provider.T) {
	t.ID("1591")
	t.AddSubSuite("命令行-提交失败结果为禅道中缺陷")
}
func (s *CbSuite) TestCbSuite(t provider.T) {
	t.Require().Equal("Success", testCb())
}

func testCb() string {
	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" cb %stest/demo/001 -p 1", constTestHelper.RootPath)
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()
	if _, err := child.Expect(continueRe, 5*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", continueRe, err.Error())
	}

	if err = child.Send("1" + constTestHelper.NewLine); err != nil {
		return err.Error()
	}

	if _, err = child.Expect(successCbRe, 30*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successCbRe, err.Error())
	}
	child.Close()
	return "Success"
}

func TestCliCb(t *testing.T) {
	suite.RunSuite(t, new(CbSuite))
}
