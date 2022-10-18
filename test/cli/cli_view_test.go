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
	"testing"
	"time"

	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type ViewSuite struct {
	suite.Suite
}

func (s *ViewSuite) BeforeEach(t provider.T) {
	t.ID("1593")
	t.AddSubSuite("命令行-查看脚本详情")
}
func (s *ViewSuite) TestViewSuite(t provider.T) {
	t.Require().Equal("Success", testView(commonTestHelper.GetZtfPath()+fmt.Sprintf(" view %stest/demo/1_string_match_fail.php", commonTestHelper.RootPath), regexp.MustCompile("check string matches pattern")))
	t.Require().Equal("Success", testView(commonTestHelper.GetZtfPath()+fmt.Sprintf(" -v %stest/demo -k 1", commonTestHelper.RootPath), regexp.MustCompile("check string matches pattern")))
	t.Require().Equal("Success", testView(commonTestHelper.GetZtfPath()+fmt.Sprintf(" view %stest/demo -k match", commonTestHelper.RootPath), regexp.MustCompile("Found 5 test cases|发现5个用例")))
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
