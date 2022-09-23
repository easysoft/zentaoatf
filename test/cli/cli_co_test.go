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
	"strconv"
	"testing"
	"time"

	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
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
)

type CoSuite struct {
	suite.Suite
}

func (s *CoSuite) BeforeEach(t provider.T) {
	t.ID("1580")
	t.AddSubSuite("命令行-co")
}
func (s *CoSuite) TestCoProduct(t provider.T) {
	t.Title("导出用例，不提供参数")
	t.Require().Equal("Success", testCoProduct())
}
func (s *CoSuite) TestCoSuite(t provider.T) {
	t.Require().Equal("Success", testCoSuite())
}
func (s *CoSuite) TestCoTask(t provider.T) {
	t.Require().Equal("Success", testCoTask())
}
func (s *CoSuite) TestCo(t provider.T) {
	t.Require().Equal("Success", testCo(fmt.Sprintf(commonTestHelper.GetZtfPath()+" co -product %d -language php", productId)))
	t.Require().Equal("Success", testCo(fmt.Sprintf(commonTestHelper.GetZtfPath()+" co -p %d -m %d -l php", productId, moduleId)))
	t.Require().Equal("Success", testCo(fmt.Sprintf(commonTestHelper.GetZtfPath()+" co -s %d -l php -i true", suiteId)))
	t.Require().Equal("Success", testCo(fmt.Sprintf(commonTestHelper.GetZtfPath()+" co -t %d -l php", taskId)))
}

func testCoProduct() string {
	cmd := commonTestHelper.GetZtfPath() + " co"
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()
	if _, err = child.Expect(typeRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", typeRe, err.Error())
	}

	if err = child.Send("1" + commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(productRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", productRe, err.Error())
	}

	if err = child.Send(strconv.Itoa(productId) + commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(moduleRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", moduleRe, err.Error())
	}

	if err = child.Send(strconv.Itoa(moduleId) + commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(separateRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", separateRe, err.Error())
	}

	if err = child.Send("n" + commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(languageCoRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", languageCoRe, err.Error())
	}

	if err = child.Send("5" + commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(storeRe, time.Second*60); err != nil {
		return fmt.Sprintf("expect %s, actual %s", storeRe, err.Error())
	}
	if err = child.Send(commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(organizeRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", organizeRe, err.Error())
	}

	if err = child.Send(commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(successCoRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successCoRe, err.Error())
	}

	return "Success"
}

func testCoSuite() string {
	cmd := commonTestHelper.GetZtfPath() + " co"
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()
	if _, err = child.Expect(typeRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", typeRe, err.Error())
	}

	if err = child.Send("2" + commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(suiteRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", suiteRe, err.Error())
	}

	if err = child.Send(strconv.Itoa(suiteId) + commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(separateRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", separateRe, err.Error())
	}

	if err = child.Send("n" + commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(languageCoRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", languageCoRe, err.Error())
	}

	if err = child.Send("5" + commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(storeRe, time.Second*60); err != nil {
		return fmt.Sprintf("expect %s, actual %s", storeRe, err.Error())
	}
	if err = child.Send(commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(organizeRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", organizeRe, err.Error())
	}

	if err = child.Send(commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(successCoRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successCoRe, err.Error())
	}

	return "Success"
}

func testCoTask() string {
	cmd := commonTestHelper.GetZtfPath() + " co"
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()
	if _, err = child.Expect(typeRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", typeRe, err.Error())
	}

	if err = child.Send("3" + commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(taskRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", taskRe, err.Error())
	}

	if err = child.Send(strconv.Itoa(taskId) + commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(separateRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", separateRe, err.Error())
	}

	if err = child.Send("n" + commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(languageCoRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", languageCoRe, err.Error())
	}

	if err = child.Send("5" + commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(storeRe, time.Second*60*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", storeRe, err.Error())
	}
	if err = child.Send(commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(organizeRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", organizeRe, err.Error())
	}

	if err = child.Send(commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(successCoRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successCoRe, err.Error())
	}

	return "Success"
}

func testCo(cmd string) string {
	fmt.Println(cmd)
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err = child.Expect(storeRe, time.Second*60); err != nil {
		return fmt.Sprintf("expect %s, actual %s", storeRe, err.Error())
	}
	if err = child.Send(commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(organizeRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", organizeRe, err.Error())
	}

	if err = child.Send(commonTestHelper.NewLine); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(successCoRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successCoRe, err.Error())
	}

	return "Success"
}

func TestCliCo(t *testing.T) {
	suite.RunSuite(t, new(CoSuite))
}
