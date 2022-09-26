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
	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	successVersionRe = regexp.MustCompile("Build TimeStamp")
)

type VersionSuite struct {
	suite.Suite
}

func (s *VersionSuite) BeforeEach(t provider.T) {
	t.ID("2")
	t.AddSubSuite("命令行-查看ZTF版本")
}
func (s *VersionSuite) TestVersionSuite(t provider.T) {
	t.Require().Equal("Success", testVersion())
}

func testVersion() string {
	cmd := commonTestHelper.GetZtfPath() + ` version`

	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err = child.Expect(successVersionRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successVersionRe, err.Error())
	}

	return "Success"
}

func TestCliVersion(t *testing.T) {
	suite.RunSuite(t, new(VersionSuite))
}
