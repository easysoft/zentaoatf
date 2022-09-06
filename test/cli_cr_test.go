package main

/**

cid=0
pid=0
timeout=10

1.提交结果到禅道 >> Success

*/
import (
	"fmt"
	"regexp"
	"runtime"
	"testing"
	"time"

	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	successCrRe = regexp.MustCompile("Submitted test results to ZenTao|提交测试结果到禅道成功")
)

type CrSuite struct {
	suite.Suite
}

func (s *CrSuite) BeforeEach(t provider.T) {
	t.ID("1590")
	t.AddSubSuite("命令行-提交测试结果到禅道")
}
func (s *CrSuite) TestCrSuite(t provider.T) {
	t.Require().Equal("Success", testCr())
}

func testCr() string {
	cmd := "ztf cr demo/001 -p 1 -y -t testcr"
	if runtime.GOOS == "windows" {
		cmd = "ztf cr demo\\001 -p 1 -y -t testcr"
	}
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err = child.Expect(successCrRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successCrRe, err.Error())
	}
	return "Success"
}

func TestCr(t *testing.T) {
	suite.RunSuite(t, new(CrSuite))
}
