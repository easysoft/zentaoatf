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
	"time"

	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
)

var (
	typeRe     = regexp.MustCompile("Import test cases from|请选择用例来源")
	productRe  = regexp.MustCompile("Please enter Product Id|请输入 产品Id")
	moduleRe   = regexp.MustCompile("Please enter Module Id|请输入 模块Id")
	suiteRe    = regexp.MustCompile("Please enter Suite Id|请输入 套件Id")
	taskRe     = regexp.MustCompile("Please enter Test Request Id|请输入 测试任务Id")
	separateRe = regexp.MustCompile("Save expected results in a separate file|是否将用例期待结果保存在独立的文件中")
	languageRe = regexp.MustCompile("Select script language|请选择脚本语言")
	storeRe    = regexp.MustCompile("Where to store scripts|请输入脚本保存目录")
	organizeRe = regexp.MustCompile("Organize test scripts by module|是否希望按模块ID组织脚本目录结构")
	successRe  = regexp.MustCompile("Successfully generated \\d+ test scripts|成功创建\\d+个测试脚本")
	productId  = 1
	moduleId   = 0
	taskId     = 1
	suiteId    = 1
	newline    = "\n"
)

func testCoProduct() {
	cmd := "ztf co"
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		fmt.Println(err)
	}
	defer child.Close()
	if _, _, err = child.Expect(typeRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", typeRe, err, newline)
		return
	}

	if err = child.Send("1" + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(productRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", productRe, err, newline)
		return
	}

	if err = child.Send(strconv.Itoa(productId) + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(moduleRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", moduleRe, err, newline)
		return
	}

	if err = child.Send(strconv.Itoa(moduleId) + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(separateRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", separateRe, err, newline)
		return
	}

	if err = child.Send("n" + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(languageRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", languageRe, err, newline)
		return
	}

	if err = child.Send("5" + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(storeRe, 10*time.Second); err != nil {
		fmt.Printf("%s: %s%s", storeRe, err, newline)
		return
	}
	if err = child.Send(newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(organizeRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", organizeRe, err, newline)
		return
	}

	if err = child.Send(newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(successRe, 10*time.Second); err != nil {
		fmt.Printf("%s: %s%s", successRe, err, newline)
		return
	}

	fmt.Println("Success")
}

func testCoSuite() {
	cmd := "ztf co"
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		fmt.Println(err)
	}
	defer child.Close()
	if _, _, err = child.Expect(typeRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", typeRe, err, newline)
		return
	}

	if err = child.Send("2" + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(suiteRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", suiteRe, err, newline)
		return
	}

	if err = child.Send(strconv.Itoa(suiteId) + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(separateRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", separateRe, err, newline)
		return
	}

	if err = child.Send("n" + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(languageRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", languageRe, err, newline)
		return
	}

	if err = child.Send("5" + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(storeRe, 10*time.Second); err != nil {
		fmt.Printf("%s: %s%s", storeRe, err, newline)
		return
	}
	if err = child.Send(newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(organizeRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", organizeRe, err, newline)
		return
	}

	if err = child.Send(newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(successRe, 10*time.Second); err != nil {
		fmt.Printf("%s: %s%s", successRe, err, newline)
		return
	}

	fmt.Println("Success")
}

func testCoTask() {
	cmd := "ztf co"
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		fmt.Println(err)
	}
	defer child.Close()
	if _, _, err = child.Expect(typeRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", typeRe, err, newline)
		return
	}

	if err = child.Send("3" + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(taskRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", taskRe, err, newline)
		return
	}

	if err = child.Send(strconv.Itoa(taskId) + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(separateRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", separateRe, err, newline)
		return
	}

	if err = child.Send("n" + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(languageRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", languageRe, err, newline)
		return
	}

	if err = child.Send("5" + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(storeRe, 10*time.Second); err != nil {
		fmt.Printf("%s: %s%s", storeRe, err, newline)
		return
	}
	if err = child.Send(newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(organizeRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", organizeRe, err, newline)
		return
	}

	if err = child.Send(newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(successRe, 10*time.Second); err != nil {
		fmt.Printf("%s: %s%s", successRe, err, newline)
		return
	}

	fmt.Println("Success")
}

func testCo(cmd string) {
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		fmt.Println(err)
	}
	defer child.Close()

	if _, _, err = child.Expect(storeRe, 10*time.Second); err != nil {
		fmt.Printf("%s: %s%s", storeRe, err, newline)
		return
	}
	if err = child.Send(newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(organizeRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", organizeRe, err, newline)
		return
	}

	if err = child.Send(newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(successRe, 10*time.Second); err != nil {
		fmt.Printf("%s: %s%s", successRe, err, newline)
		return
	}

	fmt.Println("Success")
}
func main() {
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}
	testCoProduct()
	testCoSuite()
	testCoTask()
	testCo(fmt.Sprintf("ztf co -product %d -language php", productId))
	testCo(fmt.Sprintf("ztf co -p %d -m %d -l php", productId, moduleId))
	testCo(fmt.Sprintf("ztf co -s %d -l php -i true", suiteId))
	testCo(fmt.Sprintf("ztf co -t %d -l php", taskId))
}
