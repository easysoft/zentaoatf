package main

/**
$>ztf junit -p 1 mvn clean package test          执行junit单元测试脚本，
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

	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	"github.com/go-git/go-git/v5"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type RunUnitSuit struct {
	suite.Suite
}

func (s *RunUnitSuit) BeforeEach(t provider.T) {
	if runtime.GOOS == "windows" {
		os.RemoveAll(fmt.Sprintf("%s\\test\\demo\\php\\product1", constTestHelper.RootPath))
	} else {
		os.RemoveAll(fmt.Sprintf("%s/test/demo/php/product1", constTestHelper.RootPath))
	}
	t.AddSubSuite("命令行-run")
}

func (s *RunUnitSuit) TestRunTestng(t provider.T) {
	testngDir := fmt.Sprintf("%stest/demo/ci_test_testng", constTestHelper.RootPath)
	t.ID("5432")
	t.Title("执行TestNG单元测试")
	cloneGit("https://gitee.com/ngtesting/ci_test_testng.git", testngDir)
	t.Require().Equal("Success", testRunUnitTest("mvn clean package test", testngDir, regexp.MustCompile(`Tests run\: 3, Failures\: 0, Errors\: 0, Skipped\: 0`)))
}

func (s *RunUnitSuit) TestRunPytest(t provider.T) {
	pytestDir := fmt.Sprintf(".%stest/demo/ci_test_pytest", constTestHelper.RootPath)
	t.ID("5435")
	t.Title("执行PyTest单元测试")
	cloneGit("https://gitee.com/ngtesting/ci_test_pytest.git", pytestDir)

	t.Require().Equal("Success", testRunUnitTest("pytest --junitxml=testresults.xml", pytestDir, regexp.MustCompile("1 failed, 1 passed")))
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

func TestCliRunUnit(t *testing.T) {
	suite.RunSuite(t, new(RunUnitSuit))
}
