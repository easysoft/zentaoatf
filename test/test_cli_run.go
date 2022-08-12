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
	"time"

	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	"github.com/go-git/go-git/v5"
	expect "github.com/google/goexpect"
)

var (
	scriptResMap = map[string]*regexp.Regexp{
		"ztf run demo/1_string_match_pass.php": regexp.MustCompile("Run 1 scripts in \\d+ sec, 1\\(100\\.0%\\) Pass, 0\\(0\\.0%\\) Fail, 0\\(0\\.0%\\) Skip"),
		"ztf run demo demo/all.cs":             regexp.MustCompile("Run 2 scripts in \\d+ sec, 1\\(50\\.0%\\) Pass, 1\\(50\\.0%\\) Fail, 0\\(0\\.0%\\) Skip"),
		"ztf run demo/001/result.txt":          regexp.MustCompile("Run 1 scripts in \\d+ sec, 0\\(0\\.0%\\) Pass, 1\\(100\\.0%\\) Fail, 0\\(0\\.0%\\) Skip"),
		"ztf run demo -suite 1":                regexp.MustCompile("Run 2 scripts in \\d+ sec, 1\\(50\\.0%\\) Pass, 1\\(50\\.0%\\) Fail, 0\\(0\\.0%\\) Skip"),
		"ztf run demo -task 1":                 regexp.MustCompile("Run 2 scripts in \\d+ sec, 1\\(50\\.0%\\) Pass, 1\\(50\\.0%\\) Fail, 0\\(0\\.0%\\) Skip"),
		"ztf run demo -p 1 -t task1 -cr -cb":   regexp.MustCompile("Submitted test results to ZenTao\\.[\\s\\S]+Success to report bug for case 6"),
	}
)

func testRun(cmd string, expectReg *regexp.Regexp) {
	child, _, err := expect.Spawn(cmd, -1)
	if err != nil {
		fmt.Println(err)
	}
	defer child.Close()

	if out, _, err := child.Expect(expectReg, 10*time.Second); err != nil {
		fmt.Printf("cmd: %s, %s: %s, output: %s\n", cmd, expectReg, "not found", out)
		return
	}

	fmt.Println("Success")
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

func testRunUnitTest(cmdStr, workspacePath string, successRE *regexp.Regexp) (err error) {

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	cmd.Dir = workspacePath

	if cmd == nil {
		fmt.Println("cmd is nil")
		return
	}
	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()

	if err1 != nil {
		fmt.Println(err1)
		return
	} else if err2 != nil {
		fmt.Println(err2)
		return
	}
	cmd.Start()

	isTerminal := false
	reader1 := bufio.NewReader(stdout)
	for {
		line, err3 := reader1.ReadString('\n')
		if line != "" {
			isTerminal = true
			if successRE.MatchString(line) {
				fmt.Println("Success")
				break
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
		fmt.Println(errOutput)
	}

	cmd.Wait()

	return
}
func main() {
	for cmd, expectReg := range scriptResMap {
		if runtime.GOOS == "windows" {
			cmd = strings.ReplaceAll(cmd, "/", "\\")
		}
		testRun(cmd, expectReg)
	}
	testngDir := "./demo/ci_test_testng"
	pytestDir := "./demo/ci_test_pytest"
	if runtime.GOOS == "windows" {
		testngDir = ".\\demo\\ci_test_testng"
		pytestDir = ".\\demo\\ci_test_pytest"
	}
	cloneGit("https://gitee.com/ngtesting/ci_test_testng.git", testngDir)
	testRunUnitTest("mvn clean package test", testngDir, regexp.MustCompile("Tests run\\: 3, Failures\\: 0, Errors\\: 0, Skipped\\: 0"))
	cloneGit("https://gitee.com/ngtesting/ci_test_pytest.git", pytestDir)
	testRunUnitTest("pytest --junitxml=testresult.xml", pytestDir, regexp.MustCompile("1 failed, 1 passed"))
}
