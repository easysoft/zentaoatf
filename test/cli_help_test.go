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

	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	successHelpRe = regexp.MustCompile("为了方便在任意目录中执行ztf.exe命令")
)

type HelpSuite struct {
	suite.Suite
}

func (s *HelpSuite) BeforeEach(t provider.T) {
	t.ID("1578")
	t.AddSubSuite("命令行-查看帮助")
}
func (s *HelpSuite) TestHelpSuite(t provider.T) {
	t.Require().Equal("Success", testHelp())
}

func testHelp() string {
	cmd := `ztf -h`

	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err = child.Expect(successHelpRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successHelpRe, err.Error())
	}

	return "Success"
}

func TestHelp(t *testing.T) {
	suite.RunSuite(t, new(HelpSuite))
}
