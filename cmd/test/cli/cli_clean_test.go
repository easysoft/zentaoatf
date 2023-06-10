package main

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"testing"
	"time"

	commonTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/common"

	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	successCleanRe = regexp.MustCompile("Successfully cleaned all logs|成功删除所有日志")
)

type CleanSuite struct {
	suite.Suite
}

func (s *CleanSuite) BeforeEach(t provider.T) {
	t.ID("1595")
	commonTestHelper.ReplaceLabel(t, "命令行-清除测试日志")
}
func (s *CleanSuite) TestCleanSuite(t provider.T) {
	t.Require().Equal("Success", testClean())
}

func testClean() string {
	cmd := commonTestHelper.GetZtfPath() + ` clean`
	path := "./log/test"
	if runtime.GOOS == "windows" {
		path = `.\log\test`
	}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return "Mkdir fail"
			}
		}
	}

	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err = child.Expect(successCleanRe, 10*time.Second); err != nil {
		return fmt.Sprintf("clean %s, actual %s", successCleanRe, err.Error())
	}
	_, err = os.Stat(path)
	if err == nil {
		return "Clean fail"

	}

	return "Success"
}

func TestCliClean(t *testing.T) {
	suite.RunSuite(t, new(CleanSuite))
}
