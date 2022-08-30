package main

/**

cid=0
pid=0

1.更新用例到禅道 >> Success

*/
import (
	"fmt"
	"regexp"
	"runtime"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/stretchr/testify/suite"
)

var (
	continueCiRe = regexp.MustCompile("Will commit cases below to Zentao|以下用例信息将被更新到禅道")
	successCiRe  = regexp.MustCompile("Totally commit 1 cases to Zentao|合计更新1个用例到禅道")
)

type CiSuit struct {
	suite.Suite
	testCount uint32
}

func (s *CiSuit) TestCiSuite() {
	assert.Equal(s.Suite.T(), "Success", testCi())
}

func testCi() string {
	cmd := "ztf ci ./demo/1_string_match_pass.php"
	if runtime.GOOS == "windows" {
		cmd = "ztf ci .\\demo\\1_string_match_pass.php"
	}
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()
	if _, err = child.Expect(continueCiRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", continueCiRe, err.Error())
	}

	if err = child.Send("y" + newline); err != nil {
		return err.Error()
	}

	if _, err = child.Expect(successCiRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successCiRe, err.Error())
	}

	return "Success"
}

func TestCi(t *testing.T) {
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}
	suite.Run(t, new(CiSuit))
}
