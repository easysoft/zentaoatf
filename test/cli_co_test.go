package main

/**

$>ztf set                                        根据系统提示，设置语言、禅道地址、账号等，Windows下会提示输入语言解释程序。
$>ztf co                                         交互式导出禅道测试用例，将提示用户输入导出类型和编号。
$>ztf co -product 1 -language php             导出编号为1的产品测试用例，使用php语言，缩写-p -l。
$>ztf co -p 1 -m 15 -l php                    导出产品编号为1、模块编号为15的测试用例。
$>ztf co -s 1 -l php -i true                  导出编号为1的套件所含测试用例，期待结果保存在独立文件中。
$>ztf co -t 1 -l php                          导出编号为1的测试单所含用例。

cid=0
pid=0

1.co 导出产品 >> Success
2.co 导出套件 >> Success
3.co 导出任务 >> Success
4.co 参数导出产品 >> Success
5.co 参数导出产品&模块 >> Success
6.co 参数导出套件 >> Success
7.co 参数导出任务 >> Success

*/
import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/stretchr/testify/suite"
)

var (
	coSuccessRe  = regexp.MustCompile("Success")
	typeRe       = regexp.MustCompile("Import test cases from|请选择用例来源")
	productRe    = regexp.MustCompile("Please enter Product Id|请输入 产品Id")
	moduleRe     = regexp.MustCompile("Please enter Module Id|请输入 模块Id")
	suiteRe      = regexp.MustCompile("Please enter Suite Id|请输入 套件Id")
	taskRe       = regexp.MustCompile("Please enter Test Request Id|请输入 测试任务Id")
	separateRe   = regexp.MustCompile("Save expected results in a separate file|是否将用例期待结果保存在独立的文件中")
	storeRe      = regexp.MustCompile("Where to store scripts|请输入脚本保存目录")
	organizeRe   = regexp.MustCompile("Organize test scripts by module|是否希望按模块ID组织脚本目录结构")
	successCoRe  = regexp.MustCompile("Successfully generated \\d+ test scripts|成功创建\\d+个测试脚本")
	languageCoRe = regexp.MustCompile("Select script language|请选择脚本语言")
	productId    = 1
	moduleId     = 0
	taskId       = 1
	suiteId      = 1
	coNewline    = "\n"
)

type CoSuit struct {
	suite.Suite
	testCount uint32
}

func (s *CoSuit) TestCoProduct() {
	assert.Equal(s.Suite.T(), "Success", testCoProduct())
}
func (s *CoSuit) TestCoSuite() {
	assert.Equal(s.Suite.T(), "Success", testCoSuite())
}
func (s *CoSuit) TestCoTask() {
	assert.Equal(s.Suite.T(), "Success", testCoTask())
}
func (s *CoSuit) TestCo() {
	assert.Equal(s.Suite.T(), "Success", testCo(fmt.Sprintf("ztf co -product %d -language php", productId)))
	assert.Equal(s.Suite.T(), "Success", testCo(fmt.Sprintf("ztf co -p %d -m %d -l php", productId, moduleId)))
	assert.Equal(s.Suite.T(), "Success", testCo(fmt.Sprintf("ztf co -s %d -l php -i true", suiteId)))
	assert.Equal(s.Suite.T(), "Success", testCo(fmt.Sprintf("ztf co -t %d -l php", taskId)))
}

func testCoProduct() string {
	cmd := "ztf co"
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()
	if _, err = child.Expect(typeRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", typeRe, err.Error())
	}

	if err = child.Send("1" + coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(productRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", productRe, err.Error())
	}

	if err = child.Send(strconv.Itoa(productId) + coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(moduleRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", moduleRe, err.Error())
	}

	if err = child.Send(strconv.Itoa(moduleId) + coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(separateRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", separateRe, err.Error())
	}

	if err = child.Send("n" + coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(languageCoRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", languageCoRe, err.Error())
	}

	if err = child.Send("5" + coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(storeRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", storeRe, err.Error())
	}
	if err = child.Send(coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(organizeRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", organizeRe, err.Error())
	}

	if err = child.Send(coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(coSuccessRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", coSuccessRe, err.Error())
	}

	return "Success"
}

func testCoSuite() string {
	cmd := "ztf co"
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()
	if _, err = child.Expect(typeRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", typeRe, err.Error())
	}

	if err = child.Send("2" + coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(suiteRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", suiteRe, err.Error())
	}

	if err = child.Send(strconv.Itoa(suiteId) + coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(separateRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", separateRe, err.Error())
	}

	if err = child.Send("n" + coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(languageCoRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", languageCoRe, err.Error())
	}

	if err = child.Send("5" + coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(storeRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", storeRe, err.Error())
	}
	if err = child.Send(coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(organizeRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", organizeRe, err.Error())
	}

	if err = child.Send(coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(coSuccessRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", coSuccessRe, err.Error())
	}

	return "Success"
}

func testCoTask() string {
	cmd := "ztf co"
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()
	if _, err = child.Expect(typeRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", typeRe, err.Error())
	}

	if err = child.Send("3" + coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(taskRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", taskRe, err.Error())
	}

	if err = child.Send(strconv.Itoa(taskId) + coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(separateRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", separateRe, err.Error())
	}

	if err = child.Send("n" + coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(languageCoRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", languageCoRe, err.Error())
	}

	if err = child.Send("5" + coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(storeRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", storeRe, err.Error())
	}
	if err = child.Send(coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(organizeRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", organizeRe, err.Error())
	}

	if err = child.Send(coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(coSuccessRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", coSuccessRe, err.Error())
	}

	return "Success"
}

func testCo(cmd string) string {
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err = child.Expect(storeRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", storeRe, err.Error())
	}
	if err = child.Send(coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(organizeRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", organizeRe, err.Error())
	}

	if err = child.Send(coNewline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(coSuccessRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", coSuccessRe, err.Error())
	}

	return "Success"
}

func TestCo(t *testing.T) {
	if runtime.GOOS == "windows" {
		coNewline = "\r\n"
	}
	suite.Run(t, new(CoSuit))
}
