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
$>ztf junit -p 1 mvn clean package test          执行junit单元测试脚本，

cid=0
pid=0

1.co 运行单个脚本 >> Success
2.co 运行all.cs >> Success
3.co 根据日志重新运行失败用例 >> Success
4.co 运行套件 >> Success
5.co 运行任务 >> Success
6.co 运行目录并提交结果与bug >> Success

*/
import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	"github.com/go-git/go-git/v5"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type RunSuit struct {
	suite.Suite
	testCount uint32
}

func (s *RunSuit) BeforeEach(t provider.T) {
	if runtime.GOOS == "windows" {
		os.RemoveAll(fmt.Sprintf("%s\\test\\demo\\php\\product1", commonTestHelper.RootPath))
	} else {
		os.RemoveAll(fmt.Sprintf("%s/test/demo/php/product1", commonTestHelper.RootPath))
	}
	t.AddSubSuite("命令行-run")
}
func (s *RunSuit) TestRunZtfFile(t provider.T) {
	t.ID("1584")
	t.Title("执行多个文件和目录中的脚本")
	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %s/test/demo/1_string_match_pass.php", commonTestHelper.RootPath)
	expectReg := regexp.MustCompile(`Run 1 scripts in \d+ sec, 1\(100\.0%\) Pass, 0\(0\.0%\) Fail, 0\(0\.0%\) Skip|执行1个用例，耗时\d+秒。1\(100\.0%\) 通过，0\(0\.0%\) 失败，0\(0\.0%\) 忽略`)
	t.Require().Equal("Success", testRun(cmd, expectReg))

	cmd = commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %stest/demo/1_string_match_pass.php %stest/demo/2_webpage_extract.php", commonTestHelper.RootPath, commonTestHelper.RootPath)
	expectReg = regexp.MustCompile(`Run 2 scripts in \d+ sec, 2\(100\.0%\) Pass, 0\(0\.0%\) Fail, 0\(0\.0%\) Skip|执行2个用例，耗时\d+秒。2\(100\.0%\) 通过，0\(0\.0%\) 失败，0\(0\.0%\) 忽略`)
	t.Require().Equal("Success", testRun(cmd, expectReg))

	cmd = commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %stest/demo", commonTestHelper.RootPath)
	if runtime.GOOS == "windows" {
		cmd = strings.ReplaceAll(cmd, "/", "\\")
	}
	expectReg = regexp.MustCompile(`Run \d+ scripts in \d+ sec|执行\d+个用例，耗时\d+秒`)
	t.Require().Equal("Success", testRun(cmd, expectReg))
}

func (s *RunSuit) TestRunZtfTask(t provider.T) {
	t.ID("1589")
	t.Title("执行禅道测试任务")
	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %stest/demo -task 1", commonTestHelper.RootPath)
	expectReg := regexp.MustCompile(`Run 3 scripts in \d+ sec, 1\(33\.0%\) Pass, 2\(66\.0%\) Fail, 0\(0\.0%\) Skip|执行3个用例，耗时\d+秒。1\(33\.0%\) 通过，2\(66\.0%\) 失败，0\(0\.0%\) 忽略`)
	t.Require().Equal("Success", testRun(cmd, expectReg))
}

func (s *RunSuit) TestRunZtfSuite(t provider.T) {
	t.ID("1588")
	t.Title("执行禅道测试套件")
	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %stest/demo -suite 1", commonTestHelper.RootPath)
	expectReg := regexp.MustCompile(`Run 1 scripts in \d+ sec, 0\(0\.0%\) Pass, 1\(100\.0%\) Fail, 0\(0\.0%\) Skip|执行1个用例，耗时\d+秒。0\(0\.0%\) 通过，1\(100\.0%\) 失败，0\(0\.0%\) 忽略`)
	t.Require().Equal("Success", testRun(cmd, expectReg))
}

func (s *RunSuit) TestRunZtfCsFile(t provider.T) {
	t.ID("1586")
	t.Title("执行本地套件文件中指定编号的脚本")
	cmd := commonTestHelper.GetZtfPath() + fmt.Sprintf(" run %stest/demo %stest/demo/all.cs", commonTestHelper.RootPath, commonTestHelper.RootPath)
	expectReg := regexp.MustCompile(`Run 2 scripts in \d+ sec, 1\(50\.0%\) Pass, 1\(50\.0%\) Fail, 0\(0\.0%\) Skip|执行2个用例，耗时\d+秒。1\(50\.0%\) 通过，1\(50\.0%\) 失败，0\(0\.0%\) 忽略`)
	t.Require().Equal("Success", testRun(cmd, expectReg))
}

