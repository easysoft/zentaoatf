package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	successExpectRe = regexp.MustCompile("Success to create independent expect results file|成功创建独立的期待结果文件")
)

type ExpectSuite struct {
	suite.Suite
}

func (s *ExpectSuite) BeforeEach(t provider.T) {
	t.ID("5429")
	t.AddSubSuite("命令行-生成独立的期待结果文件")
}
func (s *ExpectSuite) TestExpectSuite(t provider.T) {
	t.Require().Equal("Success", testExpect())
}

func testExpect() string {
	path := fmt.Sprintf(`%sdemo/sample/1_simple.php`, constTestHelper.RootPath)
	expPath := path[:len(path)-3] + "exp"
	cmd := commonTestHelper.GetZtfPath() + ` expect ` + path

	os.Remove(expPath)

	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err = child.Expect(successExpectRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successExpectRe, err.Error())
	}

	file, err := os.Open(expPath)
	if err != nil {
		return err.Error()
	}
	defer func() {
		file.Close()
		os.Remove(expPath)
	}()

	content, err := ioutil.ReadAll(file)
	checkResSuccess := strings.Contains(string(content), `expect 1
pass
expect 3`)
	if !checkResSuccess {
		return "Check exp error"
	}

	return "Success"
}

func TestCliExpect(t *testing.T) {
	suite.RunSuite(t, new(ExpectSuite))
}
