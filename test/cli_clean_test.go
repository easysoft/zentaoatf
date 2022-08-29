package main

/**

cid=0
pid=0

1.更新用例到禅道 >> Success

*/
import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/stretchr/testify/suite"
)

var (
	successCleanRe = regexp.MustCompile("Successfully cleaned all logs|成功删除所有日志")
)

type CleanSuit struct {
	suite.Suite
	testCount uint32
}

func (s *CleanSuit) TestCleanSuite() {
	assert.Equal(s.Suite.T(), "Success", testClean())
}

func testClean() string {
	cmd := `ztf clean`
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

func TestClean(t *testing.T) {
	suite.Run(t, new(CleanSuit))
}