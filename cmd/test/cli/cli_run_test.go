package main

/**

$>ztf run demo/lang/bat                          执行目录bat下的脚本，支持多个目录和文件参数项。
$>ztf run product01 product01/all.cs             执行all.cs测试套件的用例，脚本在product01目录中。
$>ztf run log/001/result.txt                     执行result.txt结果文件中的失败用例。
$>ztf run product01 -suite 1                     执行禅道系统中编号为1的套件，脚本在product01目录，缩写-s。
$>ztf run -task 1                                执行禅道系统中编号为1的任务，脚本在当期目录, 缩写-t。
$>ztf run demo/demo -p 1 -t task1 -cr -cb        执行目录demo下的脚本，完成后提交结果到禅道，并将失败结果提交成缺陷。cr提交结果，-cb提交缺陷; -p必填参数指定产品ID， -t可选参数指定禅道新建测试单名称。
$>ztf run demo/autoit                            执行ZTF自带AutoIT脚本。
$>ztf run demo/selenium/chrome.php --interp runtime/php/php7/php.exe 执行ZTF自带Selenium脚本，使用指定PHP解释器。
$>ztf run demo/appium/android.php --interp runtime/php/php7/php.exe 执行ZTF自带Appium脚本，使用指定PHP解释器

*/
import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type RunSuite struct {
	suite.Suite
}

func (s *RunSuite) BeforeEach(t provider.T) {
	if runtime.GOOS == "windows" {
		os.RemoveAll(fmt.Sprintf("%s\\test\\demo\\php\\product1", constTestHelper.RootPath))
	} else {
		os.RemoveAll(fmt.Sprintf("%s/test/demo/php/product1", constTestHelper.RootPath))
	}
	t.AddSubSuite("命令行-run")
}

func (s *RunSuite) TestRunZtfFile(t provider.T) {
	t.ID("1584")
	t.Title("执行多个文件和目录中的脚本")

	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %s/test/demo/1_string_match_pass.php", constTestHelper.RootPath)
	expectReg := regexp.MustCompile(`Pass:1\(100\.0%\), Fail:0\(0\.0%\), Skip:0\(0\.0%\)|通过数：1\(100\.0%\)，失败数：0\(0\.0%\)，忽略数：0\(0\.0%\)`)

	t.Require().Equal("Success", testRun(cmd, expectReg))

	cmd = commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %stest/demo/1_string_match_pass.php %stest/demo/2_webpage_extract.php", constTestHelper.RootPath, constTestHelper.RootPath)
	expectReg = regexp.MustCompile(`Pass:2\(100\.0%\), Fail:0\(0\.0%\), Skip:0\(0\.0%\)|通过数：2\(100\.0%\)，失败数：0\(0\.0%\)，忽略数：0\(0\.0%\)`)

	t.Require().Equal("Success", testRun(cmd, expectReg))

	cmd = commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %stest/demo", constTestHelper.RootPath)
	if runtime.GOOS == "windows" {
		cmd = strings.ReplaceAll(cmd, "/", "\\")
	}
	expectReg = regexp.MustCompile(`Run \d+ scripts in \d+ sec|执行\d+个用例，耗时\d+秒`)

	t.Require().Equal("Success", testRun(cmd, expectReg))
}

func (s *RunSuite) TestRunExpectFile(t provider.T) {
	t.ID("7561")
	t.Title("执行期待结果独立文件的用例")

	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %s/test/demo/expect.php", constTestHelper.RootPath)
	expectReg := regexp.MustCompile(`Pass:1\(100\.0%\), Fail:0\(0\.0%\), Skip:0\(0\.0%\)|通过数：1\(100\.0%\)，失败数：0\(0\.0%\)，忽略数：0\(0\.0%\)`)

	t.Require().Equal("Success", testRun(cmd, expectReg))
}

func (s *RunSuite) TestRunFileAndSubmit(t provider.T) {
	t.ID("7552")
	t.Title("执行后自动提交结果")

	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %s/test/demo/1_string_match_pass.php -p 1 -cr", constTestHelper.RootPath)
	expectReg := regexp.MustCompile(`Submitted test results to ZenTao|提交测试结果到禅道成功`)

	t.Require().Equal("Success", testRun(cmd, expectReg))
}

func (s *RunSuite) TestRunFileAndSubmitBug(t provider.T) {
	t.ID("7553")
	t.Title("执行后自动提交缺陷")

	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %s/test/demo/1_string_match_fail.php -p 1 -cb", constTestHelper.RootPath)
	expectReg := regexp.MustCompile(`Success to report bug for case \\d+|成功为用例\d+提交缺陷`)

	t.Require().Equal("Success", testRun(cmd, expectReg))
}

