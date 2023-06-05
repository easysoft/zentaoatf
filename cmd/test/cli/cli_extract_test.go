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
	successExtractRe = regexp.MustCompile("Success to extract test steps and results|成功从注释提取步骤和期待结果")
)

type ExtractSuite struct {
	suite.Suite
}

func (s *ExtractSuite) BeforeEach(t provider.T) {
	t.ID("5430")
	t.AddSubSuite("命令行-提取分散在脚本中的注释")
}
func (s *ExtractSuite) TestExtractSuite(t provider.T) {
	t.Require().Equal("Success", testExtract())
}

func testExtract() string {
	path := fmt.Sprintf(`%sdemo/sample/8_extract_desc.php`, constTestHelper.RootPath)
	cmd := commonTestHelper.GetZtfPath() + ` extract ` + path

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
	checkResSuccess := strings.Contains(string(content), `
title=sync step from comments
timeout=0
cid=0

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
step 10 >> expect 10`)
	if !checkResSuccess {
		return "Check steps error"
	}
	return "Success"
}

func TestCliExtract(t *testing.T) {
	suite.RunSuite(t, new(ExtractSuite))
}
