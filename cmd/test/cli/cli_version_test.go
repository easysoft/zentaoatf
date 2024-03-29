package main

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	commonTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/common"
	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
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
	t.ID("7534")
	commonTestHelper.ReplaceLabel(t, "命令行-查看ZTF版本")
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
