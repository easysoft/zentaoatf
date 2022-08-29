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
	successExtractRe = regexp.MustCompile("Success to extract test steps and results|成功从注释提取步骤和期待结果")
)

type ExtractSuit struct {
	suite.Suite
	testCount uint32
}

func (s *ExtractSuit) TestExtractSuite() {
	assert.Equal(s.Suite.T(), "Success", testExtract())
}

func testExtract() string {
	path := `../demo/sample/8_extract_desc.php`
	if runtime.GOOS == "windows" {
		path = `..\demo\sample\8_extract_desc.php`
	}
	cmd := `ztf extract ` + path

	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err = child.Expect(successExtractRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successExtractRe, err.Error())
	}

	file, err := os.Open(path)
	if err != nil {
		return "File can not open"
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	checkResSuccess := strings.Contains(string(content), `/**

title=sync step from comments
timeout=0
cid=0
pid=0

1 >> expect 1

group2
  2.1 >> expect 2.1
  2.2 >> expect 2.2
  2.3 >> expect 2.3  

multi line expect >>
  expect 3.1
  expect 3.2
>>

4 >> expect 4
5 >> expect 5
step 6 >> expect 6
step 7 >> expect 7
step 8 >> expect 8
step 9 >> expect 9
step 10 >> expect 10

*/`)
	if !checkResSuccess {
		return "Check steps error"
	}
	return "Success"
}

func TestExtract(t *testing.T) {
	suite.Run(t, new(ExtractSuit))
}