func (s *RunSuit) TestRunTestng(t provider.T) {
	testngDir := fmt.Sprintf("%stest/demo/ci_test_testng", commonTestHelper.RootPath)
	t.ID("5432")
	t.Title("执行TestNG单元测试")
	cloneGit("https://gitee.com/ngtesting/ci_test_testng.git", testngDir)
	t.Require().Equal("Success", testRunUnitTest("mvn clean package test", testngDir, regexp.MustCompile(`Tests run\: 3, Failures\: 0, Errors\: 0, Skipped\: 0`)))
}

func (s *RunSuit) TestRunPytest(t provider.T) {
	pytestDir := fmt.Sprintf(".%stest/demo/ci_test_pytest", commonTestHelper.RootPath)
	t.ID("5435")
	t.Title("执行PyTest单元测试")
	cloneGit("https://gitee.com/ngtesting/ci_test_pytest.git", pytestDir)

	t.Require().Equal("Success", testRunUnitTest("pytest --junitxml=testresult.xml", pytestDir, regexp.MustCompile("1 failed, 1 passed")))
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

func cloneGit(gitUrl string, name string) error {
	projectDir := name
	fileUtils.MkDirIfNeeded(projectDir)

	options := git.CloneOptions{
		URL:      gitUrl,
		Progress: os.Stdout,
	}
	_, err := git.PlainClone(projectDir, false, &options)
	return err
}

func testRunUnitTest(cmdStr, workspacePath string, successRe *regexp.Regexp) string {

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	cmd.Dir = workspacePath

	if cmd == nil {
		return "cmd is nil"
	}
	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()

	if err1 != nil {
		return err1.Error()
	} else if err2 != nil {
		return err2.Error()
	}
	cmd.Start()

	isTerminal := false
	reader1 := bufio.NewReader(stdout)
	for {
		line, err3 := reader1.ReadString('\n')
		if line != "" {
			isTerminal = true
			if successRe.MatchString(line) {
				return "Success"
			}
		}

		if err3 != nil || io.EOF == err3 {
			break
		}

	}

	errOutputArr := make([]string, 0)
	if !isTerminal {
		reader2 := bufio.NewReader(stderr)

		for {
			line, err2 := reader2.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				break
			}
			errOutputArr = append(errOutputArr, line)
		}
	}

	errOutput := strings.Join(errOutputArr, "")

	if errOutput != "" {
		return errOutput
	}

	cmd.Wait()

	return "Success"
}

// func (s *RunSuit) TestRunScenes(t provider.T) {
// 	sceneMap := map[string]string{
// 		"exactly empty >> ~~":                            "",
// 		"exactly start with abc >> ~f:^abc~":             "abcdvd",
// 		"exactly end with abc >> ~f:abc$~":               "dcdabc",
// 		"exactly contain abc >> ~f:abc~":                 "dvabcd",
// 		"exactly containX abc*3 >> ~f:abc*3~":            "dvabcdabcabcdds",
// 		"exactly 2 in (1,2,3) >> ~f:(1,2,3)~":            "2",
// 		"exactly match %sabc%d >> ~m:%sabc%d~":           "dabc1dad",
// 		"exactly match .*cid=.* >> ~m:.*cid=.*~":         "sdfascid=ljlkjl",
// 		"exactly equal 123 >> ~c:=123~":                  "123",
// 		"exactly less than 123 >> ~c:<123~":              "1",
// 		"exactly less than or equal 123 >> ~c:<=123~":    "123",
// 		"exactly greater than 123 >> ~c:>123~":           "124",
// 		"exactly greater than or equal 123 >> ~c:>=123~": "123",
// 		"exactly not equal 123 >> ~c:<>123~":             "120",
// 		"exactly between 12-19 >> ~c:12-19~":             "15",
// 	}
// 	template := `#!/usr/bin/env php
// 	<?php
// 	/**

// 	title=check string matches pattern
// 	cid=1
// 	pid=1

// 	%s

// 	*/

// 	print("%s\n");`

// 	path := "../../demoo/test_scene.php"
// 	if runtime.GOOS == "windows" {
// 		path = `..\..\demo\test_scene.php`
// 	}
// 	cmd := commonTestHelper.GetZtfPath()+` run ` + path
// 	for expectVal, actualVal := range sceneMap {
// 		file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
// 		if err != nil {
// 			t.Errorf("open file fail, err:%s", err)
// 			return
// 		}
// 		write := bufio.NewWriter(file)
// 		write.WriteString(fmt.Sprintf(template, expectVal, actualVal))
// 		write.Flush()
// 		t.Require().Equal("Success", testRun(cmd, regexp.MustCompile(`Run 1 scripts in \d+ sec, 1\(100\.0%\) Pass, 0\(0\.0%\) Fail, 0\(0\.0%\) Skip|执行1个用例，耗时\d+秒。1\(100\.0%\) 通过，0\(0\.0%\) 失败，0\(0\.0%\) 忽略`)))
// 		file.Close()
// 	}
// 	os.Remove(path)
// }

func TestCliRun(t *testing.T) {
	suite.RunSuite(t, new(RunSuit))
}