func (s *RunSuite) TestRunZtfTask(t provider.T) {
	t.ID("1589")
	t.Title("执行禅道测试任务")

	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %stest/demo -task 1", constTestHelper.RootPath)
	expectReg := regexp.MustCompile(`Pass:0\(0\.0%\), Fail:3\(100\.0%\), Skip:0\(0\.0%\)|通过数：0\(0\.0%\)，失败数：3\(100\.0%\)，忽略数：0\(0\.0%\)`)

	t.Require().Equal("Success", testRun(cmd, expectReg))
}

func (s *RunSuite) TestRunZtfSuite(t provider.T) {
	t.ID("1588")
	t.Title("执行禅道测试套件")

	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %stest%sdemo -suite 1", constTestHelper.RootPath, constTestHelper.FilePthSep)
	expectReg := regexp.MustCompile(`Pass:0\(0\.0%\), Fail:1\(100\.0%\), Skip:0\(0\.0%\)|通过数：0\(0\.0%\)，失败数：1\(100\.0%\)，忽略数：0\(0\.0%\)`)

	t.Require().Equal("Success", testRun(cmd, expectReg))
}

func (s *RunSuite) TestRunZtfCsFile(t provider.T) {
	t.ID("1586")
	t.Title("执行本地套件文件中指定编号的脚本")

	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %stest/demo %stest/demo/all.cs", constTestHelper.RootPath, constTestHelper.RootPath)
	expectReg := regexp.MustCompile(`Pass:0\(0\.0%\), Fail:2\(100\.0%\), Skip:0\(0\.0%\)|通过数：0\(0\.0%\)，失败数：2\(100\.0%\)，忽略数：0\(0\.0%\)`)

	t.Require().Equal("Success", testRun(cmd, expectReg))
}

func testRun(cmd string, expectReg *regexp.Regexp) string {
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err := child.Expect(expectReg, 10*time.Second); err != nil {
		return fmt.Sprintf("cmd:%s, expect %s, actual %s", cmd, expectReg, err.Error())
	}

	return "Success"
}

func (s *RunSuite) TestRunScenes(t provider.T) {
	sceneMap := map[string]string{
		"exactly empty >> ~~":                            "",
		"exactly start with abc >> ~f:^abc~":             "abcdvd",
		"exactly end with abc >> ~f:abc$~":               "dcdabc",
		"exactly contain abc >> ~f:abc~":                 "dvabcd",
		"exactly containX abc*3 >> ~f:abc*3~":            "dvabcdabcabcdds",
		"exactly 2 in (1,2,3) >> ~f:(1,2,3)~":            "2",
		"exactly match %sabc%d >> ~m:%sabc%d~":           "dabc1dad",
		"exactly match .*cid=.* >> ~m:.*cid=.*~":         "sdfascid=ljlkjl",
		"exactly equal 123 >> ~c:=123~":                  "123",
		"exactly less than 123 >> ~c:<123~":              "1",
		"exactly less than or equal 123 >> ~c:<=123~":    "123",
		"exactly greater than 123 >> ~c:>123~":           "124",
		"exactly greater than or equal 123 >> ~c:>=123~": "123",
		"exactly not equal 123 >> ~c:<>123~":             "120",
		"exactly between 12-19 >> ~c:12-19~":             "15",
	}
	template := `#!/usr/bin/env php
	<?php
	/**

	title=check string matches pattern
	cid=1
	pid=1

	%s

	*/

	print("%s\n");`

	path := "../../demo/test_scene.php"
	if runtime.GOOS == "windows" {
		path = `..\..\demo\test_scene.php`
	}
	cmd := commonTestHelper.GetZtfPath() + ` run ` + path

	for expectVal, actualVal := range sceneMap {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			t.Errorf("open file fail, err:%s", err)
			return
		}

		write := bufio.NewWriter(file)
		write.WriteString(fmt.Sprintf(template, expectVal, actualVal))
		write.Flush()
		file.Close()

		t.Require().Equal("Success", testRun(cmd, regexp.MustCompile(`Pass:1\(100\.0%\), Fail:0\(0\.0%\), Skip:0\(0\.0%\)|通过数：1\(100\.0%\)，失败数：0\(0\.0%\)，忽略数：0\(0\.0%\)`)))
	}
	os.Remove(path)
}

func TestCliRun(t *testing.T) {
	suite.RunSuite(t, new(RunSuite))
}
