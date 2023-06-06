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

type ListSuite struct {
	suite.Suite
}

func (s *ListSuite) BeforeEach(t provider.T) {
	t.ID("1592")
	t.AddSubSuite("命令行-查看脚本列表")
}
func (s *ListSuite) TestListSuite(t provider.T) {
	t.Require().Equal("Success", testLs(commonTestHelper.GetZtfPath()+fmt.Sprintf(" list %scmd/test/demo", constTestHelper.RootPath), regexp.MustCompile("Found 7 test cases|发现7个用例")))
	t.Require().Equal("Success", testLs(commonTestHelper.GetZtfPath()+fmt.Sprintf(" ls %scmd/test/demo -k 1", constTestHelper.RootPath), regexp.MustCompile("Found 3 test cases|发现3个用例")))
	t.Require().Equal("Success", testLs(commonTestHelper.GetZtfPath()+fmt.Sprintf(" ls %scmd/test/demo -k match", constTestHelper.RootPath), regexp.MustCompile("Found 4 test cases|发现4个用例")))
}

func testLs(cmd string, successRe *regexp.Regexp) string {
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

func TestCliList(t *testing.T) {
	suite.RunSuite(t, new(ListSuite))
}
