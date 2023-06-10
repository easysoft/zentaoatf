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

type ViewSuite struct {
	suite.Suite
}

func (s *ViewSuite) BeforeEach(t provider.T) {
	t.ID("1593")
	commonTestHelper.ReplaceLabel(t, "命令行-查看脚本详情")
}
func (s *ViewSuite) TestViewSuite(t provider.T) {
	t.Require().Equal("Success", testView(commonTestHelper.GetZtfPath()+fmt.Sprintf(" view %scmd/test/demo/1_string_match_fail.php", constTestHelper.RootPath), regexp.MustCompile("check string matches pattern")))

	t.Require().Equal("Success", testView(commonTestHelper.GetZtfPath()+fmt.Sprintf(" -v %scmd/test/demo -k 1", constTestHelper.RootPath), regexp.MustCompile("check string matches pattern")))

	t.Require().Equal("Success", testView(commonTestHelper.GetZtfPath()+fmt.Sprintf(" view %scmd/test/demo -k match", constTestHelper.RootPath), regexp.MustCompile("Found 5 test cases|发现5个用例")))
}

func testView(cmd string, successRe *regexp.Regexp) string {
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

func TestCliView(t *testing.T) {
	suite.RunSuite(t, new(ViewSuite))
}
