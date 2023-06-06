package main

/**

cid=0
pid=0

1.更新用例到禅道 >> Success

*/
import (
	"fmt"
	"regexp"
	"testing"
	"time"

	commonTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	apiTest "github.com/easysoft/zentaoatf/cmd/test/helper/zentao/api"
	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	continueCiRe = regexp.MustCompile("Will commit cases below to Zentao|以下用例信息将被更新到禅道")
	successCiRe  = regexp.MustCompile("Totally commit 1 cases to Zentao|合计更新1个用例到禅道")
)

type CiSuite struct {
	suite.Suite
}

func (s *CiSuite) BeforeEach(t provider.T) {
	t.ID("5431")
	t.AddSubSuite("命令行-同步用例信息到禅道")
}
func (s *CiSuite) TestCiSuite(t provider.T) {
	t.Require().Equal("Success", testCi())
}

func testCi() string {
	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" ci %scmd/test/demo/1_string_match_pass.php", constTestHelper.RootPath)
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err = child.Expect(continueCiRe, time.Second*10); err != nil {
		return fmt.Sprintf("expect %s, actual %s", continueCiRe, err.Error())
	}

	if err = child.Send("y" + constTestHelper.NewLine); err != nil {
		return err.Error()
	}

	if _, err = child.Expect(successCiRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successCiRe, err.Error())
	}

	//check zentao info
	title := apiTest.GetCaseTitleById(1)
	if title != "check string matches pattern" {
		return fmt.Sprintf("check zentao title fail, expect check string matches pattern, actual %s", title)
	}

	return "Success"
}

func TestCliCi(t *testing.T) {
	suite.RunSuite(t, new(CiSuite))
}
