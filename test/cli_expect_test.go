package main

/**

cid=0
pid=0

1.更新用例到禅道 >> Success

*/
import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/stretchr/testify/suite"
)

var (
	successExpectRe = regexp.MustCompile("Success to create independent expect results file|成功创建独立的期待结果文件")
)

type ExpectSuit struct {
	suite.Suite
	testCount uint32
}

func (s *ExpectSuit) TestExpectSuite() {
	assert.Equal(s.Suite.T(), "Success", testExpect())
}

func testExpect() string {
	path := `../demo/sample/1_simple.php`
	if runtime.GOOS == "windows" {
		path = `..\demo\sample\1_simple.php`
	}
	cmd := `ztf expect ` + path

	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err = child.Expect(successExpectRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successExpectRe, err.Error())
	}

	expPath := path[:len(path)-3] + "exp"
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

func TestExpect(t *testing.T) {
	suite.Run(t, new(ExpectSuit))
}
