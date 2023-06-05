package main

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	successCrRe = regexp.MustCompile("Submitted test results to ZenTao|提交测试结果到禅道成功")
	productIdRe = regexp.MustCompile("Zentao account|请输入 产品Id")
	taskIdRe    = regexp.MustCompile("Zentao account|请输入 测试任务Id")
)

type CrSuite struct {
	suite.Suite
}

func (s *CrSuite) BeforeEach(t provider.T) {
	t.ID("1590")
	t.AddSubSuite("命令行-提交测试结果到禅道")
}

func (s *CrSuite) TestAutoCr(t provider.T) {
	t.ID("7558")
	t.AddSubSuite("命令行-提交测试结果到禅道免确认")

	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" cr %stest/demo/001 -p 1 -y -t testcr", constTestHelper.RootPath)

	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		t.Require().Equal("Success", err.Error())
	}
	defer child.Close()

	if _, err = child.Expect(successCrRe, 10*time.Second); err != nil {
		t.Require().Equal("Success", fmt.Sprintf("expect %s, actual %s", successCrRe, err.Error()))
	}
}

func (s *CrSuite) TestCr(t provider.T) {
	t.ID("1590")
	t.AddSubSuite("命令行-提交测试结果到禅道")

	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" cr %stest/demo/001", constTestHelper.RootPath)
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		t.Require().Equal("Success", err.Error())
	}

	if _, err = child.Expect(productIdRe, time.Second*5); err != nil {
		t.Errorf("expect %s, actual %s", productIdRe, err.Error())
	}
	if err = child.Send("1" + constTestHelper.NewLine); err != nil {
		t.Error(err.Error())
	}

	if _, err = child.Expect(taskIdRe, time.Second*5); err != nil {
		t.Errorf("expect %s, actual %s", taskIdRe, err.Error())
	}
	if err = child.Send("1" + constTestHelper.NewLine); err != nil {
		t.Error(err.Error())
	}

	if _, err = child.Expect(successCrRe, 10*time.Second); err != nil {
		t.Require().Equal("Success", fmt.Sprintf("expect %s, actual %s", successCrRe, err.Error()))
	}

	defer child.Close()
}

func TestCliCr(t *testing.T) {
	suite.RunSuite(t, new(CrSuite))
}
